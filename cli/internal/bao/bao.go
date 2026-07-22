package bao

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
)

type Client struct {
	Url       string
	TokenFunc func() (string, error)
}

type KClient struct {
	Namespace string
}

func (k *KClient) getKNamespace() string {
	if k.Namespace == "" {
		return "common01"
	}
	return k.Namespace
}

func (k *KClient) GetSecretHost() (string, error) {
	ns := k.getKNamespace()
	cout := &bytes.Buffer{}
	cmd := exec.Command(
		"kubectl",
		"-n", ns,
		"get", "ingress", "openbao-"+ns,
		"-o", "jsonpath={.spec.rules[*].host}",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = cout
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("get bao secrets host: %w", err)
	}

	return "https://" + cout.String(), nil
}

func (k *KClient) GetToken() (string, error) {
	ns := k.getKNamespace()
	cout := &bytes.Buffer{}
	cmd := exec.Command(
		"kubectl",
		"-n", ns,
		"get", "secret", "secrets-root-token",
		"-o", "jsonpath={.data.token}",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = cout
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("get bao token: %w", err)
	}

	tkb, err := base64.StdEncoding.DecodeString(cout.String())
	if err != nil {
		return "", fmt.Errorf("get bao token decode secret: %w", err)
	}

	return string(tkb), nil
}

type MountsRes struct {
	Items []MountItem
}

type MountItem struct {
	Name string
	Desc string
}

func (c *Client) doReq(req *http.Request, resBody any) error {

	tk, err := c.TokenFunc()
	if err != nil {
		return err
	}

	if tk != "" {
		req.Header.Set("authorization", "Bearer "+tk)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("bao invalid status code response code=%d", res.StatusCode)
	}

	if resBody != nil {
		err = json.NewDecoder(res.Body).Decode(resBody)
		if err != nil {
			return fmt.Errorf("response decode: %w", err)
		}
	}

	return nil
}

func (c *Client) Mounts() (*MountsRes, error) {
	u, err := url.JoinPath(c.Url, "/v1/sys/mounts")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var resj map[string]any
	err = c.doReq(req, &resj)
	if err != nil {
		return nil, err
	}

	var items []MountItem
	m := resj["data"].(map[string]any)
	for k, v := range m {
		ds, _ := v.(map[string]any)["description"].(string)
		items = append(items, MountItem{k, ds})
	}
	return &MountsRes{items}, nil
}

type KeysMetadataRes struct {
	Items []KeysMetadataItem
}

type KeysMetadataItem struct {
	Name string
}

func (c *Client) KeysList(key string) (*KeysMetadataRes, error) {
	u, err := url.JoinPath(c.Url, "v1", key, "metadata")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u+"?list=true", nil)
	if err != nil {
		return nil, err
	}
	var resj map[string]any
	err = c.doReq(req, &resj)
	if err != nil {
		return nil, err
	}

	var items []KeysMetadataItem
	if arr, ok := resj["data"].(map[string]any)["keys"].([]any); ok {
		for _, v := range arr {
			if name, ok := v.(string); ok {
				items = append(items, KeysMetadataItem{name})
			}
		}
	}
	return &KeysMetadataRes{items}, nil
}

type SecretVersRes struct {
	CurrentVersion int
	OldestVersion  int
	Versions       []SecretVersItem
}

type SecretVersItem struct {
	Version int
}

func (c *Client) SecretVers(key, name string) (*SecretVersRes, error) {
	u, err := url.JoinPath(c.Url, "v1", key, "metadata", name)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	var resj map[string]any
	err = c.doReq(req, &resj)
	if err != nil {
		return nil, err
	}

	var currentVersion int
	cvf, ok := resj["data"].(map[string]any)["current_version"].(float64)
	if !ok {
		return nil, fmt.Errorf("unable to extract current_version")
	}
	currentVersion = int(cvf)

	var oldestVersion int
	ovf, ok := resj["data"].(map[string]any)["oldest_version"].(float64)
	if !ok {
		return nil, fmt.Errorf("unable to extract oldest_version")
	}
	oldestVersion = int(ovf)

	var items []SecretVersItem
	if vs, ok := resj["data"].(map[string]any)["versions"].(map[string]any); ok {
		for k := range vs {
			v, err := strconv.Atoi(k)
			if err != nil {
				return nil, fmt.Errorf("unexpected version string conversino: %w", err)
			}
			items = append(items, SecretVersItem{v})
		}
	}
	return &SecretVersRes{currentVersion, oldestVersion, items}, nil
}

type SecretVerRes struct {
	Items []SecretVerItem
}

type SecretVerItem struct {
	Key   string
	Value string
}

func (c *Client) SecretVer(key, name string, ver int) (*SecretVerRes, error) {
	u, err := url.JoinPath(c.Url, "v1", key, "data", name)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u+"?version="+strconv.Itoa(ver), nil)
	if err != nil {
		return nil, err
	}
	var resj map[string]any
	err = c.doReq(req, &resj)
	if err != nil {
		return nil, err
	}

	var items []SecretVerItem
	if vs, ok := resj["data"].(map[string]any)["data"].(map[string]any); ok {
		for k, v := range vs {
			v, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("value of secret is not a string")
			}
			items = append(items, SecretVerItem{k, v})
		}
	}
	return &SecretVerRes{items}, nil
}

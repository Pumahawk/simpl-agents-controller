package keycloak

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	Server       string
	Realm        string
	ClientId     string
	ClientSecret string
	User         string
	Password     string
}

type TokenizeResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
}

func (c *Client) Tokenize() (*TokenizeResponse, error) {

	u, err := url.JoinPath(c.Server, "/realms/", url.PathEscape(c.Realm), "/protocol/openid-connect/token")
	if err != nil {
		return nil, fmt.Errorf("create base url: %w", err)
	}
	values := url.Values{}
	values["grant_type"] = []string{"password"}
	if c.User != "" {
		values["username"] = []string{c.User}
	}
	if c.Password != "" {
		values["password"] = []string{c.Password}
	}
	if c.ClientId != "" {
		values["client_id"] = []string{c.ClientId}
	}
	if c.ClientSecret != "" {
		values["client_secret"] = []string{c.ClientSecret}
	}
	res, err := http.PostForm(u, values)
	if err != nil {
		return nil, fmt.Errorf("post form: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("tokenize status code %d", res.StatusCode)
	}

	ret := &TokenizeResponse{}
	err = json.NewDecoder(res.Body).Decode(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

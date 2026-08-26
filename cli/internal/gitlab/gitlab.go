package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	BaseUrl string
}

type PackageItem struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	PackageType string   `json:"package_type"`
	CreatedAt   string   `json:"created_at"`
	CreatorId   int      `json:"creator_id"`
	Pipeline    Pipeline `json:"pipeline"`
}

type Pipeline struct {
	Id        int    `json:"id"`
	Status    string `json:"status"`
	Ref       string `json:"ref"`
	Sha       string `json:"sha"`
	WebUrl    string `json:"web_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *Client) Packages(id, page, perpage int) ([]PackageItem, error) {
	rawUrl := c.BaseUrl +
		"/api/v4/projects/" + url.PathEscape(strconv.Itoa(id)) +
		"/packages?" +
		"sort=desc" +
		"&page=" + strconv.Itoa(page) +
		"&perpage=" + strconv.Itoa(perpage)

	res, err := http.Get(rawUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("status_code=%d url=%q", res.StatusCode, rawUrl)
	}

	var items []PackageItem
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

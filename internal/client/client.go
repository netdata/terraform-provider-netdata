package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	AuthToken  string
}

func NewClient(url, authtoken string) *Client {
	c := Client{
		HostURL:    url,
		AuthToken:  "Bearer " + authtoken,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	return &c
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.AuthToken)
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

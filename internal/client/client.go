package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrSpaceIDRequired   = errors.New("spaceID is required")
	ErrChannelIDRequired = errors.New("channelID is required")
	ErrRoomIDRequired    = errors.New("roomID is required")
	ErrMemberIDRequired  = errors.New("memberID is required")
	ErrNodeID            = errors.New("nodeID is required")
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	AuthToken  string
}

func NewClient(url, auth_token string) *Client {
	c := Client{
		HostURL:    url,
		AuthToken:  "Bearer " + auth_token,
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
		return nil, fmt.Errorf("uri: %s, method: %s, status: %d, body: %s", req.URL.RequestURI(), req.Method, res.StatusCode, body)
	}

	return body, err
}

func (c *Client) doRequestUnmarshal(req *http.Request, out any) error {
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &out)
}

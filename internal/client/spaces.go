package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetSpaces() (*[]SpaceInfo, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v3/spaces", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	var spaces []SpaceInfo

	err = c.doRequestUnmarshal(req, &spaces)
	if err != nil {
		return nil, err
	}

	return &spaces, nil
}

func (c *Client) GetSpaceByID(id string) (*SpaceInfo, error) {
	spaces, err := c.GetSpaces()
	if err != nil {
		return nil, err
	}
	for _, space := range *spaces {
		if space.ID == id {
			return &space, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) CreateSpace(name, description string) (*SpaceInfo, error) {
	reqBody, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/spaces", c.HostURL), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	var space SpaceInfo

	err = c.doRequestUnmarshal(req, &space)
	if err != nil {
		return nil, err
	}

	err = c.UpdateSpaceByID(space.ID, name, description)
	if err != nil {
		return nil, err
	}
	space.Name = name
	space.Description = description

	return &space, nil
}

func (c *Client) UpdateSpaceByID(id, name, description string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	reqBody, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/spaces/%s", c.HostURL, id), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteSpaceByID(id string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/spaces/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetSpaceClaimToken(id string) (*string, error) {
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/spaces/%s/tokens", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	err = c.doRequestUnmarshal(req, &data)
	if err != nil {
		return nil, err
	}

	token, ok := data["token"].(string)
	if !ok {
		return nil, fmt.Errorf("token not found")
	}

	return &token, nil
}

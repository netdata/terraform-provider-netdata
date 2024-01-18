package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetRooms(space_id string) (*[]RoomInfo, error) {
	if space_id == "" {
		return nil, fmt.Errorf("space_id is empty")
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/spaces/%s/rooms", c.HostURL, space_id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var rooms []RoomInfo
	err = json.Unmarshal(body, &rooms)
	if err != nil {
		return nil, err
	}

	return &rooms, nil
}

func (c *Client) GetRoomByID(id, space_id string) (*RoomInfo, error) {
	rooms, err := c.GetRooms(space_id)
	if err != nil {
		return nil, err
	}
	for _, room := range *rooms {
		if room.ID == id {
			return &room, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) CreateRoom(space_id, name, description string) (*RoomInfo, error) {
	if space_id == "" {
		return nil, fmt.Errorf("space_id is empty")
	}
	reqBody, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/spaces/%s/rooms", c.HostURL, space_id), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var room RoomInfo
	err = json.Unmarshal(respBody, &room)
	if err != nil {
		return nil, err
	}

	room.Name = name
	room.Description = description

	return &room, nil
}

func (c *Client) UpdateRoomByID(id, space_id, name, description string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	if space_id == "" {
		return fmt.Errorf("space_id is empty")
	}
	reqBody, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/spaces/%s/rooms/%s", c.HostURL, space_id, id), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteRoomByID(id, space_id string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	if space_id == "" {
		return fmt.Errorf("space_id is empty")
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/spaces/%s/rooms/%s", c.HostURL, space_id, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

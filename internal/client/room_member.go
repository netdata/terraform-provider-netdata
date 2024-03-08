package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetRoomMembers(spaceID, roomID string) (*[]RoomMember, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	if roomID == "" {
		return nil, ErrRoomIDRequired
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/rooms/%s/members", c.HostURL, spaceID, roomID), nil)
	if err != nil {
		return nil, err
	}

	var roomMembers []RoomMember

	err = c.doRequestUnmarshal(req, &roomMembers)
	if err != nil {
		return nil, err
	}

	return &roomMembers, nil
}

func (c *Client) GetRoomMemberID(spaceID, roomID, spaceMemberID string) (*RoomMember, error) {
	roomMembers, err := c.GetRoomMembers(spaceID, roomID)
	if err != nil {
		return nil, err
	}
	for _, roomMember := range *roomMembers {
		if roomMember.SpaceMemberID == spaceMemberID {
			return &roomMember, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) CreateRoomMember(spaceID, roomID, spaceMemberID string) error {
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	if roomID == "" {
		return ErrRoomIDRequired
	}
	if spaceMemberID == "" {
		return ErrMemberIDRequired
	}
	reqBody, err := json.Marshal([]string{spaceMemberID})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/spaces/%s/rooms/%s/members", c.HostURL, spaceID, roomID), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil

}

func (c *Client) DeleteRoomMember(spaceID, roomID, spaceMemberID string) error {
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	if roomID == "" {
		return ErrRoomIDRequired
	}
	if spaceMemberID == "" {
		return ErrMemberIDRequired
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v2/spaces/%s/rooms/%s/members?member_ids=%s", c.HostURL, spaceID, roomID, spaceMemberID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

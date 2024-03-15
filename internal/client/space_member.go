package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetSpaceMembers(spaceID string) (*[]SpaceMember, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/members", c.HostURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	var spaceMembers []SpaceMember

	err = c.doRequestUnmarshal(req, &spaceMembers)
	if err != nil {
		return nil, err
	}

	return &spaceMembers, nil
}

func (c *Client) GetSpaceMemberID(spaceID, memberID string) (*SpaceMember, error) {
	spaceMembers, err := c.GetSpaceMembers(spaceID)
	if err != nil {
		return nil, err
	}
	for _, spaceMember := range *spaceMembers {
		if spaceMember.MemberID == memberID {
			return &spaceMember, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) CreateSpaceMember(spaceID, email, role string) (*SpaceMember, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	if email == "" {
		return nil, fmt.Errorf("email is empty")
	}
	if role == "" {
		return nil, fmt.Errorf("role is empty")
	}
	reqBody, err := json.Marshal(map[string]string{"email": email, "role": role})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/spaces/%s/members", c.HostURL, spaceID), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	var spaceMember SpaceMember

	err = c.doRequestUnmarshal(req, &spaceMember)
	if err != nil {
		return nil, err
	}

	return &spaceMember, nil
}

func (c *Client) UpdateSpaceMemberRoleByID(spaceID, memberID, role string) error {
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	if memberID == "" {
		return ErrMemberIDRequired
	}
	if role == "" {
		return fmt.Errorf("role is empty")
	}
	reqBody, err := json.Marshal(map[string]string{"role": role})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v2/spaces/%s/members/%s", c.HostURL, spaceID, memberID), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteSpaceMember(spaceID, memberID string) error {
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	if memberID == "" {
		return ErrMemberIDRequired
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v2/spaces/%s/members?member_ids=%s", c.HostURL, spaceID, memberID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

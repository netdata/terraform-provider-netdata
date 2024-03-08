package client

import (
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetInvitations(spaceID string) (*[]Invitation, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/invitations", c.HostURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	var invitations []Invitation

	err = c.doRequestUnmarshal(req, &invitations)
	if err != nil {
		return nil, err
	}

	return &invitations, nil
}

func (c *Client) DeleteInvitations(spaceID string, invitations *[]Invitation) error {
	if spaceID == "" {
		return ErrSpaceIDRequired
	}

	if len(*invitations) == 0 {
		return nil
	}

	var invitationIDs []string
	for _, invitation := range *invitations {
		invitationIDs = append(invitationIDs, invitation.ID)
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/spaces/%s/invitations?invitation_ids=%s", c.HostURL, spaceID, strings.Join(invitationIDs, ",")), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}

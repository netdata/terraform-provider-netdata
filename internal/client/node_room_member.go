package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetRoomNodes(spaceID, roomID string) (*RoomNodes, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	if roomID == "" {
		return nil, ErrRoomIDRequired
	}

	reqBody := []byte(`{"scope":{"nodes":[]}}`)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v3/spaces/%s/rooms/%s/nodes", c.HostURL, spaceID, roomID), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	var roomNodes RoomNodes

	err = c.doRequestUnmarshal(req, &roomNodes)
	if err != nil {
		return nil, err
	}

	return &roomNodes, nil
}

func (c *Client) GetAllNodes(spaceID string) (*RoomNodes, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	allRooms, err := c.GetRooms(spaceID)
	if err != nil {
		return nil, err
	}

	var allNodesRoomID string
	for _, room := range *allRooms {
		if room.Name == "All nodes" {
			allNodesRoomID = room.ID
			break
		}
	}

	roomNodes, err := c.GetRoomNodes(spaceID, allNodesRoomID)
	if err != nil {
		return nil, err
	}

	return roomNodes, nil
}

func (c *Client) ListNodeMembershipRules(spaceID, roomID string) (*[]NodeMembershipRule, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	if roomID == "" {
		return nil, ErrRoomIDRequired
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v3/spaces/%s/rooms/%s/node-membership-rules", c.HostURL, spaceID, roomID), nil)
	if err != nil {
		return nil, err
	}

	var nodeMembershipRule []NodeMembershipRule

	err = c.doRequestUnmarshal(req, &nodeMembershipRule)
	if err != nil {
		return nil, err
	}

	return &nodeMembershipRule, nil
}

func (c *Client) GetNodeMembershipRule(spaceID, roomID, nodeMembershipID string) (*NodeMembershipRule, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	if roomID == "" {
		return nil, ErrRoomIDRequired
	}
	if nodeMembershipID == "" {
		return nil, ErrNodeMembershipIDRequired
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v3/spaces/%s/rooms/%s/node-membership-rules/%s", c.HostURL, spaceID, roomID, nodeMembershipID), nil)
	if err != nil {
		return nil, err
	}

	var nodeMembershipRule NodeMembershipRule

	err = c.doRequestUnmarshal(req, &nodeMembershipRule)
	if err != nil {
		return nil, err
	}

	return &nodeMembershipRule, nil
}

func (c *Client) CreateNodeMembershipRule(spaceID, roomID, action, description string, clauses []NodeMembershipClause) (*NodeMembershipRule, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	if roomID == "" {
		return nil, ErrRoomIDRequired
	}
	if action == "" {
		return nil, ErrNodeMembershipActionRequired
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"action":      action,
		"description": description,
		"clauses":     clauses,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v3/spaces/%s/rooms/%s/node-membership-rules", c.HostURL, spaceID, roomID), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	var nodeMembershipRule NodeMembershipRule

	err = c.doRequestUnmarshal(req, &nodeMembershipRule)
	if err != nil {
		return nil, err
	}

	return &nodeMembershipRule, nil
}

func (c *Client) UpdateNodeMembershipRule(spaceID, roomID, nodeMembershipID, action, description string, clauses []NodeMembershipClause) (*NodeMembershipRule, error) {
	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}
	if roomID == "" {
		return nil, ErrRoomIDRequired
	}
	if nodeMembershipID == "" {
		return nil, ErrNodeMembershipIDRequired
	}
	if action == "" {
		return nil, ErrNodeMembershipActionRequired
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"action":      action,
		"description": description,
		"clauses":     clauses,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v3/spaces/%s/rooms/%s/node-membership-rules/%s", c.HostURL, spaceID, roomID, nodeMembershipID), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	var nodeMembershipRule NodeMembershipRule

	err = c.doRequestUnmarshal(req, &nodeMembershipRule)
	if err != nil {
		return nil, err
	}

	return &nodeMembershipRule, nil
}

func (c *Client) DeleteNodeMembershipRule(spaceID, roomID, nodeMembershipID string) error {
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	if roomID == "" {
		return ErrRoomIDRequired
	}
	if nodeMembershipID == "" {
		return ErrNodeMembershipIDRequired
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v3/spaces/%s/rooms/%s/node-membership-rules/%s", c.HostURL, spaceID, roomID, nodeMembershipID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateNodeRoomMember(spaceID, roomID, nodeID string) error {
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	if roomID == "" {
		return ErrRoomIDRequired
	}
	if nodeID == "" {
		return ErrNodeID
	}

	reqBody, err := json.Marshal([]string{nodeID})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/spaces/%s/rooms/%s/claimed-nodes", c.HostURL, spaceID, roomID), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil

}

func (c *Client) DeleteNodeRoomMember(spaceID, roomID, nodeID string) error {
	if spaceID == "" {
		return ErrSpaceIDRequired
	}
	if roomID == "" {
		return ErrRoomIDRequired
	}
	if nodeID == "" {
		return ErrNodeID
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/spaces/%s/rooms/%s/claimed-nodes?node_ids=%s", c.HostURL, spaceID, roomID, nodeID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

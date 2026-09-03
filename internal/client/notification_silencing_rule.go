package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) CreateSilencingRule(spaceID string, silencingRule SilencingRule) (*SilencingRule, error) {

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	silencingRuleJson, err := json.Marshal(silencingRule)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/spaces/%s/notifications/silencing/rule", c.HostURL, spaceID), bytes.NewReader(silencingRuleJson))
	if err != nil {
		return nil, err
	}

	var silencingRuleResponse SilencingRule

	err = c.doRequestUnmarshal(req, &silencingRuleResponse)
	if err != nil {
		return nil, err
	}

	return &silencingRuleResponse, nil

}

func (c *Client) GetSilencingRule(silencingRuleID, spaceID string) (*SilencingRule, error) {

	if silencingRuleID == "" {
		return nil, ErrSilencingRuleIDRequired
	}

	parsedSilencingRuleID, err := uuid.Parse(silencingRuleID)
	if err != nil {
		return nil, ErrInvalidSilencingRuleID
	}

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/spaces/%s/notifications/silencing/rules", c.HostURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	var silencingRules []SilencingRule

	err = c.doRequestUnmarshal(req, &silencingRules)
	if err != nil {
		return nil, err
	}

	for _, rule := range silencingRules {
		if rule.ID == parsedSilencingRuleID {
			return &rule, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteSilencingRule(silencingRuleID, spaceID string) error {

	if silencingRuleID == "" {
		return ErrSilencingRuleIDRequired
	}

	if spaceID == "" {
		return ErrSpaceIDRequired
	}

	reqBody, err := json.Marshal([]string{silencingRuleID})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/spaces/%s/notifications/silencing/rules/delete", c.HostURL, spaceID), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil

}

func (c *Client) UpdateSilencingRule(silencingRuleID, spaceID string, silencingRule SilencingRule) (*SilencingRule, error) {

	if silencingRuleID == "" {
		return nil, ErrSilencingRuleIDRequired
	}

	if spaceID == "" {
		return nil, ErrSpaceIDRequired
	}

	silencingRuleJson, err := json.Marshal(silencingRule)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v2/spaces/%s/notifications/silencing/rule/%s", c.HostURL, spaceID, silencingRuleID), bytes.NewReader(silencingRuleJson))
	if err != nil {
		return nil, err
	}

	var silencingRuleResponse SilencingRule

	err = c.doRequestUnmarshal(req, &silencingRuleResponse)
	if err != nil {
		return nil, err
	}

	return &silencingRuleResponse, nil

}

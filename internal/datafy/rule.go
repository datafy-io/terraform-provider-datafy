package datafy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateAccountAutoscalingRuleRequest struct {
	AccountId   string          `json:"account_id"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Active      bool            `json:"active"`
	Rule        json.RawMessage `json:"rule"`
}

type CreateAccountAutoscalingRuleResponse struct {
	AutoscalingRule AutoscalingRule
}

type GetAccountAutoscalingRuleRequest struct {
	AccountId string
	RuleId    string
}

type GetAccountAutoscalingRuleResponse struct {
	AutoscalingRule AutoscalingRule
}

type UpdateAccountAutoscalingRuleRequest struct {
	AccountId   string          `json:"account_id"`
	RuleId      string          `json:"rule_id"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Active      bool            `json:"active"`
	Rule        json.RawMessage `json:"rule"`
}

type UpdateAccountAutoscalingRuleResponse struct {
	AutoscalingRule AutoscalingRule
}

type DeleteAccountAutoscalingRuleRequest struct {
	AccountId string
	RuleId    string
}

type DeleteAccountAutoscalingRuleResponse struct {
}

type AutoscalingRule struct {
	AccountId   string          `json:"accountId"`
	RuleId      string          `json:"ruleId"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Active      bool            `json:"active"`
	Rule        json.RawMessage `json:"rule"`
}

func ruleBody(name, description string, active bool, rule json.RawMessage) map[string]interface{} {
	body := map[string]interface{}{
		"active": active,
		"rule":   rule,
	}
	if name != "" {
		body["name"] = name
	}
	if description != "" {
		body["description"] = description
	}
	return body
}

func (c *Client) CreateAccountAutoscalingRule(ctx context.Context, req *CreateAccountAutoscalingRuleRequest) (*CreateAccountAutoscalingRuleResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/autoscaling/rules", req.AccountId), ruleBody(req.Name, req.Description, req.Active, req.Rule))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var autoscalingRule AutoscalingRule
	if err := json.NewDecoder(resp.Body).Decode(&autoscalingRule); err != nil {
		return nil, err
	}

	return &CreateAccountAutoscalingRuleResponse{
		AutoscalingRule: autoscalingRule,
	}, nil
}

func (c *Client) GetAccountAutoscalingRule(ctx context.Context, req *GetAccountAutoscalingRuleRequest) (*GetAccountAutoscalingRuleResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/autoscaling/rules/%s", req.AccountId, req.RuleId), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var autoscalingRule AutoscalingRule
	if err := json.NewDecoder(resp.Body).Decode(&autoscalingRule); err != nil {
		return nil, err
	}

	return &GetAccountAutoscalingRuleResponse{
		AutoscalingRule: autoscalingRule,
	}, nil
}

func (c *Client) UpdateAccountAutoscalingRule(ctx context.Context, req *UpdateAccountAutoscalingRuleRequest) (*UpdateAccountAutoscalingRuleResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPut, fmt.Sprintf("/api/v1/accounts/%s/autoscaling/rules/%s", req.AccountId, req.RuleId), ruleBody(req.Name, req.Description, req.Active, req.Rule))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var autoscalingRule AutoscalingRule
	if err := json.NewDecoder(resp.Body).Decode(&autoscalingRule); err != nil {
		return nil, err
	}

	return &UpdateAccountAutoscalingRuleResponse{
		AutoscalingRule: autoscalingRule,
	}, nil
}

func (c *Client) DeleteAccountAutoscalingRule(ctx context.Context, req *DeleteAccountAutoscalingRuleRequest) (*DeleteAccountAutoscalingRuleResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/accounts/%s/autoscaling/rules/%s", req.AccountId, req.RuleId), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	return &DeleteAccountAutoscalingRuleResponse{}, nil
}

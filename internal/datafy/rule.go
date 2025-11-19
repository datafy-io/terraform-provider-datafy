package datafy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateAccountAutoscalingRuleRequest struct {
	AccountId string
	Active    bool
	Mode      string
	Rule      string
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
	AccountId string
	RuleId    string
	Active    bool
	Mode      string
	Rule      string
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
	AccountId string `json:"account_id"`
	RuleId    string `json:"rule_id"`
	Active    bool   `json:"active"`
	Mode      string `json:"mode"`
	Rule      string `json:"rule"`
}

func (c *Client) CreateAccountAutoscalingRule(ctx context.Context, req *CreateAccountAutoscalingRuleRequest) (*CreateAccountAutoscalingRuleResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/autoscaling/rules", req.AccountId), map[string]interface{}{
		"active": req.Active,
		"mode":   req.Mode,
		"rule":   req.Rule,
	})
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
	resp, err := c.callAPI(ctx, http.MethodPut, fmt.Sprintf("/api/v1/accounts/%s/autoscaling/rules/%s", req.AccountId, req.RuleId), map[string]interface{}{
		"active": req.Active,
		"mode":   req.Mode,
		"rule":   req.Rule,
	})
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

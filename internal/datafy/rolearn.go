package datafy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateAccountRoleArnRequest struct {
	AccountId string
	Arn       string
}

type CreateAccountRoleArnResponse struct {
	AccountRoleArn AccountRoleArn
}

type GetAccountRoleArnRequest struct {
	AccountId string
}

type GetAccountRoleArnResponse struct {
	AccountRoleArn AccountRoleArn
}

type UpdateAccountRoleArnRequest struct {
	AccountId string
	Arn       string
}

type UpdateAccountRoleArnResponse struct {
	AccountRoleArn AccountRoleArn
}

type DeleteAccountRoleArnRequest struct {
	AccountId string
}

type DeleteAccountRoleArnResponse struct {
}

type AccountRoleArn struct {
	RoleArn string `json:"roleArn"`
}

func (c *Client) CreateAccountRoleArn(ctx context.Context, req *CreateAccountRoleArnRequest) (*CreateAccountRoleArnResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/role-arn", req.AccountId), map[string]interface{}{
		"roleArn": req.Arn,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var accountRoleArn AccountRoleArn
	if err := json.NewDecoder(resp.Body).Decode(&accountRoleArn); err != nil {
		return nil, err
	}

	return &CreateAccountRoleArnResponse{
		AccountRoleArn: accountRoleArn,
	}, nil
}

func (c *Client) GetAccountRoleArn(ctx context.Context, req *GetAccountRoleArnRequest) (*GetAccountRoleArnResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/role-arn", req.AccountId), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var accountRoleArn AccountRoleArn
	if err := json.NewDecoder(resp.Body).Decode(&accountRoleArn); err != nil {
		return nil, err
	}

	return &GetAccountRoleArnResponse{
		AccountRoleArn: accountRoleArn,
	}, nil
}

func (c *Client) UpdateAccountRoleArn(ctx context.Context, req *UpdateAccountRoleArnRequest) (*UpdateAccountRoleArnResponse, error) {
	res, err := c.CreateAccountRoleArn(ctx, &CreateAccountRoleArnRequest{
		AccountId: req.AccountId,
		Arn:       req.Arn,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateAccountRoleArnResponse{
		AccountRoleArn: res.AccountRoleArn,
	}, nil
}

func (c *Client) DeleteAccountRoleArn(ctx context.Context, req *DeleteAccountRoleArnRequest) (*DeleteAccountRoleArnResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/accounts/%s/role-arn", req.AccountId), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	return &DeleteAccountRoleArnResponse{}, nil
}

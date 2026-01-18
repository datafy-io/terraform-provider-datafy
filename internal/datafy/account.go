package datafy

import (
	"context"
	"encoding/json"
	"net/http"
)

type CreateAccountRequest struct {
	AccountName string
}

type CreateAccountResponse struct {
	Account Account
}

type UpdateAccountRequest struct {
	AccountId   string
	AccountName string
}

type UpdateAccountResponse struct {
	Account Account
}

type GetAccountRequest struct {
	AccountId string
}

type GetAccountResponse struct {
	Account Account
}

type DeleteAccountRequest struct {
	AccountId string
}

type DeleteAccountResponse struct {
}

type Account struct {
	AccountId       string `json:"accountId"`
	AccountName     string `json:"accountName"`
	ParentAccountId string `json:"parentAccountId"`
}

func (c *Client) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPost, "/api/v1/accounts", map[string]interface{}{
		"name": req.AccountName,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}

	return &CreateAccountResponse{
		Account: account,
	}, nil
}

func (c *Client) UpdateAccount(ctx context.Context, req *UpdateAccountRequest) (*UpdateAccountResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPut, "/api/v1/accounts/"+req.AccountId, map[string]interface{}{
		"name": req.AccountName,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}

	return &UpdateAccountResponse{
		Account: account,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodGet, "/api/v1/accounts/"+req.AccountId, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}

	return &GetAccountResponse{
		Account: account,
	}, nil
}

func (c *Client) DeleteAccount(ctx context.Context, req *DeleteAccountRequest) (*DeleteAccountResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodDelete, "/api/v1/accounts/"+req.AccountId, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	return &DeleteAccountResponse{}, nil
}

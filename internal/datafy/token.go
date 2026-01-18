package datafy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CreateAccountTokenRequest struct {
	AccountId   string
	Description string
	Ttl         time.Duration
	RoleIds     []string
}

type CreateAccountTokenResponse struct {
	AccountToken AccountToken
}

type GetAccountTokenRequest struct {
	AccountId string
	TokenId   string
}

type GetAccountTokenResponse struct {
	AccountToken AccountToken
}

type DeleteAccountTokenRequest struct {
	AccountId string
	TokenId   string
}

type DeleteAccountTokenResponse struct {
}

type AccountToken struct {
	AccountId   string    `json:"accountId"`
	TokenId     string    `json:"tokenId"`
	Description string    `json:"description"`
	Secret      string    `json:"secret,omitempty"`
	Expires     time.Time `json:"expires"`
	CreatedAt   time.Time `json:"createdAt"`
	RoleIds     []string  `json:"roleIds"`
}

func (c *Client) CreateAccountToken(ctx context.Context, req *CreateAccountTokenRequest) (*CreateAccountTokenResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/tokens", req.AccountId), map[string]interface{}{
		"description":     req.Description,
		"expireInMinutes": int(req.Ttl.Minutes()),
		"roleIds":         req.RoleIds,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var accountToken AccountToken
	if err := json.NewDecoder(resp.Body).Decode(&accountToken); err != nil {
		return nil, err
	}

	return &CreateAccountTokenResponse{
		AccountToken: accountToken,
	}, nil
}

func (c *Client) GetAccountToken(ctx context.Context, req *GetAccountTokenRequest) (*GetAccountTokenResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/tokens/%s", req.AccountId, req.TokenId), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	var accountToken AccountToken
	if err := json.NewDecoder(resp.Body).Decode(&accountToken); err != nil {
		return nil, err
	}

	return &GetAccountTokenResponse{
		AccountToken: accountToken,
	}, nil
}

func (c *Client) DeleteAccountToken(ctx context.Context, req *DeleteAccountTokenRequest) (*DeleteAccountTokenResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/accounts/%s/tokens/%s", req.AccountId, req.TokenId), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	return &DeleteAccountTokenResponse{}, nil
}

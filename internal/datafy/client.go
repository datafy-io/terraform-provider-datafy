package datafy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/datafy-io/terraform-provider-datafy/version"
)

type Client struct {
	token    string
	endpoint string

	httpClient *http.Client
}

func NewClient(token, endpoint string) *Client {
	return &Client{
		token:    token,
		endpoint: endpoint,

		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) callAPI(ctx context.Context, method, path string, headers map[string]string, body map[string]interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.endpoint, path), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("User-Agent", fmt.Sprintf("terraform-provider-datafy/%s (datafy.io)", version.ProviderVersion))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.httpClient.Do(req)
}

func (c *Client) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPost, "/api/v1/accounts", nil, map[string]interface{}{
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

func (c *Client) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodGet, "/api/v1/accounts/"+req.AccountId, nil, nil)
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
	resp, err := c.callAPI(ctx, http.MethodDelete, "/api/v1/accounts/"+req.AccountId, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	return &DeleteAccountResponse{}, nil
}

func (c *Client) CreateAccountRoleArn(ctx context.Context, req *CreateAccountRoleArnRequest) (*CreateAccountRoleArnResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/role-arn", req.AccountId), nil, map[string]interface{}{
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
	resp, err := c.callAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/role-arn", req.AccountId), nil, nil)
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
	resp, err := c.callAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/accounts/%s/role-arn", req.AccountId), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	return &DeleteAccountRoleArnResponse{}, nil
}

func (c *Client) CreateAccountToken(ctx context.Context, req *CreateAccountTokenRequest) (*CreateAccountTokenResponse, error) {
	resp, err := c.callAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/tokens", req.AccountId), nil, map[string]interface{}{
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
	resp, err := c.callAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/tokens/%s", req.AccountId, req.TokenId), nil, nil)
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
	resp, err := c.callAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/accounts/%s/tokens/%s", req.AccountId, req.TokenId), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}

	return &DeleteAccountTokenResponse{}, nil
}

func toError(res *http.Response) error {
	var errMessage struct {
		Message string `json:"message,omitempty"`
	}
	if err := json.NewDecoder(res.Body).Decode(&errMessage); err != nil {
		return err
	}

	return fmt.Errorf("status code %d: %s", res.StatusCode, errMessage.Message)
}

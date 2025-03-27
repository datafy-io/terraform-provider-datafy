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

}

func (c *Client) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {

}

func (c *Client) DeleteAccount(ctx context.Context, req *DeleteAccountRequest) (*DeleteAccountResponse, error) {

}

func (c *Client) CreateAccountRoleArn(ctx context.Context, req *CreateAccountRoleArnRequest) (*CreateAccountRoleArnResponse, error) {

}

func (c *Client) GetAccountRoleArn(ctx context.Context, req *GetAccountRoleArnRequest) (*GetAccountRoleArnResponse, error) {

}

func (c *Client) UpdateAccountRoleArn(ctx context.Context, req *UpdateAccountRoleArnRequest) (*UpdateAccountRoleArnResponse, error) {

}

func (c *Client) DeleteAccountRoleArn(ctx context.Context, req *DeleteAccountRoleArnRequest) (*DeleteAccountRoleArnResponse, error) {

}

func (c *Client) CreateAccountToken(ctx context.Context, req *CreateAccountTokenRequest) (*CreateAccountTokenResponse, error) {

}

func (c *Client) GetAccountToken(ctx context.Context, req *GetAccountTokenRequest) (*GetAccountTokenResponse, error) {

}

func (c *Client) DeleteAccountToken(ctx context.Context, req *DeleteAccountTokenRequest) (*DeleteAccountTokenResponse, error) {

}

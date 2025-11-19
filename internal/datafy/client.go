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

func (c *Client) callAPI(ctx context.Context, method, path string, body map[string]interface{}) (*http.Response, error) {
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

	return c.httpClient.Do(req)
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

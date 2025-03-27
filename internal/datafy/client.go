package datafy

import (
	"net/http"
	"time"
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

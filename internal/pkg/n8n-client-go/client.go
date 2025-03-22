// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

// Package n8n-client-go provides Go client functionalities and data structures for interacting with the n8n automation platform.
//
// It includes representations for workflows, nodes, tags, connections, and related metadata, enabling the
// management of automation tasks.
//
// The package also includes an HTTP client to facilitate communication with the n8n
// service, allowing users to handle workflows, nodes, and other platform features.
package n8n

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a client for the n8n service.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new n8n client.
// It accepts a base URL and an API key for authentication.
//
// Example:
//
//	client := n8n.NewClient("https://example.n8n.io", "your-api-key")
func NewClient(host *string, token *string) (*Client, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	if host == nil {
		return nil, fmt.Errorf("host is required")
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	c.HostURL = *host
	c.Token = *token

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	token := c.Token

	req.Header.Set("X-N8N-API-KEY", token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

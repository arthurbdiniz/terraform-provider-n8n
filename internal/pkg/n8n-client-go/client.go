// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

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

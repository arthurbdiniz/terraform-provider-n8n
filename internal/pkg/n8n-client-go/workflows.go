// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WorkflowsResponse struct {
	Data       []Workflow `json:"data"`
	NextCursor *string    `json:"nextCursor"` // Can be null in JSON, so use a pointer
}

// TODO: Make this method paginated using NextCursor
// GetWorkflows - Returns list of workflows.
func (c *Client) GetWorkflows() (*WorkflowsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/workflows", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	var response WorkflowsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

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

// GetWorkflows - Retrieve all workflows from your instance.
func (c *Client) GetWorkflows() (*WorkflowsResponse, error) {
	var allWorkflows WorkflowsResponse
	cursor := ""

	for {
		url := fmt.Sprintf("%s/api/v1/workflows", c.HostURL)
		// Only append the cursor if it's not empty
		if cursor != "" {
			url = fmt.Sprintf("%s?cursor=%s", url, cursor)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req, nil)
		if err != nil {
			return nil, err
		}

		var workflows WorkflowsResponse
		err = json.Unmarshal(body, &workflows)
		if err != nil {
			return nil, err
		}

		allWorkflows.Data = append(allWorkflows.Data, workflows.Data...)
		if workflows.NextCursor == nil {
			break
		}
		cursor = *workflows.NextCursor
	}

	return &allWorkflows, nil
}

// GetWorkflow - Retrieves a workflow.
func (c *Client) GetWorkflow(workflowID string) (*Workflow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/workflows/%s", c.HostURL, workflowID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	workflow := Workflow{}
	err = json.Unmarshal(body, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

// DeleteWorkflow - Deletes a workflow.
func (c *Client) DeleteWorkflow(workflowID string) (*Workflow, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/workflows/%s", c.HostURL, workflowID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	workflow := Workflow{}
	err = json.Unmarshal(body, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

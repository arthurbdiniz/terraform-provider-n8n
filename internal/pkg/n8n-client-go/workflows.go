// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetWorkflows retrieves all workflows from your n8n instance.
// This method supports pagination and will automatically iterate through
// all available pages by following the cursor in the response.
//
// Returns a pointer to a WorkflowsResponse containing all workflows,
// or an error if the request or response decoding fails.
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

		body, err := c.doRequest(req)
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

// GetWorkflow retrieves the details of a single workflow by its ID.
//
// Parameters:
//   - workflowID: the unique identifier of the workflow.
//
// Returns a pointer to the Workflow struct, or an error if the request or decoding fails.
func (c *Client) GetWorkflow(workflowID string) (*Workflow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/workflows/%s", c.HostURL, workflowID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
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

// DeleteWorkflow deletes a workflow from your n8n instance by its ID.
//
// Parameters:
//   - workflowID: the unique identifier of the workflow to delete.
//
// Returns the deleted Workflow object, or an error if the request or decoding fails.
func (c *Client) DeleteWorkflow(workflowID string) (*Workflow, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/workflows/%s", c.HostURL, workflowID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
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

// DeactivateWorkflow deactivates a workflow by its ID.
// This is typically used to temporarily disable workflow execution.
//
// Parameters:
//   - workflowID: the unique identifier of the workflow to deactivate.
//
// Returns the updated Workflow object, or an error if the request or decoding fails.
func (c *Client) DeactivateWorkflow(workflowID string) (*Workflow, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/workflows/%s/deactivate", c.HostURL, workflowID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workflow := Workflow{}
	if err := json.Unmarshal(body, &workflow); err != nil {
		return nil, err
	}

	return &workflow, nil
}

// ActivateWorkflow activates a workflow by its ID.
// This is typically used to enable workflow execution after creation or deactivation.
//
// Parameters:
//   - workflowID: the unique identifier of the workflow to activate.
//
// Returns the updated Workflow object, or an error if the request or decoding fails.
func (c *Client) ActivateWorkflow(workflowID string) (*Workflow, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/workflows/%s/activate", c.HostURL, workflowID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workflow := Workflow{}
	if err := json.Unmarshal(body, &workflow); err != nil {
		return nil, err
	}

	return &workflow, nil
}

// CreateWorkflow sends a request to create a new workflow in n8n.
// It accepts a CreateWorkflowRequest object and returns the created Workflow with its assigned ID and metadata.
//
// Parameters:
//   - createWorkflowRequest: the workflow data to be created.
//
// Returns the created Workflow object or an error if the request or decoding fails.
func (c *Client) CreateWorkflow(createWorkflowRequest *CreateWorkflowRequest) (*Workflow, error) {
	// Marshal the workflow into JSON
	payload, err := json.Marshal(createWorkflowRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow: %w", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/workflows", c.HostURL), bytes.NewReader(payload))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	body, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Decode the response into a Workflow
	workflow := &Workflow{}
	if err := json.Unmarshal(body, workflow); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return workflow, nil
}

// UpdateWorkflow sends a request to update an existing workflow in n8n.
// It accepts the workflow ID and an UpdateWorkflowRequest object, then returns the updated Workflow with its assigned ID and metadata.
//
// Parameters:
//   - id: the ID of the workflow to be updated.
//   - updateWorkflowRequest: the updated workflow data.
//
// Returns the updated Workflow object or an error if the request or decoding fails.
func (c *Client) UpdateWorkflow(id string, updateWorkflowRequest *UpdateWorkflowRequest) (*Workflow, error) {
	// Marshal the updated workflow into JSON
	payload, err := json.Marshal(updateWorkflowRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated workflow: %w", err)
	}

	// Create the HTTP PUT request
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/workflows/%s", c.HostURL, id), bytes.NewReader(payload))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	body, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Decode the response into a Workflow
	workflow := &Workflow{}
	if err := json.Unmarshal(body, workflow); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return workflow, nil
}

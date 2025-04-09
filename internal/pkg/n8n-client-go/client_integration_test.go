// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/arthurbdiniz/terraform-provider-n8n/internal/config"
	"github.com/arthurbdiniz/terraform-provider-n8n/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationGetWorkflows(t *testing.T) {
	// Create and start the n8n container
	container, url, err := helpers.CreateTestContainer()
	require.NoError(t, err)

	// Ensure container is cleaned up after the test
	defer func() {
		if err := helpers.TerminateContainer(context.Background(), container); err != nil {
			fmt.Println("Error when terminating container:", err)
		}
	}()

	client, err := NewClient(&url, &config.ApiToken)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	workflows, err := client.GetWorkflows()

	if err != nil {
		t.Fatalf("GetWorkflows returned an error: %v", err)
	}

	if len(workflows.Data) != 0 {
		t.Errorf("expected 0 workflows, got %d", len(workflows.Data))
	}
}

func TestIntegrationCreateWorkflow(t *testing.T) {
	// Start the n8n container for testing
	container, url, err := helpers.CreateTestContainer()
	require.NoError(t, err)

	defer func() {
		if err := helpers.TerminateContainer(context.Background(), container); err != nil {
			fmt.Println("Error when terminating container:", err)
		}
	}()

	// Create a client for the running container
	client, err := NewClient(&url, &config.ApiToken)
	require.NoError(t, err)

	// Create a sample workflow
	newWorkflow := &CreateWorkflowRequest{
		Name: "Test Workflow",
		Nodes: []Node{
			{
				ID:          "1",
				Name:        "Start",
				Type:        "n8n-nodes-base.start",
				TypeVersion: 1,
				Position:    []int{0, 0},
				Parameters:  map[string]interface{}{},
			},
		},
		Connections: map[string]Connection{},
		Settings: Settings{
			ExecutionOrder: "v1",
		},
	}

	createdWorkflow, err := client.CreateWorkflow(newWorkflow)

	require.NoError(t, err, "error creating workflow")
	require.NotNil(t, createdWorkflow, "expected non-nil workflow response")
	require.Equal(t, newWorkflow.Name, createdWorkflow.Name, "workflow name should match")
	require.NotEmpty(t, createdWorkflow.ID, "workflow ID should be set")
}

func TestIntegrationCreateWorkflowWithMultiNodeConnections(t *testing.T) {
	// Start the n8n container for testing
	container, url, err := helpers.CreateTestContainer()
	require.NoError(t, err)

	defer func() {
		if err := helpers.TerminateContainer(context.Background(), container); err != nil {
			fmt.Println("Error when terminating container:", err)
		}
	}()

	client, err := NewClient(&url, &config.ApiToken)
	require.NoError(t, err)

	// Define a multi-node workflow
	newWorkflow := &CreateWorkflowRequest{
		Name: "Integration Workflow with Connections",
		Nodes: []Node{
			{
				ID:          "1",
				Name:        "Start",
				Type:        "n8n-nodes-base.start",
				TypeVersion: 1,
				Position:    []int{0, 0},
				Parameters:  map[string]interface{}{},
			},
			{
				ID:          "2",
				Name:        "HTTP Request",
				Type:        "n8n-nodes-base.httpRequest",
				TypeVersion: 1,
				Position:    []int{300, 0},
				Parameters: map[string]interface{}{
					"url":    "https://example.com",
					"method": "GET",
				},
			},
			{
				ID:          "3",
				Name:        "Set",
				Type:        "n8n-nodes-base.set",
				TypeVersion: 1,
				Position:    []int{600, 0},
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{
						"string": []map[string]interface{}{
							{"name": "result", "value": "success"},
						},
					},
				},
			},
		},
		Connections: map[string]Connection{
			"Start": {
				Main: json.RawMessage(`[[{"node":"HTTP Request","type":"main","index":0}]]`),
			},
			"HTTP Request": {
				Main: json.RawMessage(`[[{"node":"Set","type":"main","index":0}]]`),
			},
		},
		Settings: Settings{
			ExecutionOrder: "v1",
		},
	}

	createdWorkflow, err := client.CreateWorkflow(newWorkflow)

	require.NoError(t, err, "error creating workflow")
	require.NotNil(t, createdWorkflow, "expected non-nil workflow response")
	require.Equal(t, newWorkflow.Name, createdWorkflow.Name, "workflow name should match")
	require.NotEmpty(t, createdWorkflow.ID, "workflow ID should be set")
	require.Len(t, createdWorkflow.Nodes, 3, "workflow should have 3 nodes")
}

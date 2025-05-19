// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"encoding/json"
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
	defer helpers.DeferTerminate(container)()

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

	defer helpers.DeferTerminate(container)()

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
			SaveExecutionProgress:    true,
			SaveManualExecutions:     true,
			SaveDataErrorExecution:   "all",
			SaveDataSuccessExecution: "all",
			ExecutionTimeout:         3600,
			Timezone:                 "America/New_York",
			ExecutionOrder:           "v1",
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

	defer helpers.DeferTerminate(container)()

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
			SaveExecutionProgress:    true,
			SaveManualExecutions:     true,
			SaveDataErrorExecution:   "all",
			SaveDataSuccessExecution: "all",
			ExecutionTimeout:         3600,
			Timezone:                 "America/New_York",
			ExecutionOrder:           "v1",
		},
	}

	createdWorkflow, err := client.CreateWorkflow(newWorkflow)

	require.NoError(t, err, "error creating workflow")
	require.NotNil(t, createdWorkflow, "expected non-nil workflow response")
	require.Equal(t, newWorkflow.Name, createdWorkflow.Name, "workflow name should match")
	require.NotEmpty(t, createdWorkflow.ID, "workflow ID should be set")
	require.Len(t, createdWorkflow.Nodes, 3, "workflow should have 3 nodes")
}

func TestIntegrationUpdateWorkflow(t *testing.T) {
	// Start the n8n container for testing
	container, url, err := helpers.CreateTestContainer()
	require.NoError(t, err)

	defer helpers.DeferTerminate(container)()

	client, err := NewClient(&url, &config.ApiToken)
	require.NoError(t, err)

	initialWorkflow := &CreateWorkflowRequest{
		Name: "Original Workflow",
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
			SaveExecutionProgress:    true,
			SaveManualExecutions:     true,
			SaveDataErrorExecution:   "all",
			SaveDataSuccessExecution: "all",
			ExecutionTimeout:         3600,
			Timezone:                 "America/New_York",
			ExecutionOrder:           "v1",
		},
	}

	createdWorkflow, err := client.CreateWorkflow(initialWorkflow)
	require.NoError(t, err, "error creating workflow")
	require.NotNil(t, createdWorkflow, "workflow should be created")

	updateRequest := &UpdateWorkflowRequest{
		Name: "Updated Workflow",
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
				Name:        "Set Node",
				Type:        "n8n-nodes-base.set",
				TypeVersion: 1,
				Position:    []int{300, 0},
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{
						"string": []map[string]interface{}{
							{"name": "key", "value": "value"},
						},
					},
				},
			},
		},
		Connections: map[string]Connection{
			"Start": {
				Main: json.RawMessage(`[[{"node":"Set Node","type":"main","index":0}]]`),
			},
		},
		Settings: Settings{
			SaveExecutionProgress:    true,
			SaveManualExecutions:     true,
			SaveDataErrorExecution:   "all",
			SaveDataSuccessExecution: "all",
			ExecutionTimeout:         3600,
			Timezone:                 "America/New_York",
			ExecutionOrder:           "v1",
		},
	}

	updatedWorkflow, err := client.UpdateWorkflow(createdWorkflow.ID, updateRequest)
	require.NoError(t, err, "error updating workflow")
	require.NotNil(t, updatedWorkflow, "updated workflow should not be nil")
	require.Equal(t, updateRequest.Name, updatedWorkflow.Name, "workflow name should be updated")
	require.Len(t, updatedWorkflow.Nodes, 2, "workflow should have 2 nodes after update")
	require.Equal(t, "Set Node", updatedWorkflow.Nodes[1].Name, "second node name should match")
}

func TestIntegrationDeleteWorkflow(t *testing.T) {
	container, url, err := helpers.CreateTestContainer()
	require.NoError(t, err)

	defer helpers.DeferTerminate(container)()

	client, err := NewClient(&url, &config.ApiToken)
	require.NoError(t, err)

	// Create a workflow to delete
	newWorkflow := &CreateWorkflowRequest{
		Name: "Workflow to Delete",
		Nodes: []Node{{
			ID:          "1",
			Name:        "Start",
			Type:        "n8n-nodes-base.start",
			TypeVersion: 1,
			Position:    []int{0, 0},
			Parameters:  map[string]interface{}{},
		}},
		Connections: map[string]Connection{},
		Settings: Settings{
			ExecutionOrder:           "v1",
			SaveDataErrorExecution:   "all",
			SaveDataSuccessExecution: "all",
		},
	}
	createdWorkflow, err := client.CreateWorkflow(newWorkflow)
	require.NoError(t, err)

	// Delete the workflow
	deleted, err := client.DeleteWorkflow(createdWorkflow.ID)
	require.NoError(t, err)
	require.Equal(t, createdWorkflow.ID, deleted.ID)
}

func TestIntegrationActivateDeactivateWorkflow(t *testing.T) {
	container, url, err := helpers.CreateTestContainer()
	require.NoError(t, err)

	defer helpers.DeferTerminate(container)()

	client, err := NewClient(&url, &config.ApiToken)
	require.NoError(t, err)

	// Create workflow
	newWorkflow := &CreateWorkflowRequest{
		Name: "Workflow to Activate",
		Nodes: []Node{{
			ID:          "1",
			Name:        "Schedule Trigger",
			Type:        "n8n-nodes-base.scheduleTrigger",
			TypeVersion: 1,
			Position:    []int{0, 0},
			Parameters: map[string]interface{}{
				"rule": map[string]interface{}{
					"interval": []interface{}{
						map[string]interface{}{},
					},
				},
			},
		}},
		Connections: map[string]Connection{},
		Settings: Settings{
			ExecutionOrder:           "v1",
			SaveDataErrorExecution:   "all",
			SaveDataSuccessExecution: "all",
		},
	}

	createdWorkflow, err := client.CreateWorkflow(newWorkflow)
	require.NoError(t, err)
	require.False(t, createdWorkflow.Active)

	// Activate the workflow
	activated, err := client.ActivateWorkflow(createdWorkflow.ID)
	require.NoError(t, err)
	require.True(t, activated.Active)
	require.Equal(t, createdWorkflow.ID, activated.ID)

	// Deactivate the workflow
	deactivated, err := client.DeactivateWorkflow(createdWorkflow.ID)
	require.NoError(t, err)
	require.False(t, deactivated.Active)
	require.Equal(t, createdWorkflow.ID, deactivated.ID)
}

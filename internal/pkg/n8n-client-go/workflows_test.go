// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetWorkflows(t *testing.T) {
	mockResponses := []string{
		`{"data": [{"id": "3LODqkaWPmYOi0FA", "name": "Workflow 1"}], "nextCursor": "abc"}`,
		`{"data": [{"id": "if4hSGz1GkaYMLTq", "name": "Workflow 2"}], "nextCursor": null}`,
	}
	requestCount := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		query := r.URL.Query()
		cursor := query.Get("cursor")

		if requestCount == 0 && cursor != "" {
			t.Errorf("expected cursor to be empty, got '%s'", cursor)
		}
		if requestCount == 1 && cursor != "abc" {
			t.Errorf("expected cursor 'abc', got '%s'", cursor)
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponses[requestCount])); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
		requestCount++
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	token := "test-token"
	client, err := NewClient(&ts.URL, &token)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	workflows, err := client.GetWorkflows()
	if err != nil {
		t.Fatalf("GetWorkflows returned an error: %v", err)
	}

	if len(workflows.Data) != 2 {
		t.Errorf("expected 2 workflows, got %d", len(workflows.Data))
	}
}

func TestGetWorkflow(t *testing.T) {
	mockResponse := `{"id": "3LODqkaWPmYOi0FA", "name": "Test Workflow"}`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	token := "test-token"
	client, err := NewClient(&ts.URL, &token)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	workflow, err := client.GetWorkflow("1")
	if err != nil {
		t.Fatalf("GetWorkflow returned an error: %v", err)
	}

	if workflow.ID != "3LODqkaWPmYOi0FA" {
		t.Errorf("unexpected workflow data: %+v", workflow)
	}
}

func TestDeleteWorkflow(t *testing.T) {
	mockID := "3LODqkaWPmYOi0FA"
	mockResponse := `{"id": "3LODqkaWPmYOi0FA", "name": "Test Workflow"}`

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		if path.Base(r.URL.Path) != mockID {
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(`{}`)); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	token := "test-token"
	client, err := NewClient(&ts.URL, &token)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// HTTP 404 - Not Found
	_, err = client.DeleteWorkflow("1")
	if err == nil {
		t.Fatalf("DeleteWorkflow should have returned an HTTP 404 error")
	}

	// HTTP 200 - Workflow deleted
	workflow, err := client.DeleteWorkflow("3LODqkaWPmYOi0FA")
	if err != nil {
		t.Fatalf("DeleteWorkflow returned an error: %v", err)
	}

	if workflow.ID != "3LODqkaWPmYOi0FA" {
		t.Errorf("unexpected workflow data: %+v", workflow)
	}
}

func TestDeactivateWorkflow(t *testing.T) {
	mockResponse := `{"id": "2tUt1wbLX592XDdX", "name": "Deactivated Workflow", "active": false}`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workflows/2tUt1wbLX592XDdX/deactivate" {
			t.Errorf("unexpected URL path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	token := "test-token"
	client, err := NewClient(&ts.URL, &token)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	workflow, err := client.DeactivateWorkflow("2tUt1wbLX592XDdX")
	if err != nil {
		t.Fatalf("DeactivateWorkflow returned an error: %v", err)
	}

	if workflow.ID != "2tUt1wbLX592XDdX" || workflow.Name != "Deactivated Workflow" || workflow.Active {
		t.Errorf("unexpected workflow data: %+v", workflow)
	}
}

func TestActivateWorkflow(t *testing.T) {
	mockResponse := `{"id": "2tUt1wbLX592XDdX", "name": "Activated Workflow", "active": true}`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workflows/2tUt1wbLX592XDdX/activate" {
			t.Errorf("unexpected URL path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	token := "test-token"
	client, err := NewClient(&ts.URL, &token)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	workflow, err := client.ActivateWorkflow("2tUt1wbLX592XDdX")
	if err != nil {
		t.Fatalf("ActivateWorkflow returned an error: %v", err)
	}

	if workflow.ID != "2tUt1wbLX592XDdX" || workflow.Name != "Activated Workflow" || !workflow.Active {
		t.Errorf("unexpected workflow data: %+v", workflow)
	}
}

func TestCreateWorkflow(t *testing.T) {
	mockResponse := `{
		"id": "123456",
		"name": "Test Workflow",
		"active": false,
		"nodes": [{
			"id": "1",
			"name": "Start",
			"type": "n8n-nodes-base.start",
			"typeVersion": 1,
			"position": [0, 0],
			"parameters": {}
		}],
		"connections": {},
		"settings": {
			"executionOrder": "v1"
		},
		"meta": {
			"templateCredsSetupCompleted": false
		},
		"tags": []
	}`

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workflows" {
			t.Errorf("unexpected URL path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	token := "test-token"
	client, err := NewClient(&ts.URL, &token)
	require.NoError(t, err)

	createReq := &CreateWorkflowRequest{
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

	workflow, err := client.CreateWorkflow(createReq)
	require.NoError(t, err)
	require.NotNil(t, workflow)

	if workflow.ID != "123456" || workflow.Name != "Test Workflow" || workflow.Active {
		t.Errorf("unexpected workflow data: %+v", workflow)
	}
}

func TestUpdateWorkflow(t *testing.T) {
	mockResponse := `{
		"id": "123456",
		"name": "Updated Workflow",
		"active": false,
		"nodes": [{
			"id": "1",
			"name": "Start",
			"type": "n8n-nodes-base.start",
			"typeVersion": 1,
			"position": [0, 0],
			"parameters": {}
		}, {
			"id": "2",
			"name": "Set Node",
			"type": "n8n-nodes-base.set",
			"typeVersion": 1,
			"position": [300, 0],
			"parameters": {
				"values": {
					"string": [{
						"name": "key",
						"value": "value"
					}]
				}
			}
		}],
		"connections": {
			"Start": {
				"main": [[{"node": "Set Node", "type": "main", "index": 0}]]
			}
		},
		"settings": {
			"executionOrder": "v1"
		},
		"meta": {
			"templateCredsSetupCompleted": false
		},
		"tags": []
	}`

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workflows/123456" {
			t.Errorf("unexpected URL path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	token := "test-token"
	client, err := NewClient(&ts.URL, &token)
	require.NoError(t, err)

	updateReq := &UpdateWorkflowRequest{
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
			ExecutionOrder: "v1",
		},
	}

	workflow, err := client.UpdateWorkflow("123456", updateReq)
	require.NoError(t, err)
	require.NotNil(t, workflow)

	require.Equal(t, "123456", workflow.ID)
	require.Equal(t, "Updated Workflow", workflow.Name)
	require.Len(t, workflow.Nodes, 2)
	require.Equal(t, "Set Node", workflow.Nodes[1].Name)
}

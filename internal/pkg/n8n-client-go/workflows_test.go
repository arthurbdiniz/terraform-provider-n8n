// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

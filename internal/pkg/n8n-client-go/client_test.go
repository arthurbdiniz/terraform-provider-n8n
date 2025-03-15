// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	host := "http://example.com"
	token := "test-token"

	client, err := NewClient(&host, &token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.HostURL != host {
		t.Errorf("expected host %s, got %s", host, client.HostURL)
	}

	if client.Token != token {
		t.Errorf("expected token %s, got %s", token, client.Token)
	}
}

func TestNewClientWithoutToken(t *testing.T) {
	host := "http://example.com"

	client, err := NewClient(&host, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "token is required" {
		t.Errorf("expected error message 'token is required', got %v", err)
	}

	if client != nil {
		t.Fatal("expected nil client, got a non-nil client")
	}
}

func TestNewClientWithoutHost(t *testing.T) {
	token := "test-token"

	client, err := NewClient(nil, &token)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "host is required" {
		t.Errorf("expected error message 'host is required', got %v", err)
	}

	if client != nil {
		t.Fatal("expected nil client, got a non-nil client")
	}
}

func TestDoRequest(t *testing.T) {
	mockResponse := `{"message": "success"}`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-N8N-API-KEY") != "test-token" {
			t.Errorf("expected token test-token, got %s", r.Header.Get("X-N8N-API-KEY"))
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

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	respBody, err := client.doRequest(req)
	if err != nil {
		t.Fatalf("doRequest returned an error: %v", err)
	}

	if string(respBody) != mockResponse {
		t.Errorf("expected response body %s, got %s", mockResponse, string(respBody))
	}
}

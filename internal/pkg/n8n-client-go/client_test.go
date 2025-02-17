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

	respBody, err := client.doRequest(req, nil)
	if err != nil {
		t.Fatalf("doRequest returned an error: %v", err)
	}

	if string(respBody) != mockResponse {
		t.Errorf("expected response body %s, got %s", mockResponse, string(respBody))
	}
}

func TestDoRequestWithAuthToken(t *testing.T) {
	mockResponse := `{"message": "authorized"}`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-N8N-API-KEY") != "override-token" {
			t.Errorf("expected token override-token, got %s", r.Header.Get("X-N8N-API-KEY"))
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

	overrideToken := "override-token"
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	respBody, err := client.doRequest(req, &overrideToken)
	if err != nil {
		t.Fatalf("doRequest returned an error: %v", err)
	}

	if string(respBody) != mockResponse {
		t.Errorf("expected response body %s, got %s", mockResponse, string(respBody))
	}
}

// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Custom RoundTripper to mock HTTPClient.
type mockRoundTripper struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}
func (e *errorReader) Close() error { return nil }

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

// Helper to create a Client with mocked HTTPClient.
func newMockClient(doFunc func(req *http.Request) (*http.Response, error)) *Client {
	return &Client{
		Token:   "test-token",
		HostURL: "http://example.com",
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{doFunc},
		},
	}
}

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

func TestDoRequest_HTTPClientError(t *testing.T) {
	client := newMockClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("mocked HTTP error")
	})

	req, err := http.NewRequest("GET", client.HostURL+"/test", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = client.doRequest(req)
	if err == nil || !strings.Contains(err.Error(), "mocked HTTP error") {
		t.Fatalf("expected HTTP error, got: %v", err)
	}
}

func TestDoRequest_ReadAllError(t *testing.T) {
	client := newMockClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &errorReader{},
		}, nil
	})

	req, err := http.NewRequest("GET", client.HostURL+"/test", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = client.doRequest(req)
	if err == nil || !strings.Contains(err.Error(), "read error") {
		t.Fatalf("expected read error, got: %v", err)
	}
}

func TestDoRequest_Success(t *testing.T) {
	expected := `{"message": "success"}`

	client := newMockClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(expected)),
		}, nil
	})

	req, err := http.NewRequest("GET", client.HostURL+"/test", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	body, err := client.doRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(body) != expected {
		t.Fatalf("expected body %q, got %q", expected, string(body))
	}
}

func TestDoRequest_Non200StatusCode(t *testing.T) {
	client := newMockClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader("bad request")),
		}, nil
	})

	req, _ := http.NewRequest("GET", client.HostURL+"/test", nil)

	body, err := client.doRequest(req)
	if err == nil {
		t.Fatalf("expected error due to non-200 status code")
	}

	if !strings.Contains(err.Error(), "status: 400") {
		t.Errorf("expected error to contain status code, got: %v", err)
	}

	if body != nil {
		t.Errorf("expected body to be nil on non-200 response")
	}
}

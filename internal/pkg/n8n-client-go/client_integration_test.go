// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestIntegrationDoRequest(t *testing.T) {
	// Create a test container request
	req := testcontainers.ContainerRequest{
		Image:        "docker.n8n.io/n8nio/n8n:latest", // Replace with your API container image
		ExposedPorts: []string{"5678/tcp"},
		WaitingFor:   wait.ForHTTP("/healthz").WithPort("5678").WithStartupTimeout(30 * time.Second),
	}

	ctx := context.Background()

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start container: %s", err)
	}
	defer container.Terminate(ctx) // Ensure container is cleaned up

	// Get the mapped port
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Could not get container host: %s", err)
	}

	port, err := container.MappedPort(ctx, "5678")
	if err != nil {
		t.Fatalf("Could not get mapped port: %s", err)
	}

	apiURL := fmt.Sprintf("http://%s:%s", host, port.Port())

	// Call the API
	resp, err := http.Get(apiURL + "/api/v1/workflows")
	if err != nil {
		t.Fatalf("API request failed: %s", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

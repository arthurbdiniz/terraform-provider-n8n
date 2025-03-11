// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import (
	"context"
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

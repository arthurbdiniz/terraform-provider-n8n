// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/arthurbdiniz/terraform-provider-n8n/internal/config"
	"github.com/arthurbdiniz/terraform-provider-n8n/internal/helpers"
	"github.com/arthurbdiniz/terraform-provider-n8n/internal/pkg/n8n-client-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestWorkflowsDataSource(t *testing.T) {
	// Start the n8n container for testing
	container, url, err := helpers.CreateTestContainer()
	require.NoError(t, err)

	defer helpers.DeferTerminate(container)()

	// Create a client for the running container
	client, err := n8n.NewClient(&url, &config.ApiToken)
	require.NoError(t, err)

	// Create a sample workflow
	newWorkflow := &n8n.CreateWorkflowRequest{
		Name: "Test Workflow",
		Nodes: []n8n.Node{
			{
				ID:          "1",
				Name:        "Start",
				Type:        "n8n-nodes-base.start",
				TypeVersion: 1,
				Position:    []int{0, 0},
				Parameters:  map[string]interface{}{},
			},
		},
		Connections: map[string]n8n.Connection{},
		Settings: n8n.Settings{
			SaveExecutionProgress:    true,
			SaveManualExecutions:     true,
			SaveDataErrorExecution:   "all",
			SaveDataSuccessExecution: "all",
			ExecutionTimeout:         3600,
			ErrorWorkflow:            "",
			Timezone:                 "America/New_York",
			ExecutionOrder:           "v1",
		},
	}

	createdWorkflow, err := client.CreateWorkflow(newWorkflow)

	require.NoError(t, err, "error creating workflow")
	require.NotNil(t, createdWorkflow, "expected non-nil workflow response")
	require.Equal(t, newWorkflow.Name, createdWorkflow.Name, "workflow name should match")
	require.NotEmpty(t, createdWorkflow.ID, "workflow ID should be set")

	t.Logf("n8n test container running at %s", url)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			require.NoError(t, err)
		},
		Steps: []resource.TestStep{
			{
				Config: GetProviderConfig(url) + `
					data "n8n_workflows" "test" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.n8n_workflows.test", "workflows.#"),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.id", createdWorkflow.ID),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.name", createdWorkflow.Name),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.active", fmt.Sprintf("%t", createdWorkflow.Active)),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.created_at", createdWorkflow.CreatedAt),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.updated_at", createdWorkflow.UpdatedAt),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.version_id", createdWorkflow.VersionId),

					// Settings
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.settings.save_execution_progress", fmt.Sprintf("%t", createdWorkflow.Settings.SaveExecutionProgress)),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.settings.save_manual_executions", fmt.Sprintf("%t", createdWorkflow.Settings.SaveManualExecutions)),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.settings.save_data_error_execution", createdWorkflow.Settings.SaveDataErrorExecution),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.settings.save_data_success_execution", createdWorkflow.Settings.SaveDataSuccessExecution),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.settings.execution_timeout", fmt.Sprintf("%d", createdWorkflow.Settings.ExecutionTimeout)),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.settings.error_workflow", createdWorkflow.Settings.ErrorWorkflow),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.settings.timezone", createdWorkflow.Settings.Timezone),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.settings.execution_order", createdWorkflow.Settings.ExecutionOrder),

					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.trigger_count", fmt.Sprintf("%d", createdWorkflow.TriggerCount)),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.connections", "{}"),

					resource.TestCheckNoResourceAttr("data.n8n_workflows.test", "workflows.0.tags"),
					resource.TestCheckResourceAttrSet("data.n8n_workflows.test", "workflows.0.nodes.#"),

					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.nodes.0.id", createdWorkflow.Nodes[0].ID),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.nodes.0.name", createdWorkflow.Nodes[0].Name),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.nodes.0.type", createdWorkflow.Nodes[0].Type),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.nodes.0.type_version", fmt.Sprintf("%g", createdWorkflow.Nodes[0].TypeVersion)),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.nodes.0.position.0", fmt.Sprintf("%d", createdWorkflow.Nodes[0].Position[0])),
					resource.TestCheckResourceAttr("data.n8n_workflows.test", "workflows.0.nodes.0.position.1", fmt.Sprintf("%d", createdWorkflow.Nodes[0].Position[1])),
				),
			},
		},
	})
}

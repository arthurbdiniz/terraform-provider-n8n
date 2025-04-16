// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/arthurbdiniz/terraform-provider-n8n/internal/pkg/n8n-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &workflowDataSource{}
var _ datasource.DataSourceWithConfigure = &workflowDataSource{}

// NewWorkflowDataSource returns a new data source.
func NewWorkflowDataSource() datasource.DataSource {
	return &workflowDataSource{}
}

type workflowDataSource struct {
	client *n8n.Client
}

type workflowDataSourceModel struct {
	ID           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Active       types.Bool     `tfsdk:"active"`
	VersionId    types.String   `tfsdk:"version_id"`
	TriggerCount types.Int64    `tfsdk:"trigger_count"`
	CreatedAt    types.String   `tfsdk:"created_at"`
	UpdatedAt    types.String   `tfsdk:"updated_at"`
	Nodes        []nodesModel   `tfsdk:"nodes"`
	Connections  types.String   `tfsdk:"connections"`
	Settings     *settingsModel `tfsdk:"settings"`
	Tags         []tagsModel    `tfsdk:"tags"`
}

func (d *workflowDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*n8n.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected ProviderData type", fmt.Sprintf("Expected *n8n.Client, got: %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *workflowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (d *workflowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a single workflow by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Workflow ID.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the workflow.",
			},
			"active": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the workflow is currently active.",
			},
			"version_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the current version of the workflow.",
			},
			"trigger_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of times the workflow has been triggered.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the workflow was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the workflow was last updated.",
			},
			"nodes": workflowsNodeAttr(),
			"connections": schema.StringAttribute{
				Computed:    true,
				Description: "JSON-encoded connections data.",
			},
			"settings": workflowsSettingsAttr(),
			"tags":     workflowsTagsAttr(),
		},
	}
}

func (d *workflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state workflowDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowID := state.ID.ValueString()
	workflow, err := d.client.GetWorkflow(workflowID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving workflow", err.Error())
		return
	}

	// Nodes
	var nodes []nodesModel
	for _, node := range workflow.Nodes {
		var positions []types.Int64
		for _, p := range node.Position {
			positions = append(positions, types.Int64Value(int64(p)))
		}

		parameters, err := ConvertToTerraformList(node.Parameters)
		if err != nil {
			resp.Diagnostics.AddError("Error converting node parameters", err.Error())
			return
		}

		nodes = append(nodes, nodesModel{
			ID:          types.StringValue(node.ID),
			Name:        types.StringValue(node.Name),
			Type:        types.StringValue(node.Type),
			TypeVersion: types.Float64Value(float64(node.TypeVersion)),
			Parameters:  parameters,
			Position:    positions,
		})
	}

	// Tags
	var tags []tagsModel
	for _, tag := range workflow.Tags {
		tags = append(tags, tagsModel{
			CreatedAt: types.StringValue(tag.CreatedAt),
			UpdatedAt: types.StringValue(tag.UpdatedAt),
			ID:        types.StringValue(tag.ID),
			Name:      types.StringValue(tag.Name),
		})
	}

	connectionsJSON, err := ConvertConnectionsToTerraformMap(workflow.Connections)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal connections", err.Error())
		return
	}

	state.ID = types.StringValue(workflow.ID)
	state.Name = types.StringValue(workflow.Name)
	state.Active = types.BoolValue(workflow.Active)
	state.VersionId = types.StringValue(workflow.VersionId)
	state.TriggerCount = types.Int64Value(int64(workflow.TriggerCount))
	state.CreatedAt = types.StringValue(workflow.CreatedAt)
	state.UpdatedAt = types.StringValue(workflow.UpdatedAt)
	state.Nodes = nodes
	state.Connections = connectionsJSON
	state.Settings = &settingsModel{
		SaveExecutionProgress:    types.BoolValue(workflow.Settings.SaveExecutionProgress),
		SaveManualExecutions:     types.BoolValue(workflow.Settings.SaveManualExecutions),
		SaveDataErrorExecution:   types.StringValue(workflow.Settings.SaveDataErrorExecution),
		SaveDataSuccessExecution: types.StringValue(workflow.Settings.SaveDataSuccessExecution),
		ExecutionTimeout:         types.Int64Value(int64(workflow.Settings.ExecutionTimeout)),
		ErrorWorkflow:            types.StringValue(workflow.Settings.ErrorWorkflow),
		Timezone:                 types.StringValue(workflow.Settings.Timezone),
		ExecutionOrder:           types.StringValue(workflow.Settings.ExecutionOrder),
	}

	state.Tags = tags

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

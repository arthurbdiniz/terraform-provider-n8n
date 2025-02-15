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
var (
	_ datasource.DataSource              = &workflowsDataSource{}
	_ datasource.DataSourceWithConfigure = &workflowsDataSource{}
)

// NewWorkflowsDataSource is a helper function to simplify the provider implementation.
func NewWorkflowsDataSource() datasource.DataSource {
	return &workflowsDataSource{}
}

// workflowsDataSource is the data source implementation.
type workflowsDataSource struct {
	client *n8n.Client
}

// workflowsDataSourceModel maps the data source schema data.
type workflowsDataSourceModel struct {
	Workflows []workflowsModel `tfsdk:"workflows"`
}

// workflowsModel maps workflows schema data.
type workflowsModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Active       types.Bool   `tfsdk:"active"`
	VersionId    types.String `tfsdk:"version_id"`
	TriggerCount types.Int64  `tfsdk:"trigger_count"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

// Configure adds the provider configured client to the data source.
func (d *workflowsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*n8n.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *n8n.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *workflowsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows"
}

// Schema defines the schema for the data source.
func (d *workflowsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of workflows.",
		Attributes: map[string]schema.Attribute{
			"workflows": schema.ListNestedAttribute{
				Description: "List of workflows available in the system.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the workflow.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the workflow.",
							Computed:    true,
						},
						"active": schema.BoolAttribute{
							Description: "Indicates whether the workflow is currently active.",
							Computed:    true,
						},
						"version_id": schema.StringAttribute{
							Description: "Identifier of the current version of the workflow.",
							Computed:    true,
						},
						"trigger_count": schema.Int64Attribute{
							Description: "Number of times the workflow has been triggered.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Timestamp when the workflow was created.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Timestamp when the workflow was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *workflowsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state workflowsDataSourceModel

	workflowsResponse, err := d.client.GetWorkflows()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read n8n Workflows",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, workflow := range workflowsResponse.Data {
		workflowState := workflowsModel{
			ID:           types.StringValue(workflow.ID),
			Name:         types.StringValue(workflow.Name),
			Active:       types.BoolValue(workflow.Active),
			VersionId:    types.StringValue(workflow.VersionId),
			TriggerCount: types.Int64Value(int64(workflow.TriggerCount)),
			CreatedAt:    types.StringValue(workflow.CreatedAt),
			UpdatedAt:    types.StringValue(workflow.UpdatedAt),
		}

		state.Workflows = append(state.Workflows, workflowState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

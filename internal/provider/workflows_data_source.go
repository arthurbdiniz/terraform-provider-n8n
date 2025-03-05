// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arthurbdiniz/terraform-provider-n8n/internal/pkg/n8n-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	ID           types.String  `tfsdk:"id"`
	Name         types.String  `tfsdk:"name"`
	Active       types.Bool    `tfsdk:"active"`
	VersionId    types.String  `tfsdk:"version_id"`
	TriggerCount types.Int64   `tfsdk:"trigger_count"`
	CreatedAt    types.String  `tfsdk:"created_at"`
	UpdatedAt    types.String  `tfsdk:"updated_at"`
	Nodes        []nodesModel  `tfsdk:"nodes"`
	Connections  types.String  `tfsdk:"connections"`
	Settings     settingsModel `tfsdk:"settings"`
	Meta         metaModel     `tfsdk:"meta"`
	Tags         []tagsModel   `tfsdk:"tags"`
	// PinData      types.Map      `tfsdk:"pin_data"`
	// StaticData   types.Map      `tfsdk:"static_data"`
}

type nodesModel struct {
	ID          types.String     `tfsdk:"id"`
	Name        types.String     `tfsdk:"name"`
	Type        types.String     `tfsdk:"type"`
	TypeVersion types.Float64    `tfsdk:"type_version"`
	Position    []types.Int64    `tfsdk:"position"`
	Parameters  []parameterModel `tfsdk:"parameters"`
}

type tagsModel struct {
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
}

type parameterModel struct {
	Key   types.String `tfsdk:"key"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

type settingsModel struct {
	ExecutionOrder types.String `tfsdk:"execution_order"`
}

type metaModel struct {
	TemplateCredsSetupCompleted types.Bool `tfsdk:"template_creds_setup_completed"`
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
						"nodes": schema.ListNestedAttribute{
							Description: "List of nodes in the workflow.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Node identifier.",
										Computed:    true,
									},
									"name": schema.StringAttribute{
										Description: "Node name.",
										Computed:    true,
									},
									"type": schema.StringAttribute{
										Description: "Type of the node.",
										Computed:    true,
									},
									"type_version": schema.Float64Attribute{
										Description: "Version of the node type.",
										Computed:    true,
									},
									"position": schema.ListAttribute{
										Description: "Position of the node in the workflow.",
										Computed:    true,
										ElementType: types.Int64Type,
									},
									"parameters": schema.ListNestedAttribute{
										Description: "Parameters of the node.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Description: "The parameter key.",
													Computed:    true,
												},
												"type": schema.StringAttribute{
													Description: "The type of the value.",
													Computed:    true,
												},
												"value": schema.StringAttribute{
													Description: "The value as a string.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"connections": schema.StringAttribute{
							Description: "Raw JSON representation of connections between nodes.",
							Computed:    true,
						},
						"settings": schema.SingleNestedAttribute{
							Description: "Workflow settings.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"execution_order": schema.StringAttribute{
									Description: "Execution order setting.",
									Computed:    true,
								},
							},
						},
						"meta": schema.SingleNestedAttribute{
							Description: "Workflow metadata.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"template_creds_setup_completed": schema.BoolAttribute{
									Description: "Indicates if credentials setup is completed.",
									Computed:    true,
								},
							},
						},
						"tags": schema.ListNestedAttribute{
							Description: "Tags associated with the workflow.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"created_at": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"updated_at": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"id": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"name": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ConvertToTerraformList(input map[string]interface{}) ([]parameterModel, error) {
	var paramList []parameterModel

	// Iterate over the input map to convert each element
	for k, v := range input {
		var valueStr string
		var valueType string

		// Determine the type of the value and convert it to string
		switch val := v.(type) {
		case string:
			valueType = "string"
			valueStr = val
		case int:
			valueType = "int"
			valueStr = fmt.Sprintf("%d", val)
		case float64:
			valueType = "float"
			valueStr = fmt.Sprintf("%f", val)
		case bool:
			valueType = "bool"
			valueStr = fmt.Sprintf("%t", val)
		default:
			valueType = "unknown"
			valueStr = fmt.Sprintf("%v", val)
		}

		// Append the parameter to the list
		paramList = append(paramList, parameterModel{
			Key:   types.StringValue(k),
			Type:  types.StringValue(valueType),
			Value: types.StringValue(valueStr),
		})
	}

	return paramList, nil
}

func ConvertConnectionsToTerraformMap(connections interface{}) (types.String, error) {
	data, err := json.Marshal(connections) // Convert to JSON string

	if err != nil {
		return types.StringNull(), err
	}

	return types.StringValue(string(data)), nil
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
		// Convert nodes
		var nodes []nodesModel

		for _, node := range workflow.Nodes {
			var positions []types.Int64
			for _, pos := range node.Position {
				positions = append(positions, types.Int64Value(int64(pos)))
			}
			tflog.Info(ctx, "test")

			params, err := ConvertToTerraformList(node.Parameters)
			if err != nil {
				tflog.Error(ctx, "Error converting parameters", map[string]interface{}{
					"error": err.Error(),
				})
				continue // Skip this node if there's an error
			}

			tflog.Info(ctx, "test2")
			str, _ := json.Marshal(node)
			tflog.Info(ctx, string(str))
			nodes = append(nodes, nodesModel{
				ID:          types.StringValue(node.ID),
				Name:        types.StringValue(node.Name),
				Type:        types.StringValue(node.Type),
				TypeVersion: types.Float64Value(float64(node.TypeVersion)),
				Parameters:  params,
				Position:    positions,
			})
			tflog.Info(ctx, "here")
		}

		// Convert tags
		var tags []tagsModel
		for _, tag := range workflow.Tags {
			tags = append(tags, tagsModel{
				CreatedAt: types.StringValue(tag.CreatedAt),
				UpdatedAt: types.StringValue(tag.UpdatedAt),
				ID:        types.StringValue(tag.ID),
				Name:      types.StringValue(tag.Name),
			})
		}

		// Convert connections
		connectionsJSON, err := ConvertConnectionsToTerraformMap(workflow.Connections)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to process workflow connections",
				err.Error(),
			)
			return
		}

		workflowState := workflowsModel{
			ID:           types.StringValue(workflow.ID),
			Name:         types.StringValue(workflow.Name),
			Active:       types.BoolValue(workflow.Active),
			VersionId:    types.StringValue(workflow.VersionId),
			TriggerCount: types.Int64Value(int64(workflow.TriggerCount)),
			CreatedAt:    types.StringValue(workflow.CreatedAt),
			UpdatedAt:    types.StringValue(workflow.UpdatedAt),
			Nodes:        nodes,
			Connections:  connectionsJSON,
			Settings: settingsModel{
				ExecutionOrder: types.StringValue(workflow.Settings.ExecutionOrder),
			},
			Meta: metaModel{
				TemplateCredsSetupCompleted: types.BoolValue(workflow.Meta.TemplateCredsSetupCompleted),
			},
			Tags: tags,
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

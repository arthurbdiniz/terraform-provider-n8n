// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func workflowsNodeAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
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
	}
}

func workflowsSettingsAttr() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Global execution settings for the workflow.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"save_execution_progress": schema.BoolAttribute{
				Computed:    true,
				Description: "Determines whether the execution progress is saved.",
			},
			"save_manual_executions": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether manual executions are saved.",
			},
			"save_data_error_execution": schema.StringAttribute{
				Computed:    true,
				Description: "Defines the saving behavior for executions with data errors. Options: 'all', 'none'.",
			},
			"save_data_success_execution": schema.StringAttribute{
				Computed:    true,
				Description: "Defines the saving behavior for executions with data success. Options: 'all', 'none'.",
			},
			"execution_timeout": schema.Int64Attribute{
				Computed:    true,
				Description: "Defines the execution timeout in seconds. Max value: 3600.",
			},
			"error_workflow": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the workflow that contains the error trigger node.",
			},
			"timezone": schema.StringAttribute{
				Computed:    true,
				Description: "The timezone for the workflow. Example: 'America/New_York'.",
			},
			"execution_order": schema.StringAttribute{
				Computed:    true,
				Description: "Defines the order in which the workflow nodes are executed. Valid options could include 'v1', 'v2', etc.",
			},
		},
	}
}

func workflowsTagsAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
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
	}
}

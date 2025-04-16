// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowsNodeAttr(t *testing.T) {
	attr := workflowsNodeAttr()

	assert.True(t, attr.Computed)
	assert.Equal(t, "List of nodes in the workflow.", attr.Description)

	attributes := attr.NestedObject.Attributes
	assert.Contains(t, attributes, "id")
	assert.Contains(t, attributes, "name")
	assert.Contains(t, attributes, "type")
	assert.Contains(t, attributes, "type_version")
	assert.Contains(t, attributes, "position")
	assert.Contains(t, attributes, "parameters")

	paramAttr, ok := attributes["parameters"].(schema.ListNestedAttribute)
	assert.True(t, ok)
	assert.Contains(t, paramAttr.NestedObject.Attributes, "key")
	assert.Contains(t, paramAttr.NestedObject.Attributes, "type")
	assert.Contains(t, paramAttr.NestedObject.Attributes, "value")
}

func TestWorkflowsSettingsAttr(t *testing.T) {
	attr := workflowsSettingsAttr()

	assert.True(t, attr.Computed)
	assert.Equal(t, "Global execution settings for the workflow.", attr.Description)

	attrs := attr.Attributes
	expectedKeys := []string{
		"save_execution_progress",
		"save_manual_executions",
		"save_data_error_execution",
		"save_data_success_execution",
		"execution_timeout",
		"error_workflow",
		"timezone",
		"execution_order",
	}
	for _, key := range expectedKeys {
		assert.Contains(t, attrs, key)
	}
}

func TestWorkflowsTagsAttr(t *testing.T) {
	attr := workflowsTagsAttr()

	assert.True(t, attr.Computed)
	assert.Equal(t, "Tags associated with the workflow.", attr.Description)

	attributes := attr.NestedObject.Attributes
	expectedKeys := []string{"created_at", "updated_at", "id", "name"}
	for _, key := range expectedKeys {
		assert.Contains(t, attributes, key)
	}
}

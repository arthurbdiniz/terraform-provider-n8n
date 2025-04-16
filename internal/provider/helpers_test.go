// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertToTerraformList(t *testing.T) {
	input := map[string]interface{}{
		"name":    "test",
		"enabled": true,
		"count":   3,
		"price":   9.99,
	}

	expected := []parameterModel{
		{Key: types.StringValue("name"), Type: types.StringValue("string"), Value: types.StringValue("test")},
		{Key: types.StringValue("enabled"), Type: types.StringValue("bool"), Value: types.StringValue("true")},
		{Key: types.StringValue("count"), Type: types.StringValue("int"), Value: types.StringValue("3")},
		{Key: types.StringValue("price"), Type: types.StringValue("float"), Value: types.StringValue("9.990000")},
	}

	result, err := ConvertToTerraformList(input)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, result)
}

func TestConvertToTerraformList_UnknownType(t *testing.T) {
	input := map[string]interface{}{
		"complex": []string{"unexpected", "type"},
	}

	result, err := ConvertToTerraformList(input)
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	assert.Equal(t, "complex", result[0].Key.ValueString())
	assert.Equal(t, "unknown", result[0].Type.ValueString())
	assert.Equal(t, "[unexpected type]", result[0].Value.ValueString()) // fmt.Sprintf("%v", []string{...})
}

func TestConvertConnectionsToTerraformMap(t *testing.T) {
	input := map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	}

	expected := `{"host":"localhost","port":8080}`

	result, err := ConvertConnectionsToTerraformMap(input)
	assert.NoError(t, err)
	assert.Equal(t, types.StringValue(expected), result)
}

func TestConvertConnectionsToTerraformMap_Error(t *testing.T) {
	ch := make(chan int) // JSON can't marshal channels

	_, err := ConvertConnectionsToTerraformMap(ch)
	assert.Error(t, err)
}

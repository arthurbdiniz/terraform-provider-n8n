// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

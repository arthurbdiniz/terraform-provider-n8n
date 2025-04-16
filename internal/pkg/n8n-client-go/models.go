// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import "encoding/json"

// Workflow represents a workflow in n8n, including metadata, configuration,
// nodes, connections, and tags.
type Workflow struct {
	// ID is the unique identifier of the workflow.
	ID string `json:"id"`

	// Name is the human-readable name of the workflow.
	Name string `json:"name"`

	// Active indicates whether the workflow is currently active.
	Active bool `json:"active"`

	// VersionId is the identifier for the specific version of the workflow.
	VersionId string `json:"versionId"`

	// TriggerCount tracks the number of times the workflow has been triggered.
	TriggerCount int `json:"triggerCount"`

	// CreatedAt is the timestamp when the workflow was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp when the workflow was last updated.
	UpdatedAt string `json:"updatedAt"`

	// Nodes is a list of nodes that define the steps within the workflow.
	Nodes []Node `json:"nodes"`

	// Connections maps node names to their connections, defining how nodes
	// are connected in the workflow.
	Connections map[string]Connection `json:"connections"`

	// Settings contains configuration options for workflow execution.
	Settings Settings `json:"settings"`

	// Tags is a list of tags associated with the workflow for categorization.
	Tags []Tag `json:"tags"`
	// PinData      interface{}           `json:"pinData"`  // TODO understand how this parameter is used and make it exportable to the state
	// StaticData   interface{}           `json:"staticData"` // TODO understand how this parameter is used and make it exportable to the state
}

// WorkflowsResponse represents a paginated response from an API call
// that returns a list of workflows.
type WorkflowsResponse struct {
	// Data contains the list of workflows returned in the response.
	Data []Workflow `json:"data"`

	// NextCursor is an optional cursor string used for pagination.
	// If there are more results to fetch, this field will contain the cursor
	// for the next page. It is nil when there are no additional pages.
	NextCursor *string `json:"nextCursor"`
}

// Tag represents a label assigned to a workflow for organizational purposes.
type Tag struct {
	// CreatedAt is the timestamp when the tag was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp when the tag was last updated.
	UpdatedAt string `json:"updatedAt"`

	// ID is the unique identifier of the tag.
	ID string `json:"id"`

	// Name is the name of the tag.
	Name string `json:"name"`
}

// Connection represents the connections from a node to other nodes within a workflow.
type Connection struct {
	// Main holds the raw connection data. It should be further structured for improved type safety.
	// TODO: Find a way to transform this into a concrete struct.
	Main json.RawMessage `json:"main"`
}

// ConnectionDetail provides detailed information about a specific connection between nodes.
type ConnectionDetail struct {
	// Node is the identifier of the target node in the connection.
	Node string `json:"node"`

	// Type describes the type of connection (e.g., main, conditional).
	Type string `json:"type"`

	// Index is the positional index of the connection in a list.
	Index int `json:"index"`
}

// Node represents an individual step in a workflow, including its configuration and metadata.
type Node struct {
	// Parameters is a map containing node-specific configuration options.
	Parameters map[string]interface{} `json:"parameters"`

	// Type defines the type of the node (e.g., HTTP Request, Set, Code).
	Type string `json:"type"`

	// TypeVersion indicates the version of the node type.
	TypeVersion float64 `json:"typeVersion"`

	// Position is the visual location of the node on the workflow canvas.
	Position []int `json:"position"`

	// ID is the unique identifier of the node.
	ID string `json:"id"`

	// Name is the user-defined name of the node.
	Name string `json:"name"`
}

// Settings contains global execution settings for a workflow.
type Settings struct {
	SaveExecutionProgress    bool   `json:"saveExecutionProgress"`
	SaveManualExecutions     bool   `json:"saveManualExecutions"`
	SaveDataErrorExecution   string `json:"saveDataErrorExecution"`   // Enum: "all", "none"
	SaveDataSuccessExecution string `json:"saveDataSuccessExecution"` // Enum: "all", "none"
	ExecutionTimeout         int    `json:"executionTimeout"`         // maxLength: 3600
	ErrorWorkflow            string `json:"errorWorkflow"`
	Timezone                 string `json:"timezone"`
	ExecutionOrder           string `json:"executionOrder"`
}

// CreateWorkflowRequest defines the allowed fields when creating a workflow.
type CreateWorkflowRequest struct {
	Name        string                `json:"name"`
	Nodes       []Node                `json:"nodes"`
	Connections map[string]Connection `json:"connections"`
	Settings    Settings              `json:"settings"`
	// StaticData   interface{}           `json:"staticData"` // TODO understand how this parameter is used and make it exportable to the state
}

// UpdateWorkflowRequest defines the allowed fields when updating a workflow.
type UpdateWorkflowRequest struct {
	Name        string                `json:"name"`
	Nodes       []Node                `json:"nodes"`
	Connections map[string]Connection `json:"connections"`
	Settings    Settings              `json:"settings"`
	// StaticData   interface{}           `json:"staticData"` // TODO understand how this parameter is used and make it exportable to the state
}

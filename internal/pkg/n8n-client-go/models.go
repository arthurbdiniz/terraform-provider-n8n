// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

import "encoding/json"

type Workflow struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Active       bool                  `json:"active"`
	VersionId    string                `json:"versionId"`
	TriggerCount int                   `json:"triggerCount"`
	CreatedAt    string                `json:"createdAt"`
	UpdatedAt    string                `json:"updatedAt"`
	Nodes        []Node                `json:"nodes"`
	Connections  map[string]Connection `json:"connections"`
	Settings     Settings              `json:"settings"`
	Meta         Meta                  `json:"meta"`
	Tags         []Tag                 `json:"tags"`
	// PinData      interface{}           `json:"pinData"`  // TODO understand how this parameter is used and make it exportable to the state
	// StaticData   interface{}           `json:"staticData"` // TODO understand how this parameter is used and make it exportable to the state
}

type Tag struct {
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	ID        string `json:"id"`
	Name      string `json:"name"`
}

type Connection struct {
	Main json.RawMessage `json:"main"` // TODO find a way to transform this into a struct
}

type ConnectionDetail struct {
	Node  string `json:"node"`
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type Node struct {
	Parameters  map[string]interface{} `json:"parameters"`
	Type        string                 `json:"type"`
	TypeVersion float64                `json:"typeVersion"`
	Position    []int                  `json:"position"`
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
}

type Settings struct {
	ExecutionOrder string `json:"executionOrder"`
}

type Meta struct {
	TemplateCredsSetupCompleted bool `json:"templateCredsSetupCompleted"`
}

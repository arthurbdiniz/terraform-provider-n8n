// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package n8n

type Workflow struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Active       bool        `json:"active"`
	VersionId    string      `json:"versionId"`
	TriggerCount int         `json:"triggerCount"`
	CreatedAt    string      `json:"createdAt"`
	UpdatedAt    string      `json:"updatedAt"`
	Nodes        []Node      `json:"nodes"`
	Connections  interface{} `json:"connections"`
	Settings     Settings    `json:"settings"`
	StaticData   interface{} `json:"staticData"`
	Meta         Meta        `json:"meta"`
	PinData      interface{} `json:"pinData"`
	Tags         []string    `json:"tags"`
}

type Node struct {
	Parameters  map[string]interface{} `json:"parameters"`
	Type        string                 `json:"type"`
	TypeVersion int                    `json:"typeVersion"`
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

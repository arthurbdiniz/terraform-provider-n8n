// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

//go:build ignore
// +build ignore

// This file is intended to be executed via go:generate, not built as part of your project.
//go:generate go run append_docs.go

package main

import (
	"fmt"
	"os"
)

func main() {
	const docFile = "../docs/index.md"

	f, err := os.OpenFile(docFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	customContent := `
### data-sources

- [workflow](./data-sources/workflow.md)
- [workflows](./data-sources/workflows.md)

---

## n8n-client-go

The [n8n-client-go](./client.md) provides a Go client for interacting with the n8n automation platform.

It offers data structures and methods to handle workflows, nodes, tags, connections, and other related metadata.

The client also includes an HTTP client for the communication with the n8n service, enabling you to:

- Create, update, and manage workflows.
- Handle nodes and their configurations within workflows.
- Retrieve and manipulate metadata associated with workflows and connections.
- Integrate the n8n service with your Go-based applications for automation tasks.
`

	if _, err := f.WriteString(customContent); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write content: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("n8n-client-go appended to docs/index.md")
}

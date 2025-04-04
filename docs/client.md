<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# n8n

```go
import "github.com/arthurbdiniz/terraform-provider-n8n/internal/pkg/n8n-client-go"
```

Package n8n\-client\-go provides Go client functionalities and data structures for interacting with the n8n automation platform.

It includes representations for workflows, nodes, tags, connections, and related metadata, enabling the management of automation tasks.

The package also includes an HTTP client to facilitate communication with the n8n service, allowing users to handle workflows, nodes, and other platform features.

## Index

- [type Client](<#Client>)
  - [func NewClient\(host \*string, token \*string\) \(\*Client, error\)](<#NewClient>)
  - [func \(c \*Client\) ActivateWorkflow\(workflowID string\) \(\*Workflow, error\)](<#Client.ActivateWorkflow>)
  - [func \(c \*Client\) DeactivateWorkflow\(workflowID string\) \(\*Workflow, error\)](<#Client.DeactivateWorkflow>)
  - [func \(c \*Client\) DeleteWorkflow\(workflowID string\) \(\*Workflow, error\)](<#Client.DeleteWorkflow>)
  - [func \(c \*Client\) GetWorkflow\(workflowID string\) \(\*Workflow, error\)](<#Client.GetWorkflow>)
  - [func \(c \*Client\) GetWorkflows\(\) \(\*WorkflowsResponse, error\)](<#Client.GetWorkflows>)
- [type Connection](<#Connection>)
- [type ConnectionDetail](<#ConnectionDetail>)
- [type Meta](<#Meta>)
- [type Node](<#Node>)
- [type Settings](<#Settings>)
- [type Tag](<#Tag>)
- [type Workflow](<#Workflow>)
- [type WorkflowsResponse](<#WorkflowsResponse>)


<a name="Client"></a>
## type Client

Client represents a client for the n8n service.

```go
type Client struct {
    HostURL    string
    HTTPClient *http.Client
    Token      string
}
```

<a name="NewClient"></a>
### func NewClient

```go
func NewClient(host *string, token *string) (*Client, error)
```

NewClient creates a new n8n client. It accepts a base URL and an API key for authentication.

Example:

```
client := n8n.NewClient("https://example.n8n.io", "your-api-key")
```

<a name="Client.ActivateWorkflow"></a>
### func \(\*Client\) ActivateWorkflow

```go
func (c *Client) ActivateWorkflow(workflowID string) (*Workflow, error)
```

ActivateWorkflow activates a workflow by its ID. This is typically used to enable workflow execution after creation or deactivation.

Parameters:

- workflowID: the unique identifier of the workflow to activate.

Returns the updated Workflow object, or an error if the request or decoding fails.

<a name="Client.DeactivateWorkflow"></a>
### func \(\*Client\) DeactivateWorkflow

```go
func (c *Client) DeactivateWorkflow(workflowID string) (*Workflow, error)
```

DeactivateWorkflow deactivates a workflow by its ID. This is typically used to temporarily disable workflow execution.

Parameters:

- workflowID: the unique identifier of the workflow to deactivate.

Returns the updated Workflow object, or an error if the request or decoding fails.

<a name="Client.DeleteWorkflow"></a>
### func \(\*Client\) DeleteWorkflow

```go
func (c *Client) DeleteWorkflow(workflowID string) (*Workflow, error)
```

DeleteWorkflow deletes a workflow from your n8n instance by its ID.

Parameters:

- workflowID: the unique identifier of the workflow to delete.

Returns the deleted Workflow object, or an error if the request or decoding fails.

<a name="Client.GetWorkflow"></a>
### func \(\*Client\) GetWorkflow

```go
func (c *Client) GetWorkflow(workflowID string) (*Workflow, error)
```

GetWorkflow retrieves the details of a single workflow by its ID.

Parameters:

- workflowID: the unique identifier of the workflow.

Returns a pointer to the Workflow struct, or an error if the request or decoding fails.

<a name="Client.GetWorkflows"></a>
### func \(\*Client\) GetWorkflows

```go
func (c *Client) GetWorkflows() (*WorkflowsResponse, error)
```

GetWorkflows retrieves all workflows from your n8n instance. This method supports pagination and will automatically iterate through all available pages by following the cursor in the response.

Returns a pointer to a WorkflowsResponse containing all workflows, or an error if the request or response decoding fails.

<a name="Connection"></a>
## type Connection

Connection represents the connections from a node to other nodes within a workflow.

```go
type Connection struct {
    // Main holds the raw connection data. It should be further structured for improved type safety.
    // TODO: Find a way to transform this into a concrete struct.
    Main json.RawMessage `json:"main"`
}
```

<a name="ConnectionDetail"></a>
## type ConnectionDetail

ConnectionDetail provides detailed information about a specific connection between nodes.

```go
type ConnectionDetail struct {
    // Node is the identifier of the target node in the connection.
    Node string `json:"node"`

    // Type describes the type of connection (e.g., main, conditional).
    Type string `json:"type"`

    // Index is the positional index of the connection in a list.
    Index int `json:"index"`
}
```

<a name="Meta"></a>
## type Meta

Meta contains additional metadata about the workflow's setup status.

```go
type Meta struct {
    // TemplateCredsSetupCompleted indicates whether the setup for template credentials is complete.
    TemplateCredsSetupCompleted bool `json:"templateCredsSetupCompleted"`
}
```

<a name="Node"></a>
## type Node

Node represents an individual step in a workflow, including its configuration and metadata.

```go
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
    ID  string `json:"id"`

    // Name is the user-defined name of the node.
    Name string `json:"name"`
}
```

<a name="Settings"></a>
## type Settings

Settings contains global execution settings for a workflow.

```go
type Settings struct {
    // ExecutionOrder defines how the workflow nodes should be executed.
    ExecutionOrder string `json:"executionOrder"`
}
```

<a name="Tag"></a>
## type Tag

Tag represents a label assigned to a workflow for organizational purposes.

```go
type Tag struct {
    // CreatedAt is the timestamp when the tag was created.
    CreatedAt string `json:"createdAt"`

    // UpdatedAt is the timestamp when the tag was last updated.
    UpdatedAt string `json:"updatedAt"`

    // ID is the unique identifier of the tag.
    ID  string `json:"id"`

    // Name is the name of the tag.
    Name string `json:"name"`
}
```

<a name="Workflow"></a>
## type Workflow

Workflow represents a workflow in n8n, including metadata, configuration, nodes, connections, and tags.

```go
type Workflow struct {
    // ID is the unique identifier of the workflow.
    ID  string `json:"id"`

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

    // Meta provides metadata about the workflow.
    Meta Meta `json:"meta"`

    // Tags is a list of tags associated with the workflow for categorization.
    Tags []Tag `json:"tags"`
}
```

<a name="WorkflowsResponse"></a>
## type WorkflowsResponse

WorkflowsResponse represents a paginated response from an API call that returns a list of workflows.

```go
type WorkflowsResponse struct {
    // Data contains the list of workflows returned in the response.
    Data []Workflow `json:"data"`

    // NextCursor is an optional cursor string used for pagination.
    // If there are more results to fetch, this field will contain the cursor
    // for the next page. It is nil when there are no additional pages.
    NextCursor *string `json:"nextCursor"`
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)

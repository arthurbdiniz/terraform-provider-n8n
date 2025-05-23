---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "n8n_workflows Data Source - n8n"
subcategory: ""
description: |-
  Fetches the list of workflows.
---

# n8n_workflows (Data Source)

Fetches the list of workflows.

## Example Usage

```terraform
# List all workflows.
data "n8n_workflows" "all" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `workflows` (Attributes List) List of workflows available in the system. (see [below for nested schema](#nestedatt--workflows))

<a id="nestedatt--workflows"></a>
### Nested Schema for `workflows`

Read-Only:

- `active` (Boolean) Indicates whether the workflow is currently active.
- `connections` (String) Raw JSON representation of connections between nodes.
- `created_at` (String) Timestamp when the workflow was created.
- `id` (String) Unique identifier of the workflow.
- `name` (String) Name of the workflow.
- `nodes` (Attributes List) List of nodes in the workflow. (see [below for nested schema](#nestedatt--workflows--nodes))
- `settings` (Attributes) Global execution settings for the workflow. (see [below for nested schema](#nestedatt--workflows--settings))
- `tags` (Attributes List) Tags associated with the workflow. (see [below for nested schema](#nestedatt--workflows--tags))
- `trigger_count` (Number) Number of times the workflow has been triggered.
- `updated_at` (String) Timestamp when the workflow was last updated.
- `version_id` (String) Identifier of the current version of the workflow.

<a id="nestedatt--workflows--nodes"></a>
### Nested Schema for `workflows.nodes`

Read-Only:

- `id` (String) Node identifier.
- `name` (String) Node name.
- `parameters` (Attributes List) Parameters of the node. (see [below for nested schema](#nestedatt--workflows--nodes--parameters))
- `position` (List of Number) Position of the node in the workflow.
- `type` (String) Type of the node.
- `type_version` (Number) Version of the node type.

<a id="nestedatt--workflows--nodes--parameters"></a>
### Nested Schema for `workflows.nodes.parameters`

Read-Only:

- `key` (String) The parameter key.
- `type` (String) The type of the value.
- `value` (String) The value as a string.



<a id="nestedatt--workflows--settings"></a>
### Nested Schema for `workflows.settings`

Read-Only:

- `error_workflow` (String) The ID of the workflow that contains the error trigger node.
- `execution_order` (String) Defines the order in which the workflow nodes are executed. Valid options could include 'v1', 'v2', etc.
- `execution_timeout` (Number) Defines the execution timeout in seconds. Max value: 3600.
- `save_data_error_execution` (String) Defines the saving behavior for executions with data errors. Options: 'all', 'none'.
- `save_data_success_execution` (String) Defines the saving behavior for executions with data success. Options: 'all', 'none'.
- `save_execution_progress` (Boolean) Determines whether the execution progress is saved.
- `save_manual_executions` (Boolean) Indicates whether manual executions are saved.
- `timezone` (String) The timezone for the workflow. Example: 'America/New_York'.


<a id="nestedatt--workflows--tags"></a>
### Nested Schema for `workflows.tags`

Read-Only:

- `created_at` (String)
- `id` (String)
- `name` (String)
- `updated_at` (String)

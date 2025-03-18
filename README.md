# terraform-provider-n8n

![Build Status](https://img.shields.io/github/check-runs/arthurbdiniz/terraform-provider-n8n/main?label=build) ![Codecov](https://img.shields.io/codecov/c/github/arthurbdiniz/terraform-provider-n8n) ![License](https://img.shields.io/github/license/arthurbdiniz/terraform-provider-n8n)

This repository contains the Terraform provider for [n8n](https://n8n.io), an open-source workflow automation tool. With this provider, you can manage and automate the configuration of n8n resources directly from your Terraform infrastructure-as-code setups.

> ⚠️ The n8n Terraform provider it's `under development` but the goal for the future is to allow you to manage workflows, nodes, credentials, and other resources, enabling an approach to automation workflows.

Feel free to open issues, contribute, and explore the provider's functionality.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) v1.0+ installed locally.
- [Go 1.22+](https://golang.org/doc/install) installed and configured.
- [Docker](https://www.docker.com/products/docker-desktop) and [Docker Compose](https://docs.docker.com/compose/install/) to run an instance of n8n locally.

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Once the provider gets published to Terraform registry you will be able to use the provider by:

```
terraform {
  required_providers {
    n8n = {
      source = "registry.terraform.io/arthurbdiniz/n8n"
    }
  }
}

provider "n8n" {
  host  = "http://localhost:5678"
  token = "..."
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Prepare Terraform for local provider install

Terraform installs providers and verifies their versions and checksums when you run `terraform init`. Terraform will download your providers from either the provider registry or a local registry. However, while building your provider you will want to test Terraform configuration against a local development build of the provider. The development build will not have an associated version number or an official set of checksums listed in a provider registry.

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

Terraform searches for the `.terraformrc` file in your home directory and applies any configuration settings you set.

First, find the `GOBIN` path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.

If the `GOBIN` go environment variable is not set, use the default path, `/home/<Username>/go/bin` (Linux) or `/Users/<Username>/go/bin` (Mac).

```shell
# Linux
export GOBIN=/home/<Username>/go/bin

# MacOS
export GOBIN=/Users/<Username>/go/bin
```

Create a new file called `.terraformrc` in your home directory (`~`), then add the `dev_overrides` block below. Change the `<PATH>` to the value returned from the `go env GOBIN` command above.

#### ~/.terraformrc

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/arthurbdiniz/n8n" = "/home/arthurbdiniz/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Start n8n locally

This provider requires a running instance of n8n.

In another terminal window, navigate to the `docker_compose` directory.

```shell
cd docker_compose
```

Run `docker-compose up` to spin up a local instance of n8n on port `5678`.

```shell
docker-compose up
```

Leave this process running in your terminal window. The n8n service will print out log messages in this terminal.

Open your browser and navigate to http://localhost:5678

Create the initial admin account for n8n.

Then create a API token on http://localhost:5678/settings/api after login.

Save this token for later when we declare the terrorm provider.

In the original terminal window, verify that n8n API is running by sending a request to its `/api/v1/workflows` endpoint. The n8n service will respond with the list of workflows (most likely an empty list of workflows since we don't have any created yet).


```shell
curl -X 'GET' \
  'http://localhost:5678/api/v1/workflows' \
  -H 'accept: application/json' \
  -H "X-N8N-API-KEY: $N8N_API_TOKEN"
```

## Locally install provider and verify with Terraform

Your Terraform CLI is now ready to use the locally installed provider in the `GOBIN` path. Use the `go install` command from the example repository's root directory to compile the provider into a binary and install it in your `GOBIN` path.

```shell
go install .
```

Create an `examples/provider-install-verification` directory, which will contain a terraform configuration to verify local provider installation, and navigate to it.

```shell
mkdir examples/provider-install-verification && cd "$_"
```

Create a `main.tf` file with the following.

```
terraform {
  required_providers {
    n8n = {
      source = "registry.terraform.io/arthurbdiniz/n8n"
    }
  }
}

provider "n8n" {
  host  = "http://localhost:5678"
  token = "..."
}

data "n8n_workflows" "test" {}

output "n8n_workflows" {
  value = data.n8n_workflows.test
}
```

The `main.tf` Terraform configuration file in this directory uses a "n8n_workflows" data source from the provider.

Run a Terraform plan with the data source. Terraform will respond with empty list of workflows.

```shell
terraform plan
```

> Output
```log
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/edu/hashicups in /home/arthurbdiniz/go/bin
│  - arthurbdiniz/n8n in /home/arthurbdiniz/go/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.n8n_workflows.test: Reading...
data.n8n_workflows.test: Read complete after 0s

Changes to Outputs:
  + n8n_workflows = {
      + workflows = null
    }
```

Navigate to the `terraform-provider-n8n` directory and you now can start developing the provider.

```shell
cd ../..
```

## License

Apache 2 Licensed. See [LICENSE](https://github.com/arthurbdiniz/terraform-provider-n8n/blob/master/LICENSE) for full details.

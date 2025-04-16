// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"os"

	"github.com/arthurbdiniz/terraform-provider-n8n/internal/pkg/n8n-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &n8nProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &n8nProvider{
			version: version,
		}
	}
}

// n8nProviderModel maps provider schema data to a Go type.
type n8nProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

// n8nProvider is the provider implementation.
type n8nProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *n8nProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "n8n"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *n8nProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with n8n.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for n8n API. May also be provided via `N8N_HOST` environment variable.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Token for n8n API. May also be provided via `N8N_TOKEN` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *n8nProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring n8n client")

	// Retrieve provider data from configuration
	var config n8nProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown n8n API Host",
			"The provider cannot create the n8n API client as there is an unknown configuration value for the n8n API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the N8N_HOST environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown n8n API Username",
			"The provider cannot create the n8n API client as there is an unknown configuration value for the n8n API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the N8N_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("N8N_HOST")
	token := os.Getenv("N8N_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing n8n API Host",
			"The provider cannot create the n8n API client as there is a missing or empty value for the n8n API host. "+
				"Set the host value in the configuration or use the N8N_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing n8n API Token",
			"The provider cannot create the n8n API client as there is a missing or empty value for the n8n API token. "+
				"Set the token value in the configuration or use the N8N_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "n8n_host", host)
	ctx = tflog.SetField(ctx, "n8n_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "n8n_token")

	tflog.Debug(ctx, "Creating n8n client")

	// Create a new n8n client using the configuration values
	client, err := n8n.NewClient(&host, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create n8n API Client",
			"An unexpected error occurred when creating the n8n API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"n8n Client Error: "+err.Error(),
		)
		return
	}

	// Make the n8n client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured n8n client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *n8nProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkflowsDataSource,
		NewWorkflowDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *n8nProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// Functions defines the functions implemented in the provider.
func (p *n8nProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}

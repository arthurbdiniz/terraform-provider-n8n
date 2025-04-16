// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"fmt"

	"github.com/arthurbdiniz/terraform-provider-n8n/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	//
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"n8n": providerserver.NewProtocol6WithError(New("test")()),
	}
)

// GetProviderConfig returns the provider configuration string for tests.
// providerConfig is a shared configuration to combine with the actual
// test configuration so the n8n client is properly configured.
// It is also possible to use the N8N_ environment variables instead,
// such as updating the Makefile and running the testing through that tool.
func GetProviderConfig(url string) string {
	return fmt.Sprintf(`
provider "n8n" {
  host  = "%s"
  token = "%s"
}
`, url, config.ApiToken)
}

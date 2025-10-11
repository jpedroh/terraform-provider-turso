// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var providerConfig = fmt.Sprintf(`
provider "turso" {
  api_token = "%s"
}
`, os.Getenv("TURSO_API_TOKEN"))

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"turso": providerserver.NewProtocol6WithError(New("test")()),
	}
)

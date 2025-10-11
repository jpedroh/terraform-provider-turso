// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseTokenResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "turso_database_token" "test" {
					organization_name = "jpedroh"
					database_name	  = "tfproviderdatasource"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("turso_database_token.test", "organization_name", "jpedroh"),
					resource.TestCheckResourceAttr("turso_database_token.test", "database_name", "tfproviderdatasource"),
					resource.TestCheckResourceAttr("turso_database_token.test", "authorization", "full-access"),
					resource.TestCheckResourceAttrSet("turso_database_token.test", "jwt"),
				),
			},
			{
				Config: providerConfig + `
				resource "turso_database_token" "test" {
					organization_name = "jpedroh"
					database_name	  = "tfproviderdatasource"
					authorization = "read-only"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("turso_database_token.test", "organization_name", "jpedroh"),
					resource.TestCheckResourceAttr("turso_database_token.test", "database_name", "tfproviderdatasource"),
					resource.TestCheckResourceAttr("turso_database_token.test", "authorization", "read-only"),
					resource.TestCheckResourceAttrSet("turso_database_token.test", "jwt"),
				),
			},
		},
	})
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseInstanceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "turso_database_instance" "test" {
					organization_slug = "jpedroh"
					database_name = "tfproviderdatasource"
					name = "aws-us-east-1"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.turso_database_instance.test", "organization_slug", "jpedroh"),
					resource.TestCheckResourceAttr("data.turso_database_instance.test", "database_name", "tfproviderdatasource"),
					resource.TestCheckResourceAttr("data.turso_database_instance.test", "name", "aws-us-east-1"),
					resource.TestCheckResourceAttr("data.turso_database_instance.test", "region", "aws-us-east-1"),
					resource.TestCheckResourceAttr("data.turso_database_instance.test", "type", "primary"),
					resource.TestCheckResourceAttrSet("data.turso_database_instance.test", "uuid"),
				),
			},
		},
	})
}

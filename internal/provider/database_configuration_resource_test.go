package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "turso_database_configuration" "test" {
					organization_slug = "jpedroh"
					database_name	  = "tfproviderdatasource"
					size_limit	  	  = "1gb"
					block_reads	      = true
					block_writes	  = true
					delete_protection = false
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("turso_database_configuration.test", "organization_slug", "jpedroh"),
					resource.TestCheckResourceAttr("turso_database_configuration.test", "database_name", "tfproviderdatasource"),
					resource.TestCheckResourceAttr("turso_database_configuration.test", "size_limit", "1gb"),
					resource.TestCheckResourceAttr("turso_database_configuration.test", "block_reads", "true"),
					resource.TestCheckResourceAttr("turso_database_configuration.test", "block_writes", "true"),
					resource.TestCheckResourceAttr("turso_database_configuration.test", "delete_protection", "false"),
				),
			},
			{
				ResourceName:                         "turso_database_configuration.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "jpedroh/tfproviderdatasource",
				ImportStateVerifyIdentifierAttribute: "database_name",
			},
		},
	})
}

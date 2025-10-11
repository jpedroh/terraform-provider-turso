package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseConfigurationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "turso_database_configuration" "test" {
					organization_slug = "jpedroh"
					database_name = "tfproviderdatasource"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.turso_database_configuration.test", "organization_slug", "jpedroh"),
					resource.TestCheckResourceAttr("data.turso_database_configuration.test", "database_name", "tfproviderdatasource"),
					resource.TestCheckResourceAttr("data.turso_database_configuration.test", "size_limit", "1gb"),
					resource.TestCheckResourceAttr("data.turso_database_configuration.test", "block_reads", "true"),
					resource.TestCheckResourceAttr("data.turso_database_configuration.test", "block_writes", "true"),
					resource.TestCheckResourceAttr("data.turso_database_configuration.test", "delete_protection", "false"),
				),
			},
		},
	})
}

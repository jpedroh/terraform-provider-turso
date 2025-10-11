package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "turso_database" "tf_provider_data_source" {
					name = "tfproviderdatasource"
					organization_name = "jpedroh"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.turso_database.tf_provider_data_source", "name", "tfproviderdatasource"),
					resource.TestCheckResourceAttr("data.turso_database.tf_provider_data_source", "organization_name", "jpedroh"),
					resource.TestCheckResourceAttr("data.turso_database.tf_provider_data_source", "hostname", "tfproviderdatasource-jpedroh.aws-us-east-1.turso.io"),
					resource.TestCheckResourceAttrSet("data.turso_database.tf_provider_data_source", "db_id"),
				),
			},
		},
	})
}

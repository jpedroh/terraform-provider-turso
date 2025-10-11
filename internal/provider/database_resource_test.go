package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "turso_database" "test" {
					organization_name = "jpedroh"
					name	  = "tf-provider-resource"
					group	  = "default"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("turso_database.test", "organization_name", "jpedroh"),
					resource.TestCheckResourceAttr("turso_database.test", "name", "tf-provider-resource"),
					resource.TestCheckResourceAttr("turso_database.test", "group", "default"),
					resource.TestCheckResourceAttr("turso_database.test", "hostname", "tf-provider-resource-jpedroh.aws-us-east-1.turso.io"),
					resource.TestCheckResourceAttrSet("turso_database.test", "db_id"),
				),
			},
			{
				ResourceName:      "turso_database.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "jpedroh/tf-provider-resource",
				ImportStateVerifyIdentifierAttribute: "db_id",
			},
		},
	})
}

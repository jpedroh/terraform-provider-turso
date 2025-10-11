package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "turso_organization" "test" { slug = "jpedroh" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.turso_organization.test", "slug", "jpedroh"),
					resource.TestCheckResourceAttr("data.turso_organization.test", "name", "jpedroh"),
					resource.TestCheckResourceAttr("data.turso_organization.test", "type", "personal"),
				),
			},
		},
	})
}

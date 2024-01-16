package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSpaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `resource "netdata_space" "test" { name = "testAcc" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netdata_space.test", "name", "testAcc"),
					resource.TestCheckResourceAttr("netdata_space.test", "description", ""),
				),
			},
			// Update and Read testing
			{
				Config: `
					resource "netdata_space" "test" {
						name = "testAccUpdated"
						description ="testDesc"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netdata_space.test", "name", "testAccUpdated"),
					resource.TestCheckResourceAttr("netdata_space.test", "description", "testDesc"),
				),
			},
		},
	})
}

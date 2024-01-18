package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSpaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "netdata_space" "test" {
						name = "testAcc"
					}
					data "netdata_space" "test" {
						id = netdata_space.test.id
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netdata_space.test", "name", "testAcc"),
					resource.TestMatchResourceAttr("data.netdata_space.test", "claim_token", regexp.MustCompile(`^.{135}$`)),
				),
			},
		},
	},
	)
}

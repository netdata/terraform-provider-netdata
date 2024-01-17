package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoomDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "netdata_space" "test" {
						name = "testAcc"
					}
					resource "netdata_room" "test" {
						spaceid = netdata_space.test.id
						name    = "testAcc"
					}
					data "netdata_room" "test" {
						spaceid = netdata_space.test.id
						id      = netdata_room.test.id
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netdata_room.test", "name", "testAcc"),
				),
			},
		},
	},
	)
}

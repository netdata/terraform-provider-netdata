package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoomResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netdata_space" "test" {
					name = "testAcc"
				}
				resource "netdata_room" "test" {
					space_id = netdata_space.test.id
					name    = "testAcc"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netdata_room.test", "name", "testAcc"),
					resource.TestCheckResourceAttr("netdata_room.test", "description", ""),
				),
			},
		},
	})
}

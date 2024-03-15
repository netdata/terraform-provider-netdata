package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoomMemberResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "netdata_space_member" "test" {
					email    = "foo@bar.local"
					space_id = "%s"
					role     = "admin"
				}
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_room_member" "test" {
					room_id         = netdata_room.test.id
					space_id        = "%s"
					space_member_id = netdata_space_member.test.id
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_room_member.test", "room_id"),
					resource.TestCheckResourceAttrSet("netdata_room_member.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_room_member.test", "space_member_id"),
				),
			},
		},
	})
}

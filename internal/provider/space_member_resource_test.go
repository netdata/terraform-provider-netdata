package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSpaceMemberResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "netdata_space_member" "test" {
					email    = "space@member.local"
					space_id = "%s"
					role     = "admin"
				}
				`, getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_space_member.test", "id"),
					resource.TestCheckResourceAttr("netdata_space_member.test", "email", "space@member.local"),
					resource.TestCheckResourceAttr("netdata_space_member.test", "role", "admin"),
					resource.TestCheckResourceAttrSet("netdata_space_member.test", "space_id"),
				),
			},
		},
	})
}

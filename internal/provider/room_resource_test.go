package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoomResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name    = "testAcc"
				}
				`, os.Getenv("SPACE_ID_NON_COMMUNITY")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netdata_room.test", "name", "testAcc"),
					resource.TestCheckResourceAttr("netdata_room.test", "description", ""),
				),
			},
		},
	})
}

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSpaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "netdata_space" "test" {
						id = "%s"
					}
					`, getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netdata_space.test", "name"),
					resource.TestMatchResourceAttr("data.netdata_space.test", "claim_token", regexp.MustCompile(`^.{135}$`)),
				),
			},
		},
	},
	)
}

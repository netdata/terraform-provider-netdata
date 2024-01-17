package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSpaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "netdata_space" "test" { name = "testAcc" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netdata_space.test", "name", "testAcc"),
					resource.TestCheckResourceAttr("netdata_space.test", "description", ""),
					resource.TestMatchResourceAttr("netdata_space.test", "claimtoken", regexp.MustCompile(`^.{135}$`)),
				),
			},
		},
	})
}

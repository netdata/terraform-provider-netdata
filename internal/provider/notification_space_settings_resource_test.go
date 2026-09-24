package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationSpaceSettingsResource(t *testing.T) {
	spaceID := getNonCommunitySpaceIDEnv()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Reject a delay below the API minimum
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_space_settings" "test" {
					space_id           = "%s"
					reachability_delay = 29
				}
				`, spaceID),
				ExpectError: regexp.MustCompile(`must be at least 30`),
			},
			// Create
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_space_settings" "test" {
					space_id           = "%s"
					reachability_delay = 120
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_space_settings.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_space_settings.test", "space_id", spaceID),
					resource.TestCheckResourceAttr("netdata_notification_space_settings.test", "reachability_delay", "120"),
				),
			},
			// Update the delay
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_space_settings" "test" {
					space_id           = "%s"
					reachability_delay = 660
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_space_settings.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_space_settings.test", "reachability_delay", "660"),
				),
			},
		},
	})
}

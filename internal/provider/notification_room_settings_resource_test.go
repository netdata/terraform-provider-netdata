package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationRoomSettingsResource(t *testing.T) {
	spaceID := getNonCommunitySpaceIDEnv()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Reject a delay below the API minimum
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc-notification-settings"
				}
				resource "netdata_notification_room_settings" "test" {
					space_id           = "%s"
					room_id            = netdata_room.test.id
					reachability_delay = 29
				}
				`, spaceID, spaceID),
				ExpectError: regexp.MustCompile(`must be at least 30`),
			},
			// Create
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc-notification-settings"
				}
				resource "netdata_notification_room_settings" "test" {
					space_id           = "%s"
					room_id            = netdata_room.test.id
					reachability_delay = 180
				}
				`, spaceID, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_room_settings.test", "id"),
					resource.TestCheckResourceAttrSet("netdata_notification_room_settings.test", "room_id"),
					resource.TestCheckResourceAttr("netdata_notification_room_settings.test", "space_id", spaceID),
					resource.TestCheckResourceAttr("netdata_notification_room_settings.test", "reachability_delay", "180"),
				),
			},
			// Update the delay
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc-notification-settings"
				}
				resource "netdata_notification_room_settings" "test" {
					space_id           = "%s"
					room_id            = netdata_room.test.id
					reachability_delay = 840
				}
				`, spaceID, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_room_settings.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_room_settings.test", "reachability_delay", "840"),
				),
			},
		},
	})
}

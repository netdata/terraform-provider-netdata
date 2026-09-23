package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMicrosoftTeamsNotificationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_notification_microsoftteams_channel" "test" {
					name                    = "microsoftteams"
					enabled                 = true
					space_id                = "%s"
					rooms_id                = [netdata_room.test.id]
					notifications           = ["CRITICAL","WARNING","CLEAR"]
					webhook_url             = "https://xxxxxxx.webhook.office.com/webhookb2/00000000-0000-0000-0000-000000000000@00000000-0000-0000-0000-000000000000/IncomingWebhook/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/00000000-0000-0000-0000-000000000000"
					repeat_notification_min = 30
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_microsoftteams_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "name", "microsoftteams"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("netdata_notification_microsoftteams_channel.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_microsoftteams_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "notifications.0", "CRITICAL"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "notifications.1", "WARNING"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "notifications.2", "CLEAR"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "repeat_notification_min", "30"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "webhook_url", "https://xxxxxxx.webhook.office.com/webhookb2/00000000-0000-0000-0000-000000000000@00000000-0000-0000-0000-000000000000/IncomingWebhook/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_notification_microsoftteams_channel" "test" {
					name                    = "microsoftteams"
					enabled                 = true
					space_id                = "%s"
					rooms_id                = [netdata_room.test.id]
					notifications           = ["CLEAR","CRITICAL"]
					webhook_url             = "https://xxxxxxx.webhook.office.com/webhookb2/00000000-0000-0000-0000-000000000000@00000000-0000-0000-0000-000000000000/IncomingWebhook/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/00000000-0000-0000-0000-000000000000"
					repeat_notification_min = 60
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_microsoftteams_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "name", "microsoftteams"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("netdata_notification_microsoftteams_channel.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_microsoftteams_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "notifications.0", "CLEAR"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "notifications.1", "CRITICAL"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "repeat_notification_min", "60"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "webhook_url", "https://xxxxxxx.webhook.office.com/webhookb2/00000000-0000-0000-0000-000000000000@00000000-0000-0000-0000-000000000000/IncomingWebhook/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_notification_microsoftteams_channel" "test" {
					name             = "microsoftteams"
					enabled          = false
					space_id         = "%s"
					rooms_id         = null
					notifications    = ["CRITICAL","WARNING","CLEAR"]
					webhook_url      = "https://xxxxxxx.webhook.office.com/webhookb2/00000000-0000-0000-0000-000000000000@00000000-0000-0000-0000-000000000000/IncomingWebhook/YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY/00000000-0000-0000-0000-000000000000"
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_microsoftteams_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "name", "microsoftteams"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("netdata_notification_microsoftteams_channel.test", "space_id"),
					resource.TestCheckNoResourceAttr("netdata_notification_microsoftteams_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "notifications.0", "CRITICAL"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "notifications.1", "WARNING"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "notifications.2", "CLEAR"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "repeat_notification_min", "0"),
					resource.TestCheckResourceAttr("netdata_notification_microsoftteams_channel.test", "webhook_url", "https://xxxxxxx.webhook.office.com/webhookb2/00000000-0000-0000-0000-000000000000@00000000-0000-0000-0000-000000000000/IncomingWebhook/YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY/00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

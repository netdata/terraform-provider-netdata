package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSlackNotificationResource(t *testing.T) {
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
				resource "netdata_notification_slack_channel" "test" {
				  	name        = "slack"
				  	enabled     = true
				  	space_id    = "%s"
				  	rooms_id    = [netdata_room.test.id]
				  	webhook_url = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
				  	alarms      = "ALARMS_SETTING_ALL"
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_slack_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_slack_channel.test", "name", "slack"),
					resource.TestCheckResourceAttr("netdata_notification_slack_channel.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("netdata_notification_slack_channel.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_slack_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_slack_channel.test", "alarms", "ALARMS_SETTING_ALL"),
					resource.TestCheckResourceAttr("netdata_notification_slack_channel.test", "webhook_url", "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_notification_slack_channel" "test" {
				  	name        = "slack"
				  	enabled     = false
				  	space_id    = "%s"
				  	rooms_id    = null
				  	webhook_url = "https://hooks.slack.com/services/T10000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
				  	alarms      = "ALARMS_SETTING_ALL"
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_slack_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_slack_channel.test", "name", "slack"),
					resource.TestCheckResourceAttr("netdata_notification_slack_channel.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("netdata_notification_slack_channel.test", "space_id"),
					resource.TestCheckNoResourceAttr("netdata_notification_slack_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_slack_channel.test", "alarms", "ALARMS_SETTING_ALL"),
					resource.TestCheckResourceAttr("netdata_notification_slack_channel.test", "webhook_url", "https://hooks.slack.com/services/T10000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

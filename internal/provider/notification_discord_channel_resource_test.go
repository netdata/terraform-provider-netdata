package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDiscordNotificationResource(t *testing.T) {
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
				resource "netdata_notification_discord_channel" "test" {
				  	name         = "discord"
				  	enabled      = true
				  	space_id     = "%s"
				  	rooms_id     = [netdata_room.test.id]
				  	webhook_url  = "https://discord.com/api/webhooks/0000000000000/XXXXXXXXXXXXXXXXXXXXXXXX"
				  	alarms       = "ALARMS_SETTING_ALL"
					channel_type = "text"
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_discord_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "name", "discord"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("netdata_notification_discord_channel.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_discord_channel.test", "integration_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_discord_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "alarms", "ALARMS_SETTING_ALL"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "webhook_url", "https://discord.com/api/webhooks/0000000000000/XXXXXXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "channel_type", "text"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_notification_discord_channel" "test" {
				  	name           = "discord"
				  	enabled        = false
				  	space_id       = "%s"
				  	rooms_id       = null
				  	webhook_url    = "https://discord.com/api/webhooks/1000000000000/XXXXXXXXXXXXXXXXXXXXXXXX"
				  	alarms         = "ALARMS_SETTING_ALL"
					channel_type   = "forum"
					channel_thread = "thread"
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_discord_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "name", "discord"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("netdata_notification_discord_channel.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_discord_channel.test", "integration_id"),
					resource.TestCheckNoResourceAttr("netdata_notification_discord_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "alarms", "ALARMS_SETTING_ALL"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "webhook_url", "https://discord.com/api/webhooks/1000000000000/XXXXXXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "channel_type", "forum"),
					resource.TestCheckResourceAttr("netdata_notification_discord_channel.test", "channel_thread", "thread"),
				),
			},
		},
	})
}

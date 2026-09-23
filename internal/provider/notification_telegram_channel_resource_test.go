package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTelegramNotificationResource(t *testing.T) {
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
				resource "netdata_notification_telegram_channel" "test" {
					name                    = "telegram"
					enabled                 = true
					space_id                = "%s"
					rooms_id                = [netdata_room.test.id]
					notifications           = ["CRITICAL","WARNING","CLEAR"]
					bot_token               = "0000000000:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
					chat_id                 = "-1000000000000"
					topic_id                = 10
					repeat_notification_min = 30
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_telegram_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "name", "telegram"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("netdata_notification_telegram_channel.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_telegram_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "notifications.0", "CRITICAL"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "notifications.1", "WARNING"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "notifications.2", "CLEAR"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "repeat_notification_min", "30"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "bot_token", "0000000000:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "chat_id", "-1000000000000"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "topic_id", "10"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_notification_telegram_channel" "test" {
					name                    = "telegram"
					enabled                 = true
					space_id                = "%s"
					rooms_id                = [netdata_room.test.id]
					notifications           = ["CLEAR","CRITICAL"]
					bot_token               = "0000000000:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
					chat_id                 = "-1000000000000"
					topic_id                = 10
					repeat_notification_min = 60
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_telegram_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "name", "telegram"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("netdata_notification_telegram_channel.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_telegram_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "notifications.0", "CLEAR"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "notifications.1", "CRITICAL"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "repeat_notification_min", "60"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "bot_token", "0000000000:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "chat_id", "-1000000000000"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "topic_id", "10"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_notification_telegram_channel" "test" {
					name             = "telegram"
					enabled          = false
					space_id         = "%s"
					rooms_id         = null
					notifications    = ["CRITICAL","WARNING","CLEAR"]
					bot_token        = "1111111111:YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"
					chat_id          = "-1000000000001"
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_telegram_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "name", "telegram"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("netdata_notification_telegram_channel.test", "space_id"),
					resource.TestCheckNoResourceAttr("netdata_notification_telegram_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "notifications.0", "CRITICAL"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "notifications.1", "WARNING"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "notifications.2", "CLEAR"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "repeat_notification_min", "0"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "bot_token", "1111111111:YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "chat_id", "-1000000000001"),
					resource.TestCheckResourceAttr("netdata_notification_telegram_channel.test", "topic_id", "0"),
				),
			},
		},
	})
}

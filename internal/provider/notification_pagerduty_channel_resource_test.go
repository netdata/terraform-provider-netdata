package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerdutyNotificationResource(t *testing.T) {
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
				resource "netdata_notification_pagerduty_channel" "test" {
				  	name        	 = "pagerduty"
				  	enabled     	 = true
				  	space_id    	 = "%s"
				  	rooms_id    	 = [netdata_room.test.id]
				  	alarms      	 = "ALARMS_SETTING_ALL"
					alert_events_url = "https://events.pagerduty.com/v2/enqueue"
  					integration_key  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_pagerduty_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "name", "pagerduty"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("netdata_notification_pagerduty_channel.test", "space_id"),
					resource.TestCheckResourceAttrSet("netdata_notification_pagerduty_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "alarms", "ALARMS_SETTING_ALL"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "alert_events_url", "https://events.pagerduty.com/v2/enqueue"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "integration_key", "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "netdata_room" "test" {
					space_id = "%s"
					name     = "testAcc"
				}
				resource "netdata_notification_pagerduty_channel" "test" {
				  	name        	 = "pagerduty"
				  	enabled     	 = false
				  	space_id    	 = "%s"
				  	rooms_id    	 = null
				  	alarms      	 = "ALARMS_SETTING_ALL"
					alert_events_url = "https://events.pagerduty.com/v2/enqueue"
  					integration_key  = "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"
				}
				`, getNonCommunitySpaceIDEnv(), getNonCommunitySpaceIDEnv()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_pagerduty_channel.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "name", "pagerduty"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("netdata_notification_pagerduty_channel.test", "space_id"),
					resource.TestCheckNoResourceAttr("netdata_notification_pagerduty_channel.test", "rooms_id.0"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "alarms", "ALARMS_SETTING_ALL"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "alert_events_url", "https://events.pagerduty.com/v2/enqueue"),
					resource.TestCheckResourceAttr("netdata_notification_pagerduty_channel.test", "integration_key", "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"),
				),
			},
		},
	})
}

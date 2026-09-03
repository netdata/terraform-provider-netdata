package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNotificationSilencingRuleResource(t *testing.T) {
	spaceID := getNonCommunitySpaceIDEnv()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create with minimal config
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_silencing_rule" "test" {
					space_id  = "%s"
					name      = "testAcc-minimal"
					starts_at = "2026-09-01T00:00:00Z"
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_silencing_rule.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "name", "testAcc-minimal"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "space_id", spaceID),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "starts_at", "2026-09-01T00:00:00Z"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "disabled", "false"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "delete_on_expiry", "false"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "timezone", "UTC"),
				),
			},
			// Update name and add optional fields
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_silencing_rule" "test" {
					space_id    = "%s"
					name        = "testAcc-updated"
					starts_at   = "2026-09-01T00:00:00Z"
					lasts_until = "2026-12-31T23:59:59Z"
					disabled    = true
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_silencing_rule.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "name", "testAcc-updated"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "starts_at", "2026-09-01T00:00:00Z"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "lasts_until", "2026-12-31T23:59:59Z"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "disabled", "true"),
				),
			},
			// Import
			{
				ResourceName: "netdata_notification_silencing_rule.test",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources["netdata_notification_silencing_rule.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return fmt.Sprintf("%s,%s", rs.Primary.Attributes["space_id"], rs.Primary.ID), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotificationSilencingRuleResource_WithNotificationOptions(t *testing.T) {
	spaceID := getNonCommunitySpaceIDEnv()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_silencing_rule" "test" {
					space_id             = "%s"
					name                 = "testAcc-options"
					starts_at            = "2026-09-01T00:00:00Z"
					lasts_until          = "2026-12-31T23:59:59Z"
					notification_options = ["CRITICAL", "WARNING"]
					alert_names          = ["disk.space"]
					alert_contexts       = ["disk_space_usage"]
					timezone             = "America/New_York"
					delete_on_expiry     = true
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_silencing_rule.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "notification_options.0", "CRITICAL"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "notification_options.1", "WARNING"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "alert_names.0", "disk.space"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "alert_contexts.0", "disk_space_usage"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "timezone", "America/New_York"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "delete_on_expiry", "true"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "starts_at", "2026-09-01T00:00:00Z"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "lasts_until", "2026-12-31T23:59:59Z"),
				),
			},
			// Update notification options
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_silencing_rule" "test" {
					space_id             = "%s"
					name                 = "testAcc-options"
					starts_at            = "2026-09-01T00:00:00Z"
					notification_options = ["CLEAR", "CRITICAL", "WARNING"]
					alert_names          = ["disk.space", "cpu.usage"]
					timezone             = "America/New_York"
					delete_on_expiry     = false
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "notification_options.0", "CLEAR"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "notification_options.1", "CRITICAL"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "notification_options.2", "WARNING"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "alert_names.0", "disk.space"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "alert_names.1", "cpu.usage"),
					resource.TestCheckNoResourceAttr("netdata_notification_silencing_rule.test", "alert_contexts.0"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "delete_on_expiry", "false"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "starts_at", "2026-09-01T00:00:00Z"),
				),
			},
		},
	})
}

func TestAccNotificationSilencingRuleResource_WithRRule(t *testing.T) {
	spaceID := getNonCommunitySpaceIDEnv()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create with rrule
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_silencing_rule" "test" {
					space_id  = "%s"
					name      = "testAcc-rrule"
					starts_at = "2026-09-01T07:05:00Z"
					rrule     = "RRULE:FREQ=MONTHLY;INTERVAL=1;COUNT=10;BYMONTHDAY=1"
					timezone  = "UTC"
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_silencing_rule.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "name", "testAcc-rrule"),
					// Asserts the exact value, not just presence: the API returns the
					// rrule joined with a DTSTART line that stripDTSTART has to remove
					// to reproduce the configured string.
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "rrule", "RRULE:FREQ=MONTHLY;INTERVAL=1;COUNT=10;BYMONTHDAY=1"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "timezone", "UTC"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "delete_on_expiry", "false"),
				),
			},
			// Update rrule
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_silencing_rule" "test" {
					space_id  = "%s"
					name      = "testAcc-rrule-updated"
					starts_at = "2026-09-01T07:05:00Z"
					rrule     = "RRULE:FREQ=WEEKLY;INTERVAL=1;COUNT=5;BYDAY=MO"
					timezone  = "America/New_York"
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "name", "testAcc-rrule-updated"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "rrule", "RRULE:FREQ=WEEKLY;INTERVAL=1;COUNT=5;BYDAY=MO"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "timezone", "America/New_York"),
				),
			},
			// Import
			{
				ResourceName: "netdata_notification_silencing_rule.test",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources["netdata_notification_silencing_rule.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return fmt.Sprintf("%s,%s", rs.Primary.Attributes["space_id"], rs.Primary.ID), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotificationSilencingRuleResource_WithHostLabels(t *testing.T) {
	spaceID := getNonCommunitySpaceIDEnv()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "netdata_notification_silencing_rule" "test" {
					space_id    = "%s"
					name        = "testAcc-labels"
					starts_at   = "2026-09-01T00:00:00Z"
					host_labels = {
						os  = "linux"
						env = "production"
					}
				}
				`, spaceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_notification_silencing_rule.test", "id"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "host_labels.os", "linux"),
					resource.TestCheckResourceAttr("netdata_notification_silencing_rule.test", "host_labels.env", "production"),
				),
			},
		},
	})
}

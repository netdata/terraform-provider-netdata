terraform {
  required_providers {
    netdata = {
      source = "netdata/netdata"
    }
  }
  required_version = ">= 1.1.0"
}

provider "netdata" {}

resource "netdata_node_room_member" "new" {
  room_id  = netdata_room.test.id
  space_id = netdata_space.test.id

  node_names = [
    "node1",
    "node2"
  ]

  rule {
    action      = "INCLUDE"
    description = "Description of the rule"
    clause {
      label    = "role"
      operator = "equals"
      value    = "parent"
      negate   = false
    }
    clause {
      label    = "environment"
      operator = "equals"
      value    = "production"
      negate   = false
    }
  }
  rule {
    action      = "EXCLUDE"
    description = "Description of the rule"
    clause {
      label    = "role"
      operator = "equals"
      value    = "parent"
      negate   = true
    }
    clause {
      label    = "environment"
      operator = "contains"
      value    = "production"
      negate   = false
    }
  }
}

resource "netdata_space" "test" {
  name        = "MyTestingSpace"
  description = "Created by Terraform"
}

resource "netdata_room" "test" {
  space_id    = netdata_space.test.id
  name        = "MyTestingRoom"
  description = "Created by Terraform"
}

resource "netdata_space_member" "test" {
  email    = "foo@bar.local"
  space_id = netdata_space.test.id
  role     = "admin"
}

resource "netdata_room_member" "test" {
  room_id         = netdata_room.test.id
  space_id        = netdata_space.test.id
  space_member_id = netdata_space_member.test.id
}

resource "netdata_node_room_member" "test" {
  room_id  = netdata_room.test.id
  space_id = netdata_space.test.id
  node_names = [
    "node1",
    "node2"
  ]
}

resource "netdata_notification_slack_channel" "test" {
  name = "slack"

  enabled                 = true
  space_id                = netdata_space.test.id
  rooms_id                = [netdata_room.test.id]
  notifications           = ["CRITICAL", "WARNING", "CLEAR"]
  repeat_notification_min = 60
  webhook_url             = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
}

resource "netdata_notification_discord_channel" "test" {
  name = "discord"

  enabled        = true
  space_id       = netdata_space.test.id
  rooms_id       = [netdata_room.test.id]
  notifications  = ["CRITICAL", "WARNING", "CLEAR"]
  webhook_url    = "https://discord.com/api/webhooks/0000000000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  channel_type   = "forum"
  channel_thread = "thread"
}

resource "netdata_notification_pagerduty_channel" "test" {
  name = "pagerduty"

  enabled          = true
  space_id         = netdata_space.test.id
  notifications    = ["CRITICAL", "WARNING", "CLEAR"]
  alert_events_url = "https://events.pagerduty.com/v2/enqueue"
  integration_key  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
}

resource "netdata_notification_silencing_rule" "test" {
  space_id             = netdata_space.test.id
  name                 = "Testing Rule"
  starts_at            = "2026-09-01T00:00:00Z"
  lasts_until          = "2030-12-31T23:59:59Z"
  notification_options = ["CRITICAL", "WARNING"]
  alert_names          = ["disk.space"]
  alert_contexts       = ["disk_space_usage"]
  timezone             = "America/New_York"
  rrule                = "RRULE:FREQ=WEEKLY;INTERVAL=1;COUNT=5;BYDAY=MO"
  delete_on_expiry     = false
}

data "netdata_space" "test" {
  id = netdata_space.test.id
}

data "netdata_room" "test" {
  id       = netdata_room.test.id
  space_id = netdata_space.test.id
}

output "space_name" {
  value = data.netdata_space.test.name
}

output "claim_token" {
  value = netdata_space.test.claim_token
}

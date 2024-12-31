resource "netdata_notification_pagerduty_channel" "test" {
  name = "pagerduty notifications"

  enabled                 = true
  space_id                = netdata_space.test.id
  rooms_id                = ["<room_id>"]
  alarms                  = "ALARMS_SETTING_ALL"
  repeat_notification_min = 30
  alert_events_url        = "https://events.pagerduty.com/v2/enqueue"
  integration_key         = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
}

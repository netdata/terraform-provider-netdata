resource "netdata_notification_slack_channel" "test" {
  name = "slack notifications"

  enabled                 = true
  space_id                = "<space_id>"
  rooms_id                = ["<room_id>"]
  repeat_notification_min = 30
  webhook_url             = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  alarms                  = "ALARMS_SETTING_ALL"
}

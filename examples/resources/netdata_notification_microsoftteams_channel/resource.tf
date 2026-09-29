resource "netdata_notification_microsoftteams_channel" "test" {
  name = "microsoft teams notifications"

  enabled                 = true
  space_id                = "<space_id>"
  rooms_id                = ["<room_id>"]
  notifications           = ["CRITICAL", "WARNING", "CLEAR"]
  repeat_notification_min = 30
  webhook_url             = "<incoming_webhook_url>"
}

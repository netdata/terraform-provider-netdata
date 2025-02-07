resource "netdata_notification_discord_channel" "test" {
  name = "discord notifications"

  enabled                 = true
  space_id                = "<space_id>"
  rooms_id                = ["<room_id>"]
  repeat_notification_min = 30
  webhook_url             = "https://discord.com/api/webhooks/0000000000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  notifications           = ["CRITICAL", "WARNING", "CLEAR"]
  channel_type            = "text"
}

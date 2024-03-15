resource "netdata_notification_discord_channel" "test" {
  name = "discord notifications"

  enabled      = true
  space_id     = "<space_id>"
  rooms_id     = ["<room_id>"]
  webhook_url  = "https://discord.com/api/webhooks/0000000000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  alarms       = "ALARMS_SETTING_ALL"
  channel_type = "text"
}

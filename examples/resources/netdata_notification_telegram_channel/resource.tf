resource "netdata_notification_telegram_channel" "test" {
  name = "telegram notifications"

  enabled                 = true
  space_id                = "<space_id>"
  rooms_id                = ["<room_id>"]
  notifications           = ["CRITICAL", "WARNING", "CLEAR"]
  repeat_notification_min = 30
  bot_token               = "<bot_token>"
  chat_id                 = "<chat_id>"
  topic_id                = 0
}

resource "netdata_notification_room_settings" "test" {
  space_id           = "<space_id>"
  room_id            = "<room_id>"
  reachability_delay = 180
}

resource "netdata_notification_silencing_rule" "test" {
  space_id             = "<space_id>"
  name                 = "testing"
  starts_at            = "2025-09-01T00:00:00Z"
  lasts_until          = "2025-12-31T23:59:59Z"
  notification_options = ["CRITICAL", "WARNING"]
  alert_names          = ["disk.space"]
  alert_contexts       = ["disk_space_usage"]
  timezone             = "America/New_York"
  rrule                = "RRULE:FREQ=WEEKLY;INTERVAL=1;COUNT=5;BYDAY=MO"
  delete_on_expiry     = false
}

## 0.3.0

FEATURES:

- add Node Rule-Based Room Assignment option to the `netdata_node_room_member` resource

## 0.2.2

FEATURES:

- option to specify notification repeat interval `repeat_notification_min` for the paid notification channels

## 0.2.1

BUGFIXES:

- the `integration_id` attribute is being removed because it is internally used for the create resource only and
  doesn't bring much value to store it

## 0.2.0

FEATURES:

- add `netdata_node_room_member` resource

## 0.1.3

FEATURES:

- more detailed resource descriptions

## 0.1.2

BUGFIXES:

- fix bug with empty claim token

## 0.1.1

BUGFIXES:

- empty claim token when importing space

## 0.1.0

Initial version.

FEATURES:

- support for the Netdata Cloud Spaces
- support for the Netdata Cloud Rooms
- support for the Netdata Space and Room Membership
- support for the Netdata Notifications: Discord, Slack, Pagerduty

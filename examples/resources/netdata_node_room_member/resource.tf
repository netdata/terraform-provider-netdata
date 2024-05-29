resource "netdata_node_room_member" "test" {
  space_id = "<space_id>"
  room_id  = "<room_id>"
  node_names = [
    "node1",
    "node2"
  ]
}

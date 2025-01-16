resource "netdata_node_room_member" "test" {
  space_id = "<space_id>"
  room_id  = "<room_id>"
  node_names = [
    "node1",
    "node2"
  ]
  rule {
    action      = "INCLUDE"
    description = "Description of the rule"
    clause {
      label    = "role"
      operator = "equals"
      value    = "parent"
      negate   = false
    }
    clause {
      label    = "environment"
      operator = "equals"
      value    = "production"
      negate   = false
    }
  }
}

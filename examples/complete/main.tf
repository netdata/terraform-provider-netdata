terraform {
  required_providers {
    netdata = {
      # TODO: Update this string with the published name of your provider.
      source = "netdata.cloud/todo/netdata"
    }
  }
  required_version = ">= 1.1.0"
}

provider "netdata" {}

resource "netdata_space" "test" {
  name        = "MyTestingSpace"
  description = "Created by Terraform"
}

resource "netdata_room" "test" {
  space_id    = netdata_space.test.id
  name        = "MyTestingRoom"
  description = "Created by Terraform2"
}

data "netdata_space" "test" {
  id = netdata_space.test.id
}

data "netdata_room" "test" {
  id       = netdata_room.test.id
  space_id = netdata_space.test.id
}

output "datasource" {
  value = data.netdata_space.test.name
}

output "claim_token" {
  value = netdata_space.test.claim_token
}

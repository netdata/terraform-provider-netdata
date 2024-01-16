terraform {
  required_providers {
    netdata = {
      source = "netdata.cloud/todo/netdata"
    }
  }
  required_version = ">= 1.1.0"
}

provider "netdata" {
  url       = "https://app.netdata.cloud"
  authtoken = "authtoken"
}

resource "netdata_space" "test" {
  name        = "MyTestingSpace"
  description = "Created by Terraform"
}

data "netdata_space" "test" {
  id = "ee3ec76d-0180-4ef4-93ae-c94c1e7ed2f1"
}

output "datasource" {
  value = data.netdata_space.test.name
}

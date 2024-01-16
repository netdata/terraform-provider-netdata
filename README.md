# Terraform Provider for Netdata Cloud

This is the [Terraform](https://www.terraform.io/) provider for the [Netdata Cloud](https://www.netdata.cloud/).

This provider allows you to install and manage Netdata Cloud resources using Terraform.


## Contents

* [Requirements](#requirements)
* [Getting Started](#getting-started)
* [Example](#example)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) v1.1.0 or later
- [Go](https://golang.org/doc/install) v1.20 or later (to build the provider plugin)

## Getting Started

* from terraform registry

TODO

* from source code

	* setup your [CLI configuration](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)

	```console
	$ cat ~/.terraformrc
	provider_installation {
  		dev_overrides {
		    # TODO: Update this string with the published name of your provider
  		    "netdata.cloud/todo/netdata" = "<your GOBIN directory>"
  		}
  		direct {}
	}
	```

	* build the provider

	```console
	$ make local-build
	```

## Example

```hcl
terraform {
  required_providers {
    netdata = {
      # TODO: Update this string with the published name of your provider.
      source = "netdata.cloud/todo/netdata"
    }
  }
  required_version = ">= 1.1.0"
}

provider "netdata" {
  url       = "https://app.netdata.cloud"
  authtoken = "<authtoken>"
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

```

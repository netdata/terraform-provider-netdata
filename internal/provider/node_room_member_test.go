package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNodeRoomMemberResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "netdata_space" "test" {
					name = "TestAccSpace"
				}
				resource "netdata_room" "test" {
					space_id    = netdata_space.test.id
					name        = "TestAccRoom"
				}
				resource "netdata_node_room_member" "test" {
					room_id  = netdata_room.test.id
					space_id = netdata_space.test.id
					node_names = [
					  "netdata-agent"
					]
  					rule {
  					  action      = "INCLUDE"
  					  description = "Description"
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
  					    negate   = true
  					  }
  					}
					depends_on = [
					  terraform_data.install_agent
					]
				}
				resource "terraform_data" "install_agent" {
					provisioner "local-exec" {
					  command = <<EOT
cat > docker-compose.yml <<EOF
services:
  netdata:
    image: netdata/netdata:stable
    container_name: netdata-agent
    restart: unless-stopped
    hostname: "netdata-agent"
    cap_add:
      - SYS_PTRACE
      - SYS_ADMIN
    security_opt:
      - apparmor:unconfined
    volumes:
      - /etc/passwd:/host/etc/passwd:ro
      - /etc/group:/host/etc/group:ro
      - /etc/localtime:/etc/localtime:ro
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /etc/os-release:/host/etc/os-release:ro
      - /var/log:/host/var/log:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./parent-stream.conf:/etc/netdata/stream.conf
    environment:
      - NETDATA_CLAIM_TOKEN=$${NETDATA_CLAIM_TOKEN}
      - NETDATA_CLAIM_URL=%s
EOF
docker compose up -d && sleep 30
EOT
					  environment = {
					    NETDATA_CLAIM_TOKEN = netdata_space.test.claim_token
					  }
					}
					provisioner "local-exec" {
					  when    = destroy
					  command = "docker compose down"
					}
				}
				`, getNetdataCloudURL()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netdata_node_room_member.test", "node_names.0"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.action", "INCLUDE"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.description", "Description"),
					resource.TestCheckResourceAttrSet("netdata_node_room_member.test", "rule.0.id"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.clause.0.label", "role"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.clause.0.operator", "equals"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.clause.0.value", "parent"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.clause.0.negate", "false"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.clause.1.label", "environment"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.clause.1.operator", "equals"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.clause.1.value", "production"),
					resource.TestCheckResourceAttr("netdata_node_room_member.test", "rule.0.clause.1.negate", "true"),
				),
			},
		},
	})
}

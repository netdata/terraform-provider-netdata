package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("invitations_sweeper", &resource.Sweeper{
		Name: "invitations sweeper",
		F: func(r string) error {
			url := os.Getenv("NETDATA_CLOUD_URL")
			auth_token := os.Getenv("NETDATA_CLOUD_AUTH_TOKEN")

			if url == "" {
				url = NetdataCloudURL
			}

			if auth_token == "" {
				return fmt.Errorf("auth_token must be set")
			}

			spaceID := getNonCommunitySpaceIDEnv()
			if spaceID == "" {
				return fmt.Errorf("%s must be set", nonCommunitySpaceIDEnv)
			}

			client := client.NewClient(url, auth_token)
			invitations, err := client.GetInvitations(spaceID)
			if err != nil {
				return err
			}

			err = client.DeleteInvitations(spaceID, invitations)
			if err != nil {
				return err
			}

			return nil
		},
	})
}

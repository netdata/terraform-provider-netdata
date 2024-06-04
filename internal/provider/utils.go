package provider

import (
	"os"
	"testing"
)

const (
	nonCommunitySpaceIDEnv = "SPACE_ID_NON_COMMUNITY"
)

func testAccPreCheck(t *testing.T) {
	if getNonCommunitySpaceIDEnv() == "" {
		t.Fatalf("%s must be set", nonCommunitySpaceIDEnv)
	}
}

func getNonCommunitySpaceIDEnv() string {
	return os.Getenv(nonCommunitySpaceIDEnv)
}

func getNetdataCloudURL() string {
	url, ok := os.LookupEnv("NETDATA_CLOUD_URL")
	if !ok {
		return NetdataCloudURL
	}
	return url
}

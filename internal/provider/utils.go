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

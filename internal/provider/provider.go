package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider function for Netdata
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API key for Netdata Cloud authentication",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"netdata_dashboard": resourceDashboard(),
			"netdata_server":    resourceServer(),
		},
	}
}

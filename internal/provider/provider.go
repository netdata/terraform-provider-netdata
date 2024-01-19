package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var _ provider.Provider = &netdataCloudProvider{}

const NetdataCloudURL = "https://app.netdata.cloud"

type netdataCloudProvider struct {
	version string
}

type netdataCloudProviderModel struct {
	Url       types.String `tfsdk:"url"`
	AuthToken types.String `tfsdk:"auth_token"`
}

func (p *netdataCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netdata"
	resp.Version = p.version
}

func (p *netdataCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Netdata Cloud URL Address by default is https://app.netdata.cloud. Can be also set as environment variable `NETDATA_CLOUD_URL`",
				Optional:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "Netdata Cloud Authentication Token. Can be also set as environment variable `NETDATA_CLOUD_AUTH_TOKEN`",
				Sensitive:           true,
				Optional:            true,
			},
		},
	}
}

func (p *netdataCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data netdataCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("NETDATA_CLOUD_URL")
	auth_token := os.Getenv("NETDATA_CLOUD_AUTH_TOKEN")

	if !data.AuthToken.IsNull() {
		auth_token = data.AuthToken.ValueString()
	}

	if !data.Url.IsNull() {
		url = data.Url.ValueString()
	}

	if url == "" {
		url = NetdataCloudURL
	}

	if auth_token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_token"),
			"Missing Netdata Cloud Authentication Token",
			"Provide a valid Netdata Cloud Authentication Token to authenticate with Netdata Cloud.",
		)
	}

	client := client.NewClient(url, auth_token)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *netdataCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSpaceResource,
		NewRoomResource,
	}
}

func (p *netdataCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSpaceDataSource,
		NewRoomDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &netdataCloudProvider{
			version: version,
		}
	}
}

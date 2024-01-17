package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ datasource.DataSource              = &spaceDataSource{}
	_ datasource.DataSourceWithConfigure = &spaceDataSource{}
)

func NewSpaceDataSource() datasource.DataSource {
	return &spaceDataSource{}
}

type spaceDataSource struct {
	client *client.Client
}

type spaceDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ClaimToken  types.String `tfsdk:"claimtoken"`
}

func (s *spaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (s *spaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the space",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the space",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the space",
				Computed:    true,
			},
			"claimtoken": schema.StringAttribute{
				Description: "The claim token of the space",
				Computed:    true,
			},
		},
	}
}

func (s *spaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	s.client = client
}

func (s *spaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state spaceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	tflog.Info(ctx, "Reading Space ID:"+state.ID.ValueString())

	spaceInfo, err := s.client.GetSpaceByID(state.ID.ValueString())

	switch {
	case err == client.ErrNotFound:
		resp.State.RemoveResource(ctx)
		return
	case err != nil:
		resp.Diagnostics.AddError(
			"Error Getting Space",
			"Could Not Read Space ID: "+state.ID.ValueString()+": err: "+err.Error(),
		)
		return
	default:
		state.ID = types.StringValue(spaceInfo.ID)
		state.Name = types.StringValue(spaceInfo.Name)
		state.Description = types.StringValue(spaceInfo.Description)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if state.ClaimToken.IsNull() {
		tflog.Info(ctx, "Creating Claim Token for Space ID: "+state.ID.ValueString())
		claimToken, err := s.client.GetSpaceClaimToken(spaceInfo.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Claim Token",
				"Could Not Create Claim Token for Space ID: "+state.ID.ValueString()+": err: "+err.Error(),
			)
			return
		}
		state.ClaimToken = types.StringValue(*claimToken)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

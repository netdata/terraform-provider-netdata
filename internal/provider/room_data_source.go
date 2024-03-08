package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ datasource.DataSource              = &roomDataSource{}
	_ datasource.DataSourceWithConfigure = &roomDataSource{}
)

func NewRoomDataSource() datasource.DataSource {
	return &roomDataSource{}
}

type roomDataSource struct {
	client *client.Client
}

type roomDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (s *roomDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_room"
}

func (s *roomDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Netdata Cloud Room.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the room",
				Required:    true,
			},
			"space_id": schema.StringAttribute{
				Description: "The ID of the space",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the room",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the room",
				Computed:    true,
			},
		},
	}
}

func (s *roomDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (s *roomDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state roomDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	roomInfo, err := s.client.GetRoomByID(state.ID.ValueString(), state.SpaceID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Room",
			"Could Not Read Room ID: "+state.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(roomInfo.ID)
	state.Name = types.StringValue(roomInfo.Name)
	state.Description = types.StringValue(roomInfo.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

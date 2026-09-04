package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ datasource.DataSource              = &nodeDataSource{}
	_ datasource.DataSourceWithConfigure = &nodeDataSource{}
)

func NewNodeDataSource() datasource.DataSource {
	return &nodeDataSource{}
}

type nodeDataSource struct {
	client *client.Client
}

type nodeDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	SpaceID types.String `tfsdk:"space_id"`
	RoomID  types.String `tfsdk:"room_id"`
}

func (s *nodeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (s *nodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Netdata Node Data Source allows you to retrieve information about a specific node in a Netdata Cloud space.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the node",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the node",
				Required:    true,
			},
			"space_id": schema.StringAttribute{
				Description: "The ID of the space, where the node is placed",
				Required:    true,
				Validators: []validator.String{
					UUID(),
				},
			},
			"room_id": schema.StringAttribute{
				Description: "The ID of the room, where the node is placed",
				Optional:    true,
				Validators: []validator.String{
					UUID(),
				},
			},
		},
	}
}

func (s *nodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (s *nodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		state    nodeDataSourceModel
		nodeInfo *client.RoomNodes
		err      error
	)

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.RoomID.IsNull() {
		nodeInfo, err = s.client.GetAllNodes(state.SpaceID.ValueString())
	} else {
		nodeInfo, err = s.client.GetRoomNodes(state.SpaceID.ValueString(), state.RoomID.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Node Info",
			"Could not get node info, unexpected error: "+err.Error(),
		)
		return
	}

	var matches []string
	for _, node := range nodeInfo.Nodes {
		if node.NodeName == state.Name.ValueString() {
			matches = append(matches, node.NodeID)
		}
	}

	if len(matches) == 0 {
		resp.Diagnostics.AddError(
			"Node Not Found",
			fmt.Sprintf("Could not find node with name: %s", state.Name.ValueString()),
		)
		return
	}

	if len(matches) > 1 {
		resp.Diagnostics.AddError(
			"Multiple Nodes Found",
			fmt.Sprintf("Found %d nodes with name: %s, node names must be unique to use this data source", len(matches), state.Name.ValueString()),
		)
		return
	}

	state.ID = types.StringValue(matches[0])

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

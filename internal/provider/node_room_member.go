package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &roomMemberResource{}
	_ resource.ResourceWithConfigure = &roomMemberResource{}
)

func NewNodeRoomMemberResource() resource.Resource {
	return &nodeRoomMemberResource{}
}

type nodeRoomMemberResource struct {
	client *client.Client
}

type nodeRoomMemberResourceModel struct {
	RoomID    types.String `tfsdk:"room_id"`
	SpaceID   types.String `tfsdk:"space_id"`
	NodeNames types.List   `tfsdk:"node_names"`
}

func (s *nodeRoomMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node_room_member"
}

func (s *nodeRoomMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a Netdata Cloud Node Room Member resource. Use this resource to manage node membership to the room in the selected space, only reachable nodes can be added to the room.
This resource is useful in the case of [Netdata Streaming and Replication](https://learn.netdata.cloud/docs/observability-centralization-points/metrics-centralization-points/) when you want to spread
the Netdata child agents across different rooms because by default all of them end in the same room like the Netdata parent.`,
		Attributes: map[string]schema.Attribute{
			"room_id": schema.StringAttribute{
				Description: "The Room ID of the space",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"space_id": schema.StringAttribute{
				Description: "Space ID of the member",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_names": schema.ListAttribute{
				Description: "List of node names to add to the room",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (s *nodeRoomMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *nodeRoomMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan nodeRoomMemberResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating node room member for space_id/room_id/node_names: %s/%s/%s", plan.SpaceID.ValueString(), plan.RoomID.ValueString(), plan.NodeNames.String()))

	allNodes, err := s.client.GetAllNodes(plan.SpaceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting All Nodes",
			"err: "+err.Error(),
		)
		return
	}

	planNodes := make([]types.String, 0, len(plan.NodeNames.Elements()))
	diags = plan.NodeNames.ElementsAs(ctx, &planNodes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// precheck if all nodes exist
	for _, planNode := range planNodes {
		exist, _ := checkNodeExists(planNode.ValueString(), allNodes, true)
		if !exist {
			resp.Diagnostics.AddError(
				"Error Creating Node Room Member",
				fmt.Sprintf("Reachable node %s not found in the space %s", planNode.String(), plan.SpaceID.ValueString()),
			)
			return
		}
	}

	for _, planNode := range planNodes {
		_, nodeID := checkNodeExists(planNode.ValueString(), allNodes, true)
		err := s.client.CreateNodeRoomMember(plan.SpaceID.ValueString(), plan.RoomID.ValueString(), nodeID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Node Room Member",
				"err: "+err.Error(),
			)
			return
		}
	}

	plan.RoomID = types.StringValue(plan.RoomID.ValueString())
	plan.SpaceID = types.StringValue(plan.SpaceID.ValueString())
	plan.NodeNames, _ = types.ListValueFrom(ctx, types.StringType, plan.NodeNames)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *nodeRoomMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state nodeRoomMemberResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeRoomMember, err := s.client.GetRoomNodes(state.SpaceID.ValueString(), state.RoomID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Node Room Member",
			fmt.Sprintf("Could not read node room member for space_id/room_id/node_names: %s/%s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), state.NodeNames.String(), err.Error()),
		)
		return
	}

	var refreshedItems []string
	var foundItems bool

	stateNodes := make([]types.String, 0, len(state.NodeNames.Elements()))
	diags = state.NodeNames.ElementsAs(ctx, &stateNodes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, stateNode := range stateNodes {
		exist, _ := checkNodeExists(stateNode.ValueString(), nodeRoomMember, false)
		if exist {
			refreshedItems = append(refreshedItems, stateNode.ValueString())
			foundItems = true
		}
	}

	if !foundItems {
		resp.State.RemoveResource(ctx)
		return
	}

	state.NodeNames, _ = types.ListValueFrom(ctx, types.StringType, refreshedItems)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *nodeRoomMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan nodeRoomMemberResourceModel
	var state nodeRoomMemberResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	allNodes, err := s.client.GetAllNodes(plan.SpaceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting All Nodes",
			"err: "+err.Error(),
		)
		return
	}

	planNodes := make([]types.String, 0, len(plan.NodeNames.Elements()))
	diags = plan.NodeNames.ElementsAs(ctx, &planNodes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// precheck if all nodes exist
	for _, planNode := range planNodes {
		exist, _ := checkNodeExists(planNode.ValueString(), allNodes, true)
		if !exist {
			resp.Diagnostics.AddError(
				"Error Creating Node Room Member",
				fmt.Sprintf("Reachable node %s not found in the space %s", planNode.String(), plan.SpaceID.ValueString()),
			)
			return
		}
	}

	stateNodes := make([]types.String, 0, len(plan.NodeNames.Elements()))
	diags = state.NodeNames.ElementsAs(ctx, &stateNodes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, stateNode := range stateNodes {
		foundState := false
		for _, planNode := range planNodes {
			if stateNode.ValueString() == planNode.ValueString() {
				foundState = true
			}
		}
		if !foundState {
			exist, nodeID := checkNodeExists(stateNode.ValueString(), allNodes, false)
			if exist {
				err := s.client.DeleteNodeRoomMember(state.SpaceID.ValueString(), state.RoomID.ValueString(), nodeID)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Deleting Node Room Member",
						"err: "+err.Error(),
					)
					return
				}
			}
		}
	}

	for _, planNode := range planNodes {
		_, nodeID := checkNodeExists(planNode.ValueString(), allNodes, true)
		err := s.client.CreateNodeRoomMember(plan.SpaceID.ValueString(), plan.RoomID.ValueString(), nodeID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Node Room Member",
				"err: "+err.Error(),
			)
			return
		}
	}

	plan.RoomID = types.StringValue(plan.RoomID.ValueString())
	plan.SpaceID = types.StringValue(plan.SpaceID.ValueString())
	plan.NodeNames, _ = types.ListValueFrom(ctx, types.StringType, plan.NodeNames)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *nodeRoomMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state nodeRoomMemberResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateNodes := make([]types.String, 0, len(state.NodeNames.Elements()))
	diags = state.NodeNames.ElementsAs(ctx, &stateNodes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeRoomMember, err := s.client.GetRoomNodes(state.SpaceID.ValueString(), state.RoomID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Node Room Member",
			fmt.Sprintf("Could not read node room member for space_id/room_id/node_names: %s/%s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), state.NodeNames.String(), err.Error()),
		)
		return
	}

	for _, stateNode := range stateNodes {
		exist, nodeID := checkNodeExists(stateNode.ValueString(), nodeRoomMember, false)
		if exist {
			err := s.client.DeleteNodeRoomMember(state.SpaceID.ValueString(), state.RoomID.ValueString(), nodeID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Deleting Node Room Member",
					"err: "+err.Error(),
				)
				return
			}
		}
	}
}

func (s *nodeRoomMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: space_id,room_id Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("room_id"), idParts[1])...)
}

func checkNodeExists(searchingForNodeName string, nodes *client.RoomNodes, reachableOnly bool) (bool, string) {
	for _, node := range nodes.Nodes {
		if searchingForNodeName == node.NodeName {
			if node.State == "reachable" && reachableOnly {
				return true, node.NodeID
			} else if !reachableOnly {
				return true, node.NodeID
			}
		}
	}
	return false, ""
}

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &roomMemberResource{}
	_ resource.ResourceWithConfigure = &roomMemberResource{}
)

func NewRoomMemberResource() resource.Resource {
	return &roomMemberResource{}
}

type roomMemberResource struct {
	client *client.Client
}

type roomMemberResourceModel struct {
	RoomID        types.String `tfsdk:"room_id"`
	SpaceID       types.String `tfsdk:"space_id"`
	SpaceMemberID types.String `tfsdk:"space_member_id"`
}

func (s *roomMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_room_member"
}

func (s *roomMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Netdata Cloud Room Member resource.",
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
			"space_member_id": schema.StringAttribute{
				Description: "The Space Member ID of the space",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (s *roomMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *roomMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roomMemberResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating room member for space_id/room_id/space_member_id: %s/%s/%s", plan.SpaceID.ValueString(), plan.RoomID.ValueString(), plan.SpaceMemberID.ValueString()))

	err := s.client.CreateRoomMember(plan.SpaceID.ValueString(), plan.RoomID.ValueString(), plan.SpaceMemberID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Room Member",
			"err: "+err.Error(),
		)
		return
	}

	plan.RoomID = types.StringValue(plan.RoomID.ValueString())
	plan.SpaceID = types.StringValue(plan.SpaceID.ValueString())
	plan.SpaceMemberID = types.StringValue(plan.SpaceMemberID.ValueString())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *roomMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roomMemberResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roomMemberInfo, err := s.client.GetRoomMemberID(state.SpaceID.ValueString(), state.RoomID.ValueString(), state.SpaceMemberID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Room Member",
			fmt.Sprintf("Could not read room member for space_id/room_id/space_member_id: %s/%s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), state.SpaceMemberID.ValueString(), err.Error()),
		)
		return
	}

	state.SpaceMemberID = types.StringValue(roomMemberInfo.SpaceMemberID)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *roomMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roomMemberResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roomMemberInfo, err := s.client.GetRoomMemberID(plan.SpaceID.ValueString(), plan.RoomID.ValueString(), plan.SpaceMemberID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Room Member",
			fmt.Sprintf("Could not read room member for space_id/room_id/space_member_id: %s/%s/%s err: %v", plan.SpaceID.ValueString(), plan.RoomID.ValueString(), plan.SpaceMemberID.ValueString(), err.Error()),
		)
		return
	}

	plan.SpaceMemberID = types.StringValue(roomMemberInfo.SpaceMemberID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *roomMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roomMemberResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteRoomMember(state.SpaceID.ValueString(), state.RoomID.ValueString(), state.SpaceMemberID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Room Member",
			fmt.Sprintf("Could not delete room member for space_id/room_id: %s/%s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), state.SpaceMemberID.ValueString(), err.Error()),
		)
		return
	}
}

func (s *roomMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: space_id,room_id,space_member_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("room_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_member_id"), idParts[2])...)
}

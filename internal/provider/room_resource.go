package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &roomResource{}
	_ resource.ResourceWithConfigure = &roomResource{}
)

func NewRoomResource() resource.Resource {
	return &roomResource{}
}

type roomResource struct {
	client *client.Client
}

type roomResourceModel struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (s *roomResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_room"
}

func (s *roomResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the room",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.StringAttribute{
				Description: "The ID of the space",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the room",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the room",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (s *roomResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *roomResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roomResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	roomInfo, err := s.client.CreateRoom(plan.SpaceID.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Room",
			"err: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(roomInfo.ID)
	plan.Name = types.StringValue(roomInfo.Name)
	plan.Description = types.StringValue(roomInfo.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *roomResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roomResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roomInfo, err := s.client.GetRoomByID(state.ID.ValueString(), state.SpaceID.ValueString())

	switch {
	case err == client.ErrNotFound:
		resp.State.RemoveResource(ctx)
		return
	case err != nil:
		resp.Diagnostics.AddError(
			"Error Getting Room",
			"Could Not Read Room ID: "+state.ID.ValueString()+": err: "+err.Error(),
		)
		return
	default:
		state.ID = types.StringValue(roomInfo.ID)
		state.Name = types.StringValue(roomInfo.Name)
		state.Description = types.StringValue(roomInfo.Description)
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

func (s *roomResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roomResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.UpdateRoomByID(plan.ID.ValueString(), plan.SpaceID.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating room",
			"Could Not Update Room ID: "+plan.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}

	roomInfo, err := s.client.GetRoomByID(plan.ID.ValueString(), plan.SpaceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Room",
			"Could Not Read Room ID: "+plan.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(roomInfo.ID)
	plan.Name = types.StringValue(roomInfo.Name)
	plan.Description = types.StringValue(roomInfo.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *roomResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roomResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteRoomByID(state.ID.ValueString(), state.SpaceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Room",
			"Could Not Delete Room ID: "+state.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}
}

func (s *roomResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: space_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

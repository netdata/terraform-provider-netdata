package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &notificationRoomSettingsResource{}
	_ resource.ResourceWithConfigure = &notificationRoomSettingsResource{}
)

func NewNotificationRoomSettingsResource() resource.Resource {
	return &notificationRoomSettingsResource{}
}

type notificationRoomSettingsResource struct {
	client *client.Client
}

type notificationRoomSettingsResourceModel struct {
	ID                types.String `tfsdk:"id"`
	SpaceID           types.String `tfsdk:"space_id"`
	RoomID            types.String `tfsdk:"room_id"`
	ReachabilityDelay types.Int64  `tfsdk:"reachability_delay"`
}

func (s *notificationRoomSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_room_settings"
}

func (s *notificationRoomSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a Netdata Cloud room notification settings resource. The values set here override the space-wide values from ` + "`netdata_notification_space_settings`" + ` for a single room.
Destroying this resource removes the override, and the room falls back to the space-wide value.
Available only in paid plans.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the notification settings.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.StringAttribute{
				Description: "The ID of the space.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					UUID(),
				},
			},
			"room_id": schema.StringAttribute{
				Description: "The ID of the room the override applies to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					UUID(),
				},
			},
			"reachability_delay": schema.Int64Attribute{
				Description: fmt.Sprintf("Node reachability notification delay in seconds for this room. Overrides the space-wide value set by `netdata_notification_space_settings`. Must be at least %d seconds.", minReachabilityDelaySeconds),
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(minReachabilityDelaySeconds),
				},
			},
		},
	}
}

func (s *notificationRoomSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *notificationRoomSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationRoomSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := s.upsert(&plan); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Notification Room Settings",
			fmt.Sprintf("Could not set notification settings for space_id/room_id: %s/%s err: %v", plan.SpaceID.ValueString(), plan.RoomID.ValueString(), err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *notificationRoomSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationRoomSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := s.client.GetRoomNotificationSettings(state.SpaceID.ValueString(), state.RoomID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Notification Room Settings",
			fmt.Sprintf("Could not read notification settings for space_id/room_id: %s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), err.Error()),
		)
		return
	}

	// A null delay means the room no longer overrides the space-wide value, which
	// is the same thing as this resource no longer existing.
	if settings.ReachabilityDelay == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(settings.ID.String())
	state.ReachabilityDelay = types.Int64Value(int64(*settings.ReachabilityDelay))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *notificationRoomSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationRoomSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := s.upsert(&plan); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Notification Room Settings",
			fmt.Sprintf("Could not update notification settings for space_id/room_id: %s/%s err: %v", plan.SpaceID.ValueString(), plan.RoomID.ValueString(), err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *notificationRoomSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationRoomSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteNotificationSettings(state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Notification Room Settings",
			fmt.Sprintf("Could not delete notification settings for space_id/room_id: %s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), err.Error()),
		)
		return
	}
}

func (s *notificationRoomSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

// upsert writes the planned delay and takes the stored row from the response, so
// state carries the ID the API assigned and the value it actually kept.
func (s *notificationRoomSettingsResource) upsert(plan *notificationRoomSettingsResourceModel) error {
	roomID, err := uuid.Parse(plan.RoomID.ValueString())
	if err != nil {
		return err
	}

	delay := int(plan.ReachabilityDelay.ValueInt64())

	response, err := s.client.UpdateNotificationSettings(plan.SpaceID.ValueString(), []client.NotificationSettingsUpsert{
		{
			RoomID:            &roomID,
			ReachabilityDelay: &delay,
		},
	})
	if err != nil {
		return err
	}

	settings, err := response.Room(roomID)
	if err != nil {
		return err
	}

	if settings.ReachabilityDelay == nil {
		return fmt.Errorf("the API returned no reachability delay for room %s after writing %d", plan.RoomID.ValueString(), delay)
	}

	plan.ID = types.StringValue(settings.ID.String())
	plan.ReachabilityDelay = types.Int64Value(int64(*settings.ReachabilityDelay))

	return nil
}

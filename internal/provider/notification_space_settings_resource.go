package provider

import (
	"context"
	"errors"
	"fmt"

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

// minReachabilityDelaySeconds is the lowest delay the API accepts, shared by the
// space-wide and the room-scoped notification settings resources.
const minReachabilityDelaySeconds = 30

var (
	_ resource.Resource              = &notificationSpaceSettingsResource{}
	_ resource.ResourceWithConfigure = &notificationSpaceSettingsResource{}
)

func NewNotificationSpaceSettingsResource() resource.Resource {
	return &notificationSpaceSettingsResource{}
}

type notificationSpaceSettingsResource struct {
	client *client.Client
}

type notificationSpaceSettingsResourceModel struct {
	ID                types.String `tfsdk:"id"`
	SpaceID           types.String `tfsdk:"space_id"`
	ReachabilityDelay types.Int64  `tfsdk:"reachability_delay"`
}

func (s *notificationSpaceSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_space_settings"
}

func (s *notificationSpaceSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a Netdata Cloud space-wide notification settings resource. The values set here apply to every room in the space that does not define its own override through ` + "`netdata_notification_room_settings`" + `.
A room override always wins for the room it targets, so changing the space-wide value does not affect that room. Destroying this resource removes the space-wide setting and the space falls back to the Netdata Cloud default.
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
			"reachability_delay": schema.Int64Attribute{
				Description: fmt.Sprintf("Node reachability notification delay in seconds, applied to every room in the space that does not have its own override. Rooms managed by a `netdata_notification_room_settings` resource ignore this value. Must be at least %d seconds.", minReachabilityDelaySeconds),
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(minReachabilityDelaySeconds),
				},
			},
		},
	}
}

func (s *notificationSpaceSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *notificationSpaceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationSpaceSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := s.upsert(&plan); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Notification Space Settings",
			fmt.Sprintf("Could not set notification settings for space_id: %s err: %v", plan.SpaceID.ValueString(), err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *notificationSpaceSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationSpaceSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := s.client.GetSpaceDefaultNotificationSettings(state.SpaceID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Notification Space Settings",
			fmt.Sprintf("Could not read notification settings for space_id: %s err: %v", state.SpaceID.ValueString(), err.Error()),
		)
		return
	}

	// The space-wide row exists for reasons beyond this resource, so a null delay
	// means the value this resource manages is no longer set.
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

func (s *notificationSpaceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationSpaceSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := s.upsert(&plan); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Notification Space Settings",
			fmt.Sprintf("Could not update notification settings for space_id: %s err: %v", plan.SpaceID.ValueString(), err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *notificationSpaceSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationSpaceSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteNotificationSettings(state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Notification Space Settings",
			fmt.Sprintf("Could not delete notification settings for space_id/id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (s *notificationSpaceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: space_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), req.ID)...)
}

// upsert writes the planned delay and takes the stored row from the response, so
// state carries the ID the API assigned and the value it actually kept.
func (s *notificationSpaceSettingsResource) upsert(plan *notificationSpaceSettingsResourceModel) error {
	delay := int(plan.ReachabilityDelay.ValueInt64())

	response, err := s.client.UpdateNotificationSettings(plan.SpaceID.ValueString(), []client.NotificationSettingsUpsert{
		{
			ReachabilityDelay: &delay,
		},
	})
	if err != nil {
		return err
	}

	settings, err := response.SpaceDefault()
	if err != nil {
		return err
	}

	if settings.ReachabilityDelay == nil {
		return fmt.Errorf("the API returned no reachability delay for the space-wide settings after writing %d", delay)
	}

	plan.ID = types.StringValue(settings.ID.String())
	plan.ReachabilityDelay = types.Int64Value(int64(*settings.ReachabilityDelay))

	return nil
}

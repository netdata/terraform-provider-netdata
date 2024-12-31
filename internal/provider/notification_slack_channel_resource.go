package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &slackChannelResource{}
	_ resource.ResourceWithConfigure = &slackChannelResource{}
)

func NewSlackChannelResource() resource.Resource {
	return &slackChannelResource{}
}

type slackChannelResource struct {
	client *client.Client
}

type slackChannelResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	SpaceID                  types.String `tfsdk:"space_id"`
	RoomsID                  types.List   `tfsdk:"rooms_id"`
	Alarms                   types.String `tfsdk:"alarms"`
	RepeatNotificationMinute types.Int64  `tfsdk:"repeat_notification_min"`
	WebhookURL               types.String `tfsdk:"webhook_url"`
}

func (s *slackChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_slack_channel"
}

func (s *slackChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fullSchema := commonNotificationSchema("Slack")
	fullSchema.Attributes["webhook_url"] = schema.StringAttribute{
		Description: "Slack webhook URL",
		Required:    true,
		Sensitive:   true,
	}
	resp.Schema = fullSchema
}

func (s *slackChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *slackChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan slackChannelResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notificationIntegration, err := s.client.GetNotificationIntegrationByType(plan.SpaceID.ValueString(), "slack")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Slack Notification",
			"err: "+err.Error(),
		)
		return
	}

	var roomsID []string
	plan.RoomsID.ElementsAs(ctx, &roomsID, false)

	commonParams := client.NotificationChannel{
		Name:                     plan.Name.ValueString(),
		Integration:              *notificationIntegration,
		Rooms:                    roomsID,
		Alarms:                   plan.Alarms.ValueString(),
		Enabled:                  plan.Enabled.ValueBool(),
		RepeatNotificationMinute: plan.RepeatNotificationMinute.ValueInt64(),
	}

	slackParams := client.NotificationSlackChannel{
		URL: plan.WebhookURL.ValueString(),
	}

	notificationChannel, err := s.client.CreateSlackChannel(plan.SpaceID.ValueString(), commonParams, slackParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Slack Notification",
			"err: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(notificationChannel.ID)
	plan.Name = types.StringValue(notificationChannel.Name)
	plan.Enabled = types.BoolValue(notificationChannel.Enabled)
	plan.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	plan.Alarms = types.StringValue(notificationChannel.Alarms)
	plan.RepeatNotificationMinute = types.Int64Value(notificationChannel.RepeatNotificationMinute)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *slackChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state slackChannelResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notificationChannel, err := s.client.GetNotificationChannelByIDAndType(state.SpaceID.ValueString(), state.ID.ValueString(), "slack")
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Slack Notification",
			fmt.Sprintf("Could not read slack notification for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}

	var notificationSecrets client.NotificationSlackChannel
	err = json.Unmarshal(notificationChannel.Secrets, &notificationSecrets)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Slack Notification",
			fmt.Sprintf("Could not unmarshal slack notification secrets for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}
	state.Name = types.StringValue(notificationChannel.Name)
	state.Enabled = types.BoolValue(notificationChannel.Enabled)
	state.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	state.Alarms = types.StringValue(notificationChannel.Alarms)
	state.RepeatNotificationMinute = types.Int64Value(notificationChannel.RepeatNotificationMinute)
	state.WebhookURL = types.StringValue(notificationSecrets.URL)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *slackChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan slackChannelResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var roomsID []string
	plan.RoomsID.ElementsAs(ctx, &roomsID, false)

	commonParams := client.NotificationChannel{
		ID:                       plan.ID.ValueString(),
		Name:                     plan.Name.ValueString(),
		Rooms:                    roomsID,
		Alarms:                   plan.Alarms.ValueString(),
		Enabled:                  plan.Enabled.ValueBool(),
		RepeatNotificationMinute: plan.RepeatNotificationMinute.ValueInt64(),
	}

	slackParams := client.NotificationSlackChannel{
		URL: plan.WebhookURL.ValueString(),
	}

	notificationChannel, err := s.client.UpdateSlackChannelByID(plan.SpaceID.ValueString(), commonParams, slackParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Slack Notification",
			fmt.Sprintf("Could not update slack notification for space_id/channel_id: %s/%s err: %v", plan.SpaceID.ValueString(), plan.ID.ValueString(), err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(notificationChannel.ID)
	plan.Name = types.StringValue(notificationChannel.Name)
	plan.Enabled = types.BoolValue(notificationChannel.Enabled)
	plan.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	plan.Alarms = types.StringValue(notificationChannel.Alarms)
	plan.RepeatNotificationMinute = types.Int64Value(notificationChannel.RepeatNotificationMinute)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *slackChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state slackChannelResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteChannelByID(state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Slack Notification",
			fmt.Sprintf("Could not delete slack notification for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (s *slackChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: space_id,channel_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

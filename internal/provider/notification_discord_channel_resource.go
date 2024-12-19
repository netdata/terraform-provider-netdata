package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &discordChannelResource{}
	_ resource.ResourceWithConfigure = &discordChannelResource{}
)

func NewDiscordChannelResource() resource.Resource {
	return &discordChannelResource{}
}

type discordChannelResource struct {
	client *client.Client
}

type discordChannelResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	SpaceID       types.String `tfsdk:"space_id"`
	RoomsID       types.List   `tfsdk:"rooms_id"`
	Alarms        types.String `tfsdk:"alarms"`
	WebhookURL    types.String `tfsdk:"webhook_url"`
	ChannelType   types.String `tfsdk:"channel_type"`
	ChannelThread types.String `tfsdk:"channel_thread"`
}

func (s *discordChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_discord_channel"
}

func (s *discordChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fullSchema := commonNotificationSchema("Discord")
	fullSchema.Attributes["webhook_url"] = schema.StringAttribute{
		Description: "Discord webhook URL",
		Required:    true,
		Sensitive:   true,
	}
	fullSchema.Attributes["channel_type"] = schema.StringAttribute{
		Description: "Discord channel type. Valid values are: `text`, `forum`",
		Required:    true,
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^(text|forum)$`),
				"Invalid channel type",
			),
		},
	}
	fullSchema.Attributes["channel_thread"] = schema.StringAttribute{
		Description: "Discord channel thread name required if channel type is `forum`",
		Optional:    true,
	}
	resp.Schema = fullSchema
}

func (s *discordChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *discordChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan discordChannelResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ChannelType.ValueString() == "forum" && plan.ChannelThread.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Error Creating Discord Notification",
			"channel_thread is required if channel_type is forum",
		)
		return
	}

	notificationIntegration, err := s.client.GetNotificationIntegrationByType(plan.SpaceID.ValueString(), "discord")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Discord Notification",
			"err: "+err.Error(),
		)
		return
	}

	var roomsID []string
	plan.RoomsID.ElementsAs(ctx, &roomsID, false)

	commonParams := client.NotificationChannel{
		Name:        plan.Name.ValueString(),
		Integration: *notificationIntegration,
		Rooms:       roomsID,
		Alarms:      plan.Alarms.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
	}

	discordParams := client.NotificationDiscordChannel{
		URL: plan.WebhookURL.ValueString(),
	}

	switch plan.ChannelType.ValueString() {
	case "text":
		discordParams.ChannelParams.Selection = "text"
	case "forum":
		discordParams.ChannelParams.Selection = "forum"
		discordParams.ChannelParams.ThreadName = plan.ChannelThread.ValueString()
	}

	notificationChannel, err := s.client.CreateDiscordChannel(plan.SpaceID.ValueString(), commonParams, discordParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Discord Notification",
			"err: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(notificationChannel.ID)
	plan.Name = types.StringValue(notificationChannel.Name)
	plan.Enabled = types.BoolValue(notificationChannel.Enabled)
	plan.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	plan.Alarms = types.StringValue(notificationChannel.Alarms)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *discordChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state discordChannelResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notificationChannel, err := s.client.GetNotificationChannelByIDAndType(state.SpaceID.ValueString(), state.ID.ValueString(), "discord")
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Discord Notification",
			fmt.Sprintf("Could not read discord notification for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}

	var notificationSecrets client.NotificationDiscordChannel
	err = json.Unmarshal(notificationChannel.Secrets, &notificationSecrets)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Discord Notification",
			fmt.Sprintf("Could not unmarshal discord notification secrets for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}
	state.Name = types.StringValue(notificationChannel.Name)
	state.Enabled = types.BoolValue(notificationChannel.Enabled)
	state.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	state.Alarms = types.StringValue(notificationChannel.Alarms)
	state.WebhookURL = types.StringValue(notificationSecrets.URL)
	state.ChannelType = types.StringValue(notificationSecrets.ChannelParams.Selection)
	if notificationSecrets.ChannelParams.Selection == "forum" {
		state.ChannelThread = types.StringValue(notificationSecrets.ChannelParams.ThreadName)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *discordChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan discordChannelResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ChannelType.ValueString() == "forum" && plan.ChannelThread.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Error Creating Discord Notification",
			"channel_thread is required if channel_type is forum",
		)
		return
	}

	var roomsID []string
	plan.RoomsID.ElementsAs(ctx, &roomsID, false)

	commonParams := client.NotificationChannel{
		ID:      plan.ID.ValueString(),
		Name:    plan.Name.ValueString(),
		Rooms:   roomsID,
		Alarms:  plan.Alarms.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
	}

	discordParams := client.NotificationDiscordChannel{
		URL: plan.WebhookURL.ValueString(),
	}

	switch plan.ChannelType.ValueString() {
	case "text":
		discordParams.ChannelParams.Selection = "text"
	case "forum":
		discordParams.ChannelParams.Selection = "forum"
		discordParams.ChannelParams.ThreadName = plan.ChannelThread.ValueString()
	}

	notificationChannel, err := s.client.UpdateDiscordChannelByID(plan.SpaceID.ValueString(), commonParams, discordParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Discord Notification",
			fmt.Sprintf("Could not update discord notification for space_id/channel_id: %s/%s err: %v", plan.SpaceID.ValueString(), plan.ID.ValueString(), err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(notificationChannel.ID)
	plan.Name = types.StringValue(notificationChannel.Name)
	plan.Enabled = types.BoolValue(notificationChannel.Enabled)
	plan.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	plan.Alarms = types.StringValue(notificationChannel.Alarms)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *discordChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state discordChannelResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteChannelByID(state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Discord Notification",
			fmt.Sprintf("Could not delete discord notification for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (s *discordChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

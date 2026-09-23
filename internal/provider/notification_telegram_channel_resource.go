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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &telegramChannelResource{}
	_ resource.ResourceWithConfigure = &telegramChannelResource{}
)

func NewTelegramChannelResource() resource.Resource {
	return &telegramChannelResource{}
}

type telegramChannelResource struct {
	client *client.Client
}

type telegramChannelResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	SpaceID                  types.String `tfsdk:"space_id"`
	RoomsID                  types.List   `tfsdk:"rooms_id"`
	NotificationOptions      types.List   `tfsdk:"notifications"`
	RepeatNotificationMinute types.Int64  `tfsdk:"repeat_notification_min"`
	BotToken                 types.String `tfsdk:"bot_token"`
	ChatID                   types.String `tfsdk:"chat_id"`
	TopicID                  types.Int64  `tfsdk:"topic_id"`
}

func (s *telegramChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_telegram_channel"
}

func (s *telegramChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fullSchema := commonNotificationSchema("Telegram")
	fullSchema.Attributes["bot_token"] = schema.StringAttribute{
		Description: "Telegram bot token",
		Required:    true,
		Sensitive:   true,
	}
	fullSchema.Attributes["chat_id"] = schema.StringAttribute{
		Description: "Unique identifier for the target chat",
		Required:    true,
	}
	fullSchema.Attributes["topic_id"] = schema.Int64Attribute{
		Description: "Unique identifier for the target topic of a chat. If omitted or 0, messages are sent to the General topic. If topics are not supported messages are simply sent to the chat.",
		Computed:    true,
		Optional:    true,
		Default:     int64default.StaticInt64(0),
	}
	resp.Schema = fullSchema
}

func (s *telegramChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *telegramChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan telegramChannelResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notificationIntegration, err := s.client.GetNotificationIntegrationByType(plan.SpaceID.ValueString(), "telegram")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Telegram Notification",
			"err: "+err.Error(),
		)
		return
	}

	var roomsID []string
	plan.RoomsID.ElementsAs(ctx, &roomsID, false)

	var notificationOptions []string
	plan.NotificationOptions.ElementsAs(ctx, &notificationOptions, false)

	commonParams := client.NotificationChannel{
		Name:                     plan.Name.ValueString(),
		Integration:              *notificationIntegration,
		Rooms:                    roomsID,
		NotificationOptions:      notificationOptions,
		Enabled:                  plan.Enabled.ValueBool(),
		RepeatNotificationMinute: plan.RepeatNotificationMinute.ValueInt64(),
	}

	telegramParams := client.NotificationTelegramChannel{
		BotToken: plan.BotToken.ValueString(),
		ChatID:   plan.ChatID.ValueString(),
		TopicID:  plan.TopicID.ValueInt64(),
	}

	notificationChannel, err := s.client.CreateTelegramChannel(plan.SpaceID.ValueString(), commonParams, telegramParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Telegram Notification",
			"err: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(notificationChannel.ID)
	plan.Name = types.StringValue(notificationChannel.Name)
	plan.Enabled = types.BoolValue(notificationChannel.Enabled)
	plan.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	plan.NotificationOptions, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.NotificationOptions)
	plan.RepeatNotificationMinute = types.Int64Value(notificationChannel.RepeatNotificationMinute)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *telegramChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state telegramChannelResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notificationChannel, err := s.client.GetNotificationChannelByIDAndType(state.SpaceID.ValueString(), state.ID.ValueString(), "telegram")
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Telegram Notification",
			fmt.Sprintf("Could not read telegram notification for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}

	var notificationSecrets client.NotificationTelegramChannel
	err = json.Unmarshal(notificationChannel.Secrets, &notificationSecrets)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Telegram Notification",
			fmt.Sprintf("Could not unmarshal telegram notification secrets for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}
	state.Name = types.StringValue(notificationChannel.Name)
	state.Enabled = types.BoolValue(notificationChannel.Enabled)
	state.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	state.NotificationOptions, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.NotificationOptions)
	state.RepeatNotificationMinute = types.Int64Value(notificationChannel.RepeatNotificationMinute)
	state.BotToken = types.StringValue(notificationSecrets.BotToken)
	state.ChatID = types.StringValue(notificationSecrets.ChatID)
	state.TopicID = types.Int64Value(notificationSecrets.TopicID)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *telegramChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan telegramChannelResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var roomsID []string
	plan.RoomsID.ElementsAs(ctx, &roomsID, false)

	var notificationOptions []string
	plan.NotificationOptions.ElementsAs(ctx, &notificationOptions, false)

	commonParams := client.NotificationChannel{
		ID:                       plan.ID.ValueString(),
		Name:                     plan.Name.ValueString(),
		Rooms:                    roomsID,
		NotificationOptions:      notificationOptions,
		Enabled:                  plan.Enabled.ValueBool(),
		RepeatNotificationMinute: plan.RepeatNotificationMinute.ValueInt64(),
	}

	telegramParams := client.NotificationTelegramChannel{
		BotToken: plan.BotToken.ValueString(),
		ChatID:   plan.ChatID.ValueString(),
		TopicID:  plan.TopicID.ValueInt64(),
	}

	notificationChannel, err := s.client.UpdateTelegramChannelByID(plan.SpaceID.ValueString(), commonParams, telegramParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Telegram Notification",
			fmt.Sprintf("Could not update telegram notification for space_id/channel_id: %s/%s err: %v", plan.SpaceID.ValueString(), plan.ID.ValueString(), err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(notificationChannel.ID)
	plan.Name = types.StringValue(notificationChannel.Name)
	plan.Enabled = types.BoolValue(notificationChannel.Enabled)
	plan.RoomsID, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.Rooms)
	plan.NotificationOptions, _ = types.ListValueFrom(ctx, types.StringType, notificationChannel.NotificationOptions)
	plan.RepeatNotificationMinute = types.Int64Value(notificationChannel.RepeatNotificationMinute)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *telegramChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state telegramChannelResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteChannelByID(state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Telegram Notification",
			fmt.Sprintf("Could not delete telegram notification for space_id/channel_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (s *telegramChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &notificationSilencingRuleResource{}
	_ resource.ResourceWithConfigure = &notificationSilencingRuleResource{}
)

func NewNotificationSilencingRule() resource.Resource {
	return &notificationSilencingRuleResource{}
}

type notificationSilencingRuleResource struct {
	client *client.Client
}

type notificationSilencingRuleResourceModel struct {
	ID                  types.String      `tfsdk:"id"`
	Name                types.String      `tfsdk:"name"`
	SpaceID             types.String      `tfsdk:"space_id"`
	RoomIDs             types.List        `tfsdk:"room_ids"`
	NodeIDs             types.List        `tfsdk:"node_ids"`
	HostLabels          types.Map         `tfsdk:"host_labels"`
	AlertNames          types.List        `tfsdk:"alert_names"`
	AlertContexts       types.List        `tfsdk:"alert_contexts"`
	AlertInstances      types.List        `tfsdk:"alert_instances"`
	AlertRoles          types.List        `tfsdk:"alert_roles"`
	NotificationOptions types.List        `tfsdk:"notification_options"`
	IntegrationIDs      types.List        `tfsdk:"integration_ids"`
	StartsAt            timetypes.RFC3339 `tfsdk:"starts_at"`
	LastsUntil          timetypes.RFC3339 `tfsdk:"lasts_until"`
	DeleteOnExpiry      types.Bool        `tfsdk:"delete_on_expiry"`
	Disabled            types.Bool        `tfsdk:"disabled"`
	RRule               types.String      `tfsdk:"rrule"`
	Timezone            types.String      `tfsdk:"timezone"`
}

func (s *notificationSilencingRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_silencing_rule"
}

func (s *notificationSilencingRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a Netdata Cloud Notification Silencing Rule resource. Use this resource to manage notification silencing rules in Netdata Cloud.
A notification silencing rule allows you to silence notifications for specific alerts, nodes, rooms, or spaces based on various criteria.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the silencing rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the silencing rule.",
				Required:    true,
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
			"room_ids": schema.ListAttribute{
				Description: "List of room IDs to apply the silencing rule to. When empty it applies to all rooms in the space.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						UUID(),
					),
				},
			},
			"node_ids": schema.ListAttribute{
				Description: "List of node IDs to apply the silencing rule to. When empty it applies to all nodes in the space.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						UUID(),
					),
				},
			},
			"host_labels": schema.MapAttribute{
				Description: "Host labels to filter nodes the silencing rule applies to. When empty it applies to any label. Example { os = \"linux\" }.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			},
			"alert_names": schema.ListAttribute{
				Description: "List of alert names to silence. When empty it applies to all alert names. Example [\"disk.space\"].",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"alert_contexts": schema.ListAttribute{
				Description: "List of alert contexts to silence. When empty it applies to all alert contexts. Example [\"disk_space_usage\"].",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"alert_instances": schema.ListAttribute{
				Description: "List of alert instances to silence. When empty it applies to all alert instances. Example [\"disk_space./\"].",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"alert_roles": schema.ListAttribute{
				Description: "List of alert roles to silence. When empty it applies to all alert roles. Example [\"webmaster\"].",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"notification_options": schema.ListAttribute{
				Description: "List of notification options to silence. Valid values: CRITICAL, WARNING, CLEAR, REACHABLE, UNREACHABLE.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("CRITICAL", "WARNING", "CLEAR", "REACHABLE", "UNREACHABLE"),
					),
				},
			},
			"integration_ids": schema.ListAttribute{
				Description: "List of integration (notification methods) IDs to apply the silencing rule to. When empty it applies to all integrations in the space.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						UUID(),
					),
				},
			},
			"starts_at": schema.StringAttribute{
				Description: "The start time of the silencing rule in RFC3339 format. Must be expressed in UTC, e.g. \"2026-09-01T00:00:00Z\". A value carrying a non-UTC offset such as \"+02:00\" is rejected at plan time, because the API stores and returns timestamps in UTC. Use `timezone` to control the schedule's local time.",
				CustomType:  timetypes.RFC3339Type{},
				Required:    true,
				Validators: []validator.String{
					UTCTimestamp(),
				},
			},
			"lasts_until": schema.StringAttribute{
				Description: "The end time of the silencing rule in RFC3339 format. Must be expressed in UTC, e.g. \"2026-12-31T23:59:59Z\". A value carrying a non-UTC offset such as \"+02:00\" is rejected at plan time, because the API stores and returns timestamps in UTC. If not set, the rule does not expire.",
				CustomType:  timetypes.RFC3339Type{},
				Optional:    true,
				Validators: []validator.String{
					UTCTimestamp(),
				},
			},
			"delete_on_expiry": schema.BoolAttribute{
				Description: "Whether to delete the silencing rule when it expires. If enabled, then it requires setting lasts_until. It is not compatible with an RRule.. then Defaults to false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disabled": schema.BoolAttribute{
				Description: "Whether the silencing rule is disabled. Defaults to false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"rrule": schema.StringAttribute{
				Description: "The recurrence rule (RRULE) for the silencing rule. The format orignates from the iCalendar specification (RFC 5545). Must not include a DTSTART line - the recurrence start is taken from starts_at. Example: \"RRULE:FREQ=MONTHLY;INTERVAL=1;COUNT=10;BYMONTHDAY=1\".",
				Optional:    true,
				Validators: []validator.String{
					RRuleOnly(),
				},
			},
			"timezone": schema.StringAttribute{
				Description: "The timezone for the silencing rule schedule. Must be a valid IANA timezone (e.g. \"America/New_York\", \"UTC\"). Defaults to \"UTC\".",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("UTC"),
				Validators: []validator.String{
					IANATimezone(),
				},
			},
		},
	}
}

func (s *notificationSilencingRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *notificationSilencingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationSilencingRuleResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating notification silencing rule with name: %s in space: %s", plan.Name.ValueString(), plan.SpaceID.ValueString()))

	silencingRule, _ := planToSilencingRule(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	silencingRuleResponse, err := s.client.CreateSilencingRule(plan.SpaceID.ValueString(), silencingRule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Notification Silencing Rule",
			"err: "+err.Error(),
		)
		return
	}

	silencingRuleResponseToplan(ctx, silencingRuleResponse, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *notificationSilencingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationSilencingRuleResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading notification silencing rule with ID: %s in space: %s", state.ID.ValueString(), state.SpaceID.ValueString()))

	silencingRules, err := s.client.GetSilencingRule(state.ID.ValueString(), state.SpaceID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Notification Silencing Rules",
			fmt.Sprintf("Could not read notification silencing rules for space_id: %s err: %v", state.SpaceID.ValueString(), err.Error()),
		)
		return
	}

	silencingRuleResponseToplan(ctx, silencingRules, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	// SpaceID is not part of the API response — preserve it from state
	state.SpaceID = types.StringValue(state.SpaceID.ValueString())

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (s *notificationSilencingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationSilencingRuleResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating notification silencing rule with ID: %s in space: %s", plan.ID.ValueString(), plan.SpaceID.ValueString()))

	silencingRule, _ := planToSilencingRule(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	silencingRuleResponse, err := s.client.UpdateSilencingRule(plan.ID.ValueString(), plan.SpaceID.ValueString(), silencingRule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Notification Silencing Rule",
			"Could Not Update Notification Silencing Rule ID: "+plan.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}

	silencingRuleResponseToplan(ctx, silencingRuleResponse, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *notificationSilencingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationSilencingRuleResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteSilencingRule(state.ID.ValueString(), state.SpaceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Notification Silencing Rule",
			"Could Not Delete Notification Silencing Rule ID: "+state.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}
}

func (s *notificationSilencingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: space_id,silencing_rule_id Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

// planToSilencingRule converts a notificationSilencingRuleResourceModel to a client.SilencingRule.
func planToSilencingRule(ctx context.Context, plan notificationSilencingRuleResourceModel, diags *diag.Diagnostics) (client.SilencingRule, error) {
	rule := client.SilencingRule{
		Name:           plan.Name.ValueString(),
		DeleteOnExpiry: plan.DeleteOnExpiry.ValueBool(),
		Disabled:       plan.Disabled.ValueBool(),
	}

	// StartsAt (required)
	startsAt, d := plan.StartsAt.ValueRFC3339Time()
	diags.Append(d...)
	if d.HasError() {
		return rule, fmt.Errorf("invalid starts_at: %q", plan.StartsAt.ValueString())
	}
	rule.StartsAt = startsAt

	// LastsUntil (optional)
	if !plan.LastsUntil.IsNull() && !plan.LastsUntil.IsUnknown() {
		lastsUntil, d := plan.LastsUntil.ValueRFC3339Time()
		diags.Append(d...)
		if d.HasError() {
			return rule, fmt.Errorf("invalid lasts_until: %q", plan.LastsUntil.ValueString())
		}
		rule.LastsUntil = &lastsUntil
	}

	// RRule (optional)
	if !plan.RRule.IsNull() && !plan.RRule.IsUnknown() {
		rrule := plan.RRule.ValueString()
		rule.RRule = &rrule
	}

	// Timezone (optional)
	if !plan.Timezone.IsNull() && !plan.Timezone.IsUnknown() {
		tz := plan.Timezone.ValueString()
		rule.Timezone = &tz
	}

	// RoomIDs (optional, []uuid.UUID)
	if !plan.RoomIDs.IsNull() && !plan.RoomIDs.IsUnknown() {
		var roomIDStrs []string
		diags.Append(plan.RoomIDs.ElementsAs(ctx, &roomIDStrs, false)...)
		for _, s := range roomIDStrs {
			id, err := uuid.Parse(s)
			if err != nil {
				diags.AddError("Invalid room_id", fmt.Sprintf("%q is not a valid UUID: %s", s, err))
				return rule, err
			}
			rule.RoomIDs = append(rule.RoomIDs, id)
		}
	}

	// NodeIDs (optional, []uuid.UUID)
	if !plan.NodeIDs.IsNull() && !plan.NodeIDs.IsUnknown() {
		var nodeIDStrs []string
		diags.Append(plan.NodeIDs.ElementsAs(ctx, &nodeIDStrs, false)...)
		for _, s := range nodeIDStrs {
			id, err := uuid.Parse(s)
			if err != nil {
				diags.AddError("Invalid node_id", fmt.Sprintf("%q is not a valid UUID: %s", s, err))
				return rule, err
			}
			rule.NodeIDs = append(rule.NodeIDs, id)
		}
	}

	// IntegrationIDs (optional, []uuid.UUID)
	if !plan.IntegrationIDs.IsNull() && !plan.IntegrationIDs.IsUnknown() {
		var integrationIDStrs []string
		diags.Append(plan.IntegrationIDs.ElementsAs(ctx, &integrationIDStrs, false)...)
		for _, s := range integrationIDStrs {
			id, err := uuid.Parse(s)
			if err != nil {
				diags.AddError("Invalid integration_id", fmt.Sprintf("%q is not a valid UUID: %s", s, err))
				return rule, err
			}
			rule.IntegrationIDs = append(rule.IntegrationIDs, id)
		}
	}

	// HostLabels (optional, map[string]string)
	if !plan.HostLabels.IsNull() && !plan.HostLabels.IsUnknown() {
		hostLabels := make(map[string]string, len(plan.HostLabels.Elements()))
		diags.Append(plan.HostLabels.ElementsAs(ctx, &hostLabels, false)...)
		rule.HostLabels = hostLabels
	}

	// String lists (optional)
	if !plan.AlertNames.IsNull() && !plan.AlertNames.IsUnknown() {
		diags.Append(plan.AlertNames.ElementsAs(ctx, &rule.AlertNames, false)...)
	}
	if !plan.AlertContexts.IsNull() && !plan.AlertContexts.IsUnknown() {
		diags.Append(plan.AlertContexts.ElementsAs(ctx, &rule.AlertContexts, false)...)
	}
	if !plan.AlertInstances.IsNull() && !plan.AlertInstances.IsUnknown() {
		diags.Append(plan.AlertInstances.ElementsAs(ctx, &rule.AlertInstances, false)...)
	}
	if !plan.AlertRoles.IsNull() && !plan.AlertRoles.IsUnknown() {
		diags.Append(plan.AlertRoles.ElementsAs(ctx, &rule.AlertRoles, false)...)
	}

	// NotificationOptions (optional, *[]string)
	if !plan.NotificationOptions.IsNull() && !plan.NotificationOptions.IsUnknown() {
		var opts []string
		diags.Append(plan.NotificationOptions.ElementsAs(ctx, &opts, false)...)
		rule.NotificationOptions = opts
	}

	return rule, nil
}

// silencingRuleResponseToplan maps a client.SilencingRule response back into the plan model.
func silencingRuleResponseToplan(ctx context.Context, r *client.SilencingRule, plan *notificationSilencingRuleResourceModel) {
	plan.ID = types.StringValue(r.ID.String())
	plan.Name = types.StringValue(r.Name)
	plan.DeleteOnExpiry = types.BoolValue(r.DeleteOnExpiry)
	plan.Disabled = types.BoolValue(r.Disabled)
	plan.StartsAt = timetypes.NewRFC3339TimeValue(r.StartsAt)
	plan.LastsUntil = timetypes.NewRFC3339TimePointerValue(r.LastsUntil)

	if r.RRule != nil {
		if rrule := stripDTSTART(*r.RRule); rrule != "" {
			plan.RRule = types.StringValue(rrule)
		} else {
			plan.RRule = types.StringNull()
		}
	} else {
		plan.RRule = types.StringNull()
	}

	if r.Timezone != nil {
		plan.Timezone = types.StringValue(*r.Timezone)
	} else {
		plan.Timezone = types.StringNull()
	}

	// UUID lists → []string
	roomIDs := make([]string, len(r.RoomIDs))
	for i, id := range r.RoomIDs {
		roomIDs[i] = id.String()
	}
	plan.RoomIDs, _ = types.ListValueFrom(ctx, types.StringType, roomIDs)

	nodeIDs := make([]string, len(r.NodeIDs))
	for i, id := range r.NodeIDs {
		nodeIDs[i] = id.String()
	}
	plan.NodeIDs, _ = types.ListValueFrom(ctx, types.StringType, nodeIDs)

	integrationIDs := make([]string, len(r.IntegrationIDs))
	for i, id := range r.IntegrationIDs {
		integrationIDs[i] = id.String()
	}
	plan.IntegrationIDs, _ = types.ListValueFrom(ctx, types.StringType, integrationIDs)

	// HostLabels
	if len(r.HostLabels) > 0 {
		plan.HostLabels, _ = types.MapValueFrom(ctx, types.StringType, r.HostLabels)
	} else {
		plan.HostLabels, _ = types.MapValueFrom(ctx, types.StringType, map[string]string{})
	}

	// String lists
	plan.AlertNames = stringListValue(ctx, r.AlertNames)
	plan.AlertContexts = stringListValue(ctx, r.AlertContexts)
	plan.AlertInstances = stringListValue(ctx, r.AlertInstances)
	plan.AlertRoles = stringListValue(ctx, r.AlertRoles)
	plan.NotificationOptions = stringListValue(ctx, r.NotificationOptions)
}

// stringListValue converts a []string into a types.List, mapping a nil slice to
// an empty list instead of null. The API omits empty collections from its
// response, and returning null for an attribute the plan set to [] makes
// Terraform reject the applied value as inconsistent with the plan.
func stringListValue(ctx context.Context, values []string) types.List {
	if values == nil {
		values = []string{}
	}
	list, _ := types.ListValueFrom(ctx, types.StringType, values)
	return list
}

// stripDTSTART removes the leading DTSTART line(s) the API prepends to an RRULE.
// The configuration only ever carries the recurrence rule itself (see RRuleOnly),
// so the schedule start has to be dropped again to keep state consistent with config.
func stripDTSTART(rrule string) string {
	for {
		trimmed := strings.TrimLeft(rrule, " \t\r\n")
		if !strings.HasPrefix(strings.ToUpper(trimmed), "DTSTART") {
			return strings.TrimSpace(trimmed)
		}
		idx := strings.IndexAny(trimmed, "\r\n")
		if idx == -1 {
			return ""
		}
		rrule = trimmed[idx+1:]
	}
}

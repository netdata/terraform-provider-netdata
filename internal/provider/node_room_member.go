package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
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
	RoomID    types.String             `tfsdk:"room_id"`
	SpaceID   types.String             `tfsdk:"space_id"`
	NodeNames types.List               `tfsdk:"node_names"`
	Rules     []nodeRoomMembershipRule `tfsdk:"rule"`
}
type nodeRoomMembershipRule struct {
	ID          types.String               `tfsdk:"id"`
	Action      types.String               `tfsdk:"action"`
	Description types.String               `tfsdk:"description"`
	Clauses     []nodeRoomMembershipClause `tfsdk:"clause"`
}

type nodeRoomMembershipClause struct {
	Label    types.String `tfsdk:"label"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
	Negate   types.Bool   `tfsdk:"negate"`
}

func (s *nodeRoomMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node_room_member"
}

func (s *nodeRoomMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a Netdata Cloud Node Room Member resource. Use this resource to manage node membership to the room in the selected space.
There are two options to add nodes to the room:
- providing the node names directly, but only reachable nodes will be added to the room, use node_names attribute for this
- creating rules that will automatically add nodes to the room based on the rule, use rule block for this
`,
		Attributes: map[string]schema.Attribute{
			"room_id": schema.StringAttribute{
				Description: "The Room ID of the space.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"space_id": schema.StringAttribute{
				Description: "Space ID of the member.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_names": schema.ListAttribute{
				Description: "List of node names to add to the room. At least one node name is required.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListNull(types.StringType)),
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "The node rule to apply to the room. The logical relation between multiple rules is OR. More info [here](https://learn.netdata.cloud/docs/netdata-cloud/spaces-and-rooms/node-rule-based-room-assignment).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the rule.",
							Computed:    true,
						},
						"action": schema.StringAttribute{
							Description: "Determines whether matching nodes will be included or excluded from the room. Valid values: INCLUDE or EXCLUDE. EXCLUDE action always takes precedence against INCLUDE.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"INCLUDE", "EXCLUDE"}...),
							},
						},
						"description": schema.StringAttribute{
							Description: "The description of the rule.",
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"clause": schema.ListNestedBlock{
							Description: "The clause to apply to the rule. The logical relation between multiple clauses is AND. It should be a least one clause.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"label": schema.StringAttribute{
										Description: "The host label to check.",
										Required:    true,
									},
									"operator": schema.StringAttribute{
										Description: "Operator to compare. Valid values: equals, starts_with, ends_with, contains.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.OneOf([]string{"equals", "starts_with", "ends_with", "contains"}...),
										},
									},
									"value": schema.StringAttribute{
										Description: "The value to compare against.",
										Required:    true,
									},
									"negate": schema.BoolAttribute{
										Description: "Negate the clause.",
										Required:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.AtLeastOneOf(
									path.MatchRoot("clause"),
								),
							},
						},
					},
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

	for i, rule := range plan.Rules {
		var nodeMembershipClauses []client.NodeMembershipClause
		for _, clause := range rule.Clauses {
			nodeMembershipClauses = append(nodeMembershipClauses, client.NodeMembershipClause{
				Label:    clause.Label.ValueString(),
				Operator: clause.Operator.ValueString(),
				Value:    clause.Value.ValueString(),
				Negate:   clause.Negate.ValueBool(),
			})
		}
		nodeMembershipRule, err := s.client.CreateNodeMembershipRule(plan.SpaceID.ValueString(),
			plan.RoomID.ValueString(),
			rule.Action.ValueString(),
			rule.Description.ValueString(),
			nodeMembershipClauses)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Node Membership Rule",
				"err: "+err.Error(),
			)
			return
		}
		plan.Rules[i].ID = types.StringValue(nodeMembershipRule.ID.String())
		plan.Rules[i].Action = types.StringValue(nodeMembershipRule.Action)
		plan.Rules[i].Description = types.StringValue(nodeMembershipRule.Description)
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

	stateNodes := make([]types.String, 0, len(state.NodeNames.Elements()))
	diags = state.NodeNames.ElementsAs(ctx, &stateNodes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var refreshedRoomNodes []string
	for _, stateNode := range stateNodes {
		exist, _ := checkNodeExists(stateNode.ValueString(), nodeRoomMember, false)
		if exist {
			refreshedRoomNodes = append(refreshedRoomNodes, stateNode.ValueString())
		}
	}

	state.NodeNames, _ = types.ListValueFrom(ctx, types.StringType, refreshedRoomNodes)

	nodeMembershipRules, err := s.client.ListNodeMembershipRules(state.SpaceID.ValueString(), state.RoomID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Node Room Membership Rules",
			fmt.Sprintf("Could not read node room membership rules for space_id/room_id: %s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), err.Error()),
		)
	}

	var refreshedNodeMembershipRules []nodeRoomMembershipRule

	for _, rule := range state.Rules {
		var ruleExist bool
		for _, currentRule := range nodeMembershipRules {
			if rule.ID.ValueString() == currentRule.ID.String() {
				ruleExist = true
				break
			}
		}
		if ruleExist {
			nodeMembershipRule, err := s.client.GetNodeMembershipRule(state.SpaceID.ValueString(), state.RoomID.ValueString(), rule.ID.String())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Getting Node Room Membership Rule",
					fmt.Sprintf("Could not read node room membership rule for space_id/room_id/rule_id: %s/%s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), rule.ID.String(), err.Error()),
				)
				return
			}
			var refreshedNodeMembershipRulesClauses []nodeRoomMembershipClause
			for _, clause := range nodeMembershipRule.Clauses {
				refreshedNodeMembershipRulesClauses = append(refreshedNodeMembershipRulesClauses, nodeRoomMembershipClause{
					Label:    types.StringValue(clause.Label),
					Operator: types.StringValue(clause.Operator),
					Value:    types.StringValue(clause.Value),
					Negate:   types.BoolValue(clause.Negate),
				})
			}
			refreshedNodeMembershipRules = append(refreshedNodeMembershipRules, nodeRoomMembershipRule{
				ID:          types.StringValue(nodeMembershipRule.ID.String()),
				Action:      types.StringValue(nodeMembershipRule.Action),
				Description: types.StringValue(nodeMembershipRule.Description),
				Clauses:     refreshedNodeMembershipRulesClauses,
			})
		}
	}

	state.Rules = refreshedNodeMembershipRules

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

	for _, stateRule := range state.Rules {
		exist := checkNodeMembershipRule(stateRule.ID.ValueString(), plan.Rules)
		if !exist {
			err = s.client.DeleteNodeMembershipRule(state.SpaceID.ValueString(), state.RoomID.ValueString(), stateRule.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Deleting Node Membership Rule",
					fmt.Sprintf("Could not delete node membership rule for space_id/room_id/rule_id: %s/%s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), stateRule.ID.ValueString(), err.Error()),
				)
				return
			}
		}
	}

	for i, planRule := range plan.Rules {
		var nodeMembershipRule *client.NodeMembershipRule
		var nodeMembershipClauses []client.NodeMembershipClause

		exist := checkNodeMembershipRule(planRule.ID.ValueString(), state.Rules)
		for _, clause := range planRule.Clauses {
			nodeMembershipClauses = append(nodeMembershipClauses, client.NodeMembershipClause{
				Label:    clause.Label.ValueString(),
				Operator: clause.Operator.ValueString(),
				Value:    clause.Value.ValueString(),
				Negate:   clause.Negate.ValueBool(),
			})
		}
		if exist {
			nodeMembershipRule, err = s.client.UpdateNodeMembershipRule(plan.SpaceID.ValueString(),
				plan.RoomID.ValueString(),
				planRule.ID.ValueString(),
				planRule.Action.ValueString(),
				planRule.Description.ValueString(),
				nodeMembershipClauses)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Node Room Membership Rules",
					fmt.Sprintf("Could not read node room membership rules for space_id/room_id/membership_id: %s/%s/%s err: %v", plan.SpaceID.ValueString(), plan.RoomID.ValueString(), planRule.ID, err.Error()),
				)
				return
			}
		} else {
			nodeMembershipRule, err = s.client.CreateNodeMembershipRule(plan.SpaceID.ValueString(),
				plan.RoomID.ValueString(),
				planRule.Action.ValueString(),
				planRule.Description.ValueString(),
				nodeMembershipClauses)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Creating Node Membership Rule",
					"err: "+err.Error(),
				)
				return
			}
		}
		plan.Rules[i].ID = types.StringValue(nodeMembershipRule.ID.String())
		plan.Rules[i].Action = types.StringValue(nodeMembershipRule.Action)
		plan.Rules[i].Description = types.StringValue(nodeMembershipRule.Description)
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

	for _, rule := range state.Rules {
		err = s.client.DeleteNodeMembershipRule(state.SpaceID.ValueString(), state.RoomID.ValueString(), rule.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Node Membership Rule",
				fmt.Sprintf("Could not delete node membership rule for space_id/room_id/rule_id: %s/%s/%s err: %v", state.SpaceID.ValueString(), state.RoomID.ValueString(), rule.ID.ValueString(), err.Error()),
			)
			return
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

func checkNodeMembershipRule(ruleID string, rules []nodeRoomMembershipRule) bool {
	for _, rule := range rules {
		if ruleID == rule.ID.ValueString() {
			return true
		}
	}
	return false
}

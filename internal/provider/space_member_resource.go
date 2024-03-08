package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource              = &spaceMemberResource{}
	_ resource.ResourceWithConfigure = &spaceMemberResource{}
)

func NewSpaceMemberResource() resource.Resource {
	return &spaceMemberResource{}
}

type spaceMemberResource struct {
	client *client.Client
}

type spaceMemberResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Email   types.String `tfsdk:"email"`
	Role    types.String `tfsdk:"role"`
	SpaceID types.String `tfsdk:"space_id"`
}

func (s *spaceMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_member"
}

func (s *spaceMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Netdata Cloud Space Member resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The Member ID of the space",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Description: "Email of the member",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
						"Invalid email format",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Description: "Role of the member. The community plan can only set the role to `admin`",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9]+$`),
						"Role should be lowercase",
					),
				},
			},
			"space_id": schema.StringAttribute{
				Description: "Space ID of the member",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (s *spaceMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *spaceMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan spaceMemberResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating space member for email: "+plan.Email.ValueString())

	spaceMemberInfo, err := s.client.CreateSpaceMember(plan.SpaceID.ValueString(), plan.Email.ValueString(), plan.Role.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Space Member",
			"err: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(spaceMemberInfo.MemberID)
	plan.Email = types.StringValue(plan.Email.ValueString())
	plan.Role = types.StringValue(spaceMemberInfo.Role)
	plan.SpaceID = types.StringValue(plan.SpaceID.ValueString())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *spaceMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state spaceMemberResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading space member for space_id/id: "+state.Email.ValueString()+"/"+state.ID.ValueString())

	spaceMemberInfo, err := s.client.GetSpaceMemberID(state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Space Member",
			fmt.Sprintf("Could not read space member for space_id/space_member_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}

	state.ID = types.StringValue(spaceMemberInfo.MemberID)
	state.Email = types.StringValue(spaceMemberInfo.Email)
	state.Role = types.StringValue(spaceMemberInfo.Role)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *spaceMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan spaceMemberResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.UpdateSpaceMemberRoleByID(plan.SpaceID.ValueString(), plan.ID.ValueString(), plan.Role.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Space Member Role",
			fmt.Sprintf("Could not update space member for space_id/space_member_id: %s/%s err: %v", plan.SpaceID.ValueString(), plan.ID.ValueString(), err.Error()),
		)
		return
	}

	spaceMemberInfo, err := s.client.GetSpaceMemberID(plan.SpaceID.ValueString(), plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Space Member",
			fmt.Sprintf("Could not read space member for space_id/space_member_id: %s/%s err: %v", plan.SpaceID.ValueString(), plan.ID.ValueString(), err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(spaceMemberInfo.MemberID)
	plan.Email = types.StringValue(spaceMemberInfo.Email)
	plan.Role = types.StringValue(spaceMemberInfo.Role)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *spaceMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state spaceMemberResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteSpaceMember(state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Space Member",
			fmt.Sprintf("Could not delete space member for space_id/space_member_id: %s/%s err: %v", state.SpaceID.ValueString(), state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (s *spaceMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

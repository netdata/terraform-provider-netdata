package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netdata/terraform-provider-netdata/internal/client"
)

var (
	_ resource.Resource              = &spaceResource{}
	_ resource.ResourceWithConfigure = &spaceResource{}
)

func NewSpaceResource() resource.Resource {
	return &spaceResource{}
}

type spaceResource struct {
	client *client.Client
}

type spaceResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ClaimToken  types.String `tfsdk:"claim_token"`
}

func (s *spaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (s *spaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Netdata Cloud Space resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the space",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the space",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(5),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the space",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"claim_token": schema.StringAttribute{
				Description: "The claim token of the space",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (s *spaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (s *spaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan spaceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating space: "+plan.Name.ValueString())

	spaceInfo, err := s.client.CreateSpace(plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Space",
			"err: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(spaceInfo.ID)
	plan.Name = types.StringValue(spaceInfo.Name)
	plan.Description = types.StringValue(spaceInfo.Description)

	tflog.Info(ctx, "Creating Claim Token for Space ID: "+spaceInfo.ID)

	claimToken, err := s.client.GetSpaceClaimToken(spaceInfo.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Claim Token",
			"Could Not Create Claim Token for Space ID: "+spaceInfo.ID+": err: "+err.Error(),
		)
		return
	}

	plan.ClaimToken = types.StringValue(*claimToken)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *spaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state spaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceInfo, err := s.client.GetSpaceByID(state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Getting Space",
			"Could Not Read Space ID: "+state.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(spaceInfo.Name)
	state.Description = types.StringValue(spaceInfo.Description)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *spaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan spaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.UpdateSpaceByID(plan.ID.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Space",
			"Could Not Update Space ID: "+plan.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}

	spaceInfo, err := s.client.GetSpaceByID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Space",
			"Could Not Read Space ID: "+plan.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(spaceInfo.ID)
	plan.Name = types.StringValue(spaceInfo.Name)
	plan.Description = types.StringValue(spaceInfo.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *spaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state spaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteSpaceByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Space",
			"Could Not Delete Space ID: "+state.ID.ValueString()+": err: "+err.Error(),
		)
		return
	}
}

func (s *spaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

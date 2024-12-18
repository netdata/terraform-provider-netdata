package provider

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func commonNotificationSchema(notificationType string) schema.Schema {
	return schema.Schema{
		Description: fmt.Sprintf("Resource for managing centralized notifications for %s. Available only in paid plans.", notificationType),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: fmt.Sprintf("The ID of the %s notification", notificationType),
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: fmt.Sprintf("The name of the %s notification", notificationType),
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: fmt.Sprintf("The enabled status of the %s notification", notificationType),
				Required:    true,
			},
			"space_id": schema.StringAttribute{
				Description: fmt.Sprintf("The ID of the space for the %s notification", notificationType),
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rooms_id": schema.ListAttribute{
				Description: fmt.Sprintf("The list of room IDs to set the %s notification. If the rooms list is null, the %s notification will be applied to `All rooms`", notificationType, notificationType),
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"alarms": schema.StringAttribute{
				Description: fmt.Sprintf("The alarms setting to set the %s notification. Valid values are: `ALARMS_SETTING_ALL`, `ALARMS_SETTING_CRITICAL`, `ALARMS_SETTING_ALL_BUT_UNREACHABLE`, `ALARMS_SETTING_UNREACHABLE`", notificationType),
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(ALARMS_SETTING_ALL|ALARMS_SETTING_CRITICAL|ALARMS_SETTING_ALL_BUT_UNREACHABLE|ALARMS_SETTING_UNREACHABLE)$`),
						"Invalid alarms setting",
					),
				},
			},
		},
	}
}

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
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
			"repeat_notification_min": schema.Int64Attribute{
				Description: fmt.Sprintf("The time interval for the %s notification to be repeated. The interval is presented in minutes and should be between 30 and 1440, or 0 to avoid repetition, which is the default.", notificationType),
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.OneOf(0),
						int64validator.Between(30, 1440),
					),
				},
			},
			"notifications": schema.ListAttribute{
				Description: fmt.Sprintf("The notification options for the %s. Valid values are: `CRITICAL`, `WARNING`, `CLEAR`, `REACHABLE`, `UNREACHABLE`", notificationType),
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"CRITICAL", "WARNING", "CLEAR", "REACHABLE", "UNREACHABLE"}...)),
				},
			},
		},
	}
}

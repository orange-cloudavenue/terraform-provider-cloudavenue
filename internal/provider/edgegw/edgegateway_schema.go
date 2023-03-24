package edgegw

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
)

type edgeGatewaysResourceModel struct {
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
	ID                  types.String   `tfsdk:"id"`
	Tier0VrfID          types.String   `tfsdk:"tier0_vrf_name"`
	Name                types.String   `tfsdk:"name"`
	OwnerType           types.String   `tfsdk:"owner_type"`
	OwnerName           types.String   `tfsdk:"owner_name"`
	Description         types.String   `tfsdk:"description"`
	EnableLoadBalancing types.Bool     `tfsdk:"lb_enabled"`
}

type edgeGatewayDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Tier0VrfID          types.String `tfsdk:"tier0_vrf_name"`
	Name                types.String `tfsdk:"name"`
	OwnerType           types.String `tfsdk:"owner_type"`
	OwnerName           types.String `tfsdk:"owner_name"`
	Description         types.String `tfsdk:"description"`
	EnableLoadBalancing types.Bool   `tfsdk:"lb_enabled"`
}

/*
edgegwSchema

This function is used to create the schema for the edgegateway resource and datasource.
*/
func edgegwSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Edge Gateway ",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: " resource allows you to create and delete Edge Gateways in CloudAvenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: " data source allows you to show the details of an Edge Gateways in CloudAvenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"timeouts": &superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Delete: true,
					Update: true,
				},
			},
			"id": &superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"name": &superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"tier0_vrf_name": &superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Tier0 VRF to which the Edge Gateway is attached.",
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: helpers.ForceNewDescription,
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"owner_type": &superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of the owner of the Edge Gateway (vdc|vdc-group).",
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: helpers.ForceNewDescription,
					Required:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(vdc|vdc-group)$`),
							"must be vdc or vdc-group",
						),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"owner_name": &superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the owner of the Edge Gateway.",
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: helpers.ForceNewDescription,
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"description": &superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the owner of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"lb_enabled": &superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Load Balancing state on the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: " Default to `true`.",
					Optional:            true,
					Default:             booldefault.StaticBool(true),
				},
			},
		},
	}
}

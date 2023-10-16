package edgegw

import (
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

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
			MarkdownDescription: "resource allows you to create and delete Edge Gateways in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allows you to show the details of an Edge Gateways in Cloud Avenue.",
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
			"id": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": &superschema.SuperStringAttribute{
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
			"tier0_vrf_name": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Tier-0 VRF to which the Edge Gateway is attached.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"owner_type": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of the Edge Gateway owner.",
					Validators: []validator.String{
						stringvalidator.OneOf("vdc", "vdc-group"),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"owner_name": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway owner.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"description": &superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"bandwidth": &superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The bandwidth in Mbps of the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "If no value is not specified, the bandwidth is automatically calculated based on the remaining bandwidth of the Tier-0 VRF.",
				},
			},
			"lb_enabled": &superschema.SuperBoolAttribute{
				Deprecated: &superschema.Deprecated{
					DeprecationMessage:                "Remove the lb_enabled attribute configuration and the attribute will be removed in the version 0.16.0 of the provider. This field have does not work and will be replaced soon by a new resource.",
					ComputeMarkdownDeprecationMessage: true,
					Removed:                           true,
					FromAttributeName:                 "lb_enabled",
					TargetRelease:                     "v0.16.0",
					LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/8",
					LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/567",
				},
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Load Balancing state on the Edge Gateway.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
				},
			},
		},
	}
}

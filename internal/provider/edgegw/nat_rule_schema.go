package edgegw

import (
	"context"

	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func natRuleSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_edgegateway_nat_rule` resource allows you to manage EdgeGateway NAT rule. To change the source IP address from a private to a public IP address, you create a source NAT (SNAT) rule. To change the destination IP address from a public to a private IP address, you create a destination NAT (DNAT) rule.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_edgegateway_nat_rule` data source allows you to retrieve informations about an EdgeGateway NAT rule.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the Nat Rule.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Name of the Nat Rule.",
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway",
					Computed:            true,
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Name of the Edge Gateway.",
					Computed:            true,
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A description of the NAT rule",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Enable or Disable the Nat Rule.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Default:  booldefault.StaticBool(true),
					Optional: true,
				},
			},
			"rule_type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Nat Rule type.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						fstringvalidator.OneOfWithDescription(
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "DNAT",
								Description: "Rule translates the external IP to an internal IP and is used for inbound traffic.",
							},
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "NO_DNAT",
								Description: "Prevents external IP translation.",
							},
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "SNAT",
								Description: "Translates an internal IP to an external IP and is used for outbound traffic.",
							},
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "NO_SNAT",
								Description: "Prevents internal IP translation.",
							},
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "REFLEXIVE",
								Description: "This translates an internal IP to an external IP and vice versa.",
							},
						),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"external_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The external address for the NAT Rule. This must be supplied as a single IP or Network CIDR. For a DNAT rule, this is the external facing IP Address for incoming traffic. For an SNAT rule, this is the external facing IP Address for outgoing traffic. These IPs are typically allocated/suballocated IP Addresses on the Edge Gateway. For a REFLEXIVE rule, these are the external facing IPs.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("rule_type"), []attr.Value{types.StringValue("NO_SNAT")}),
						fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("rule_type"), []attr.Value{types.StringValue("DNAT"), types.StringValue("SNAT"), types.StringValue("NO_DNAT"), types.StringValue("REFLEXIVE")}),
					},
					// TODO - Validator of IP or IP/CIDR
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"internal_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The internal address for the NAT Rule. This must be supplied as a single IP or Network CIDR. For a DNAT rule, this is the internal IP address for incoming traffic. For an SNAT rule, this is the internal IP Address for outgoing traffic. For a REFLEXIVE rule, these are the internal IPs. These IPs are typically the Private IPs that are allocated to workloads.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("rule_type"), []attr.Value{types.StringValue("NO_DNAT")}),
						fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("rule_type"), []attr.Value{types.StringValue("DNAT"), types.StringValue("NO_SNAT"), types.StringValue("REFLEXIVE")}),
					},
					// TODO - Validator of IP or IP/CIDR
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"app_port_profile_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Application Port Profile ID to which the rule applies.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.PrefixContains(urn.AppPortProfile.String()),
						fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("rule_type"), []attr.Value{types.StringValue("SNAT"), types.StringValue("NO_SNAT"), types.StringValue("NO_DNAT"), types.StringValue("REFLEXIVE")}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"dnat_external_port": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "This represents the external port number or port range when doing DNAT port forwarding from external to internal. If not specify, all ports are translated",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("rule_type"), []attr.Value{types.StringValue("SNAT"), types.StringValue("NO_SNAT"), types.StringValue("REFLEXIVE")}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"snat_destination_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The destination addresses to match in the SNAT Rule. This must be supplied as a single IP or Network CIDR. Providing no value for this field results in match with ANY destination network.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("rule_type"), []attr.Value{types.StringValue("DNAT"), types.StringValue("NO_DNAT"), types.StringValue("REFLEXIVE")}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			// Not Supported on CloudAvenue
			// "logging": superschema.SuperBoolAttribute{
			// 	Common: &schemaR.BoolAttribute{
			// 		MarkdownDescription: "Enable to have the address translation performed by this rule logged ",
			// 	},
			// 	Resource: &schemaR.BoolAttribute{
			// 		Optional: true,
			// 		Default:  booldefault.StaticBool(true),
			// 	},
			// 	DataSource: &schemaD.BoolAttribute{
			// 		Computed: true,
			// 	},
			// },
			"priority": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "If an address has multiple NAT rule, you can assign these rule different priorities to determine the order in which they are applied. A lower value means a higher priority for this rule.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Default:  int64default.StaticInt64(0),
					Optional: true,
				},
			},
			"firewall_match": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "You can set a firewall match rule to determine how firewall is applied during NAT.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.OneOfWithDescription(
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "MATCH_INTERNAL_ADDRESS",
								Description: "Applies firewall rule to the internal address of a NAT rule.",
							},
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "MATCH_EXTERNAL_ADDRESS",
								Description: "Applies firewall rule to the external address of a NAT rule.",
							},
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "BYPASS",
								Description: "Skip applying firewall rule to NAT rule.",
							},
						),
					},
				},
			},
		},
	}
}

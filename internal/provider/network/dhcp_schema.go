package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fsetvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/setvalidator"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
)

func dhcpSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `network_dhcp` resource allows you to manage DHCP servers on Org Network.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `network_dhcp` data source allows you to retrieve information about an existing DHCP server on Org Network.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the DHCP server.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"org_network_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "The ID of the network.",
					Validators: []validator.String{
						fstringvalidator.IsURN(),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"mode": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The mode of the DHCP server.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Default:  stringdefault.StaticString("EDGE"),
					Validators: []validator.String{
						fstringvalidator.OneOfWithDescription(
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "EDGE",
								Description: "The Edge's DHCP service is used to obtain DHCP IP addresses.",
							},
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "NETWORK",
								Description: "A new DHCP service directly associated with this network is used to obtain DHCP IP addresses. Use Network Mode if the network is isolated or if you plan to detach this network from the Edge",
							},
							fstringvalidator.OneOfWithDescriptionValues{
								Value:       "RELAY",
								Description: "DHCP messages are relayed from virtual machines to designated DHCP servers in your physical DHCP infrastructure.",
							},
						),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"pools": superschema.SetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "IP ranges used for DHCP pool allocation in the network",
					Computed:            true,
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
					Validators: []validator.Set{
						fsetvalidator.NullIfAttributeIsOneOf(path.MatchRoot("mode"), []attr.Value{types.StringValue("RELAY")}),
						fsetvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("mode"), []attr.Value{types.StringValue("EDGE"), types.StringValue("NETWORK")}),
					},
				},
				Attributes: superschema.Attributes{
					"start_address": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The start address of the DHCP pool IP range.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								fstringvalidator.IsIP(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"end_address": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The end address of the DHCP pool IP range.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								fstringvalidator.IsIP(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"dns_servers": superschema.ListAttribute{
				Common: &schemaR.ListAttribute{
					MarkdownDescription: "The DNS server IPs to be assigned by this DHCP service.",
					ElementType:         types.StringType,
				},
				Resource: &schemaR.ListAttribute{
					Optional: true,
					Validators: []validator.List{
						listvalidator.SizeAtMost(2),
					},
				},
				DataSource: &schemaD.ListAttribute{
					Computed: true,
				},
			},
			"lease_time": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The lease time in seconds for the DHCP service.",
					Computed:            true,
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Default:  int64default.StaticInt64(86400),
					Validators: []validator.Int64{
						int64validator.AtLeast(60),
					},
				},
			},
			"listener_ip_address": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The IP address of the DHCP listener.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
						fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("mode"), []attr.Value{types.StringValue("RELAY"), types.StringValue("EDGE")}),
					},
					PlanModifiers: []planmodifier.String{
						// API still does not allow to change IP address in 10.4.0, but the error is human
						// readable and it might allow changing in future.
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
				},
			},
		},
	}
}

package network

import (
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

type networkSchemaOpts func(*networkSchemaParams)

type networkSchemaParams struct {
	resource bool

	typeNetwork Type
}

type Type int

const (
	ISOLATED   Type = iota
	NAT_ROUTED      //nolint:revive,stylecheck
	ISOLATEDVAPP
	ROUTEDVAPP
)

// Set bool to true to create a schema for a routed network.
func SetRouted() networkSchemaOpts {
	return func(params *networkSchemaParams) {
		params.typeNetwork = NAT_ROUTED
	}
}

// Set bool to true to create a schema for an isolated network.
func SetIsolated() networkSchemaOpts {
	return func(params *networkSchemaParams) {
		params.typeNetwork = ISOLATED
	}
}

// Set bool to true to create a schema for an isolated vApp network.
func SetIsolatedVapp() networkSchemaOpts {
	return func(params *networkSchemaParams) {
		params.typeNetwork = ISOLATEDVAPP
	}
}

func SetRoutedVapp() networkSchemaOpts {
	return func(params *networkSchemaParams) {
		params.typeNetwork = ROUTEDVAPP
	}
}

/*
networkSchema

This function is used to create the schema for the catalog resource and datasource.
Default is to create a resource schema.  If you want to create a datasource schema
you must pass in the withDataSource() option.
*/
func GetSchema(opts ...networkSchemaOpts) superschema.Schema {
	params := &networkSchemaParams{}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(params)
		}
	} else {
		params.resource = true
	}

	// _schema := superschema.Schema{}
	_schema := superschema.Schema{
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the network.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the network. This value must be unique within the `VDC` or `VDC Group` that owns the network.",
					Required:            true,
				},
			},
			"description": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A description of the network.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"gateway": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The gateway IP address for the network. This value define also the network IP range with the prefix length.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"prefix_length": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The prefix length for the network. This value must be a valid prefix length for the network IP range. (e.g. /24 for netmask 255.255.255.0)",
				},
				Resource: &schemaR.Int64Attribute{
					Required: true,
					Validators: []validator.Int64{
						int64validator.Between(1, 32),
					},
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"dns1": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The primary DNS server IP address for the network.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"dns2": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The secondary DNS server IP address for the network.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"dns_suffix": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The DNS suffix for the network.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"static_ip_pool": superschema.SetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "A set of static IP pools to be used for this network.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"start_address": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The start address of the IP pool. This value must be a valid IP address in the network IP range.",
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
							MarkdownDescription: "The end address of the IP pool. This value must be a valid IP address in the network IP range.",
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
		},
	}

	switch params.typeNetwork {
	case NAT_ROUTED:
		// Add routed network specific attributes to the schema
		_schema.Resource.MarkdownDescription = "Provides a Cloud Avenue vDC routed Network. This can be used to create, modify, and delete VDC routed networks."
		_schema.DataSource.MarkdownDescription = "Provides a Cloud Avenue vDC routed Network data source to read data or reference existing network"
		_schema.Attributes["edge_gateway_id"] = superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The ID of the edge gateway in which the routed network should be located.",
				Optional:            true,
				Computed:            true,
			},
			Resource: &schemaR.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
					stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_id"), path.MatchRoot("edge_gateway_name")),
				},
			},
		}
		_schema.Attributes["edge_gateway_name"] = superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The name of the edge gateway in which the routed network should be located.",
				Computed:            true,
			},
			Resource: &schemaR.StringAttribute{
				MarkdownDescription: "The name of the edge gateway in which the routed network should be located.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		}
		_schema.Attributes["interface_type"] = superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "An interface for the network",
				Computed:            true,
			},
			Resource: &schemaR.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("INTERNAL"),
				Validators: []validator.String{
					stringvalidator.OneOf("INTERNAL", "SUBINTERFACE", "DISTRIBUTED"),
				},
			},
		}

	case ISOLATED:
		// Add isolated network specific attributes to the schema
		_schema.Resource.MarkdownDescription = "Provides a Cloud Avenue VDC isolated Network. This can be used to create, modify, and delete VDC isolated networks."
		_schema.DataSource.MarkdownDescription = "Provides a Cloud Avenue VDC isolated Network data source to read data or reference existing network."
		_schema.Attributes["vdc"] = vdc.SuperSchema()

		// * Add deprecated resource and data source
		_schema.Resource.Deprecated = superschema.DeprecatedResource{
			DeprecationMessage:                "The `cloudavenue_network_isolated` resource is deprecated. Please use the `cloudavenue_vdc_network_isolated` resource instead.",
			ComputeMarkdownDeprecationMessage: true,
			Renamed:                           true,
			TargetResourceName:                "cloudavenue_vdc_network_isolated",
			TargetRelease:                     "v0.32.0",
			LinkToMigrationGuide:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/network_isolated#how-to-migrate-existing-resources",
			LinkToNewResourceDoc:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc_network_isolated",
			LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/20",
			LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/880",
		}

		_schema.DataSource.Deprecated = superschema.DeprecatedResource{
			DeprecationMessage:                "The `cloudavenue_network_isolated` data source is deprecated. Please use the `cloudavenue_vdc_network_isolated` data source instead.",
			ComputeMarkdownDeprecationMessage: true,
			Renamed:                           true,
			TargetResourceName:                "cloudavenue_vdc_network_isolated",
			TargetRelease:                     "v0.32.0",
			LinkToNewResourceDoc:              "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/data-sources/vdc_network_isolated",
			LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/20",
			LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/880",
		}

	case ISOLATEDVAPP:
		// Add isolated vApp network specific attributes to the schema
		delete(_schema.Attributes, "prefix_length")
		_schema.Resource.MarkdownDescription = "Provides a Cloud Avenue isolated vAPP Network resource. This can be used to create, modify, and delete isolated vAPP Network."
		_schema.DataSource.MarkdownDescription = "Provides a Cloud Avenue isolated vAPP Network data source to read data or reference existing network."
		_schema.Attributes["netmask"] = superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The netmask of the network.",
				Computed:            true,
			},
			Resource: &schemaR.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("255.255.255.0"),
				Validators: []validator.String{
					fstringvalidator.IsNetmask(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		}
		_schema.Attributes["vdc"] = vdc.SuperSchema()
		_schema.Attributes["vapp_id"] = vapp.SuperSchema()["vapp_id"]
		_schema.Attributes["vapp_name"] = vapp.SuperSchema()["vapp_name"]
		_schema.Attributes["guest_vlan_allowed"] = superschema.BoolAttribute{
			DataSource: &schemaD.BoolAttribute{
				MarkdownDescription: "Return True if Network allows guest VLAN.",
				Computed:            true,
			},
		}
		_schema.Attributes["retain_ip_mac_enabled"] = superschema.BoolAttribute{
			DataSource: &schemaD.BoolAttribute{
				MarkdownDescription: "Return network resources such as IP/MAC of router will be retained across deployments.",
				Computed:            true,
			},
		}

	case ROUTEDVAPP:
		// Add routed vApp network specific attributes to the schema
		delete(_schema.Attributes, "description")
		delete(_schema.Attributes, "prefix_length")
		delete(_schema.Attributes, "gateway")
		delete(_schema.Attributes, "dns1")
		delete(_schema.Attributes, "dns2")
		delete(_schema.Attributes, "dns_suffix")
		delete(_schema.Attributes, "static_ip_pool")
		delete(_schema.Attributes, "name")
		_schema.Resource.MarkdownDescription = "Provides a Cloud Avenue routed vAPP Org Network resource. This can be used to create, modify, and delete routed vAPP Network."
		_schema.DataSource.MarkdownDescription = "Provides a Cloud Avenue routed vAPP Org Network data source to read data or reference existing network."
		_schema.Attributes["vdc"] = vdc.SuperSchema()
		_schema.Attributes["vapp_id"] = vapp.SuperSchema()["vapp_id"]
		_schema.Attributes["vapp_name"] = vapp.SuperSchema()["vapp_name"]
		_schema.Attributes["network_name"] = superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "Organization network name to which vApp network is connected to.",
				Required:            true,
			},
			Resource: &schemaR.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		}
	}
	return _schema
}

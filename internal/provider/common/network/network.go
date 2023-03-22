package network

import (
	fstringplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/stringplanmodifier"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

type networkSchemaOpts func(*networkSchemaParams)

type networkSchemaParams struct {
	resource bool

	typeNetwork networkType
}

type networkType int

const (
	Isolated networkType = iota
	Routed
	IsolatedVapp
	RoutedVapp
)

// TODO - plan to use this to create a datasource schema
// func withDataSource() networkSchemaOpts {
//	return func(params *networkSchemaParams) {
//		params.datasource = true
//	}
//}

// Set bool to true to create a schema for a routed network.
func SetRouted() networkSchemaOpts {
	return func(params *networkSchemaParams) {
		params.typeNetwork = Routed
	}
}

// Set bool to true to create a schema for an isolated network.
func SetIsolated() networkSchemaOpts {
	return func(params *networkSchemaParams) {
		params.typeNetwork = Isolated
	}
}

// Set bool to true to create a schema for an isolated vApp network.
func SetIsolatedVapp() networkSchemaOpts {
	return func(params *networkSchemaParams) {
		params.typeNetwork = IsolatedVapp
	}
}

func SetRoutedVapp() networkSchemaOpts {
	return func(params *networkSchemaParams) {
		params.typeNetwork = RoutedVapp
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

	_schema := superschema.Schema{}
	_schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The ID of the network. This is a generated value and cannot be specified during creation. This value is used to identify the network in other resources.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name of the network. This value must be unique within the `VDC` or `VDC Group` that owns the network.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "A description of the network.",
		},
		"gateway": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "(Force replacement) The gateway IP address for the network. This value define also the network IP range with the prefix length.",
			Validators: []validator.String{
				fstringvalidator.IsIP(),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"prefix_length": schema.Int64Attribute{
			Required:            true,
			MarkdownDescription: "(Force replacement) The prefix length for the network. This value must be a valid prefix length for the network IP range.(e.g. /24 for netmask 255.255.255.0)",
			Validators: []validator.Int64{
				int64validator.Between(1, 32),
			},
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"dns1": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The primary DNS server IP address for the network.",
			Validators: []validator.String{
				fstringvalidator.IsIP(),
			},
		},
		"dns2": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The secondary DNS server IP address for the network.",
			Validators: []validator.String{
				fstringvalidator.IsIP(),
			},
		},
		"dns_suffix": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The DNS suffix for the network.",
		},
		"static_ip_pool": schema.SetNestedAttribute{
			Optional:            true,
			MarkdownDescription: "A set of static IP pools to be used for this network.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"start_address": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The start address of the IP pool. This value must be a valid IP address in the network IP range.",
						Validators: []validator.String{
							fstringvalidator.IsIP(),
						},
					},
					"end_address": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The end address of the IP pool. This value must be a valid IP address in the network IP range.",
						Validators: []validator.String{
							fstringvalidator.IsIP(),
						},
					},
				},
			},
		},
	}

	switch params.typeNetwork {
	case Routed:
		// Add routed network specific attributes to the schema
		_schema.MarkdownDescription = "Provides a Cloud Avenue Org VDC routed Network. This can be used to create, modify, and delete routed VDC networks."
		_schema.Attributes["edge_gateway_id"] = schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
			MarkdownDescription: "Edge gateway ID in which Routed network should be located.",
		}
		_schema.Attributes["interface_type"] = schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Optional interface type (only for NSX-V networks). One of `INTERNAL` (default), `DISTRIBUTED`, `SUBINTERFACE`",
			PlanModifiers: []planmodifier.String{
				fstringplanmodifier.SetDefault("INTERNAL"),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("INTERNAL", "SUBINTERFACE", "DISTRIBUTED"),
			},
		}

	case Isolated:
		// Add isolated network specific attributes to the schema
		_schema.MarkdownDescription = "Provides a Cloud Avenue Org VDC isolated Network. This can be used to create, modify, and delete isolated VDC networks"
		_schema.Attributes["vdc"] = vdc.Schema()

	case IsolatedVapp:
		// Add isolated vApp network specific attributes to the schema
		delete(_schema.Attributes, "prefix_length")
		_schema.Attributes["netmask"] = schema.StringAttribute{
			MarkdownDescription: "(ForceNew) The netmask for the network. Default is `255.255.255.0`",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				fstringvalidator.IsNetmask(),
			},
			PlanModifiers: []planmodifier.String{
				fstringplanmodifier.SetDefault("255.255.255.0"),
				stringplanmodifier.RequiresReplace(),
			},
		}

	case RoutedVapp:
		// Add routed vApp network specific attributes to the schema
		delete(_schema.Attributes, "description")
		delete(_schema.Attributes, "prefix_length")
		delete(_schema.Attributes, "gateway")
		delete(_schema.Attributes, "dns1")
		delete(_schema.Attributes, "dns2")
		delete(_schema.Attributes, "dns_suffix")
		delete(_schema.Attributes, "static_ip_pool")
		delete(_schema.Attributes, "name")
		_schema.Attributes["network_name"] = schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			MarkdownDescription: "Organization network name to which vApp network is connected to.",
		}
	}

	return _schema
}

package vapp

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

func isolatedNetworkSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Provides capability to attach an existing Org VDC Network to a vApp and toggle network features.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a Cloud Avenue isolated vAPP Network data source to read data or reference existing network.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the isolated network.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the network. This value must be unique within the `VDC` or `VDC Group` that owns the network.",
					Required:            true,
				},
			},
			"description": superschema.SuperStringAttribute{
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
			"gateway": superschema.SuperStringAttribute{
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
			"netmask": superschema.SuperStringAttribute{
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
			},
			"dns1": superschema.SuperStringAttribute{
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
			"dns2": superschema.SuperStringAttribute{
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
			"dns_suffix": superschema.SuperStringAttribute{
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
			"static_ip_pool": superschema.SuperSetNestedAttributeOf[isolatedNetworkModelStaticIPPool]{
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
					"start_address": superschema.SuperStringAttribute{
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
					"end_address": superschema.SuperStringAttribute{
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
			"vdc": vdc.SuperSchemaSuperType(),
			"vapp_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "ID of the vApp.",
					Computed:            true,
					Optional:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
					},
				},
			},
			"vapp_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Name of the vApp.",
					Computed:            true,
					Optional:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_id"), path.MatchRoot("vapp_name")),
					},
				},
			},
			"guest_vlan_allowed": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "True if Network allows guest VLAN.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Default:  booldefault.StaticBool(false),
				},
			},
			"retain_ip_mac_enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Specifies whether the network resources such as IP/MAC of router will be retained across deployments.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Default:  booldefault.StaticBool(false),
				},
			},
		},
	}
}

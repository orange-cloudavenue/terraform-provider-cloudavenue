package vm

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
)

// DiskSuperSchema returns the super schema.
func DiskSuperSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The virtual machine (vm) resource allows you to manage a virtual machine in the CloudAvenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the VM.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"vdc":       vdc.SuperSchema(),
			"vapp_id":   vapp.SuperSchema()["vapp_id"],
			"vapp_name": vapp.SuperSchema()["vapp_name"],
			"vm_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the VM where the disk will be attached.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.Any(
							fstringvalidator.RequireIfAttributeIsOneOf(
								path.MatchRoot("is_detachable"),
								[]attr.Value{
									types.BoolValue(false),
								},
							),
							stringvalidator.ExactlyOneOf(path.MatchRoot("vm_id")),
						),
					},
				},
			},
			"vm_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the VM where the disk will be attached.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.Any(
							fstringvalidator.RequireIfAttributeIsOneOf(
								path.MatchRoot("is_detachable"),
								[]attr.Value{
									types.BoolValue(false),
								},
							),
							stringvalidator.ExactlyOneOf(path.MatchRoot("vm_name")),
						),
					},
				},
			},
			"is_detachable": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "If set to `true`, the disk will be detached from the VM. If set to `false`, the disk will be attached to the VM.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Default:  booldefault.StaticBool(false),
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
						boolplanmodifier.RequiresReplace(),
					},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the disk.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						// This attribute is required if the disk is detachable.
						fstringvalidator.RequireIfAttributeIsOneOf(
							path.MatchRoot("is_detachable"),
							[]attr.Value{
								types.BoolValue(true),
							},
						),
					},
				},
			},
			"storage_profile": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the storage profile. If not set, the default storage profile will be used.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(diskparams.StorageProfileValues...),
					},
				},
			},
			"size_in_mb": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The size of the disk in MB.",
				},
				Resource: &schemaR.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
				},
			},
			"bus_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of disk controller.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Default:  stringdefault.StaticString(diskparams.BusTypeSCSI.Name()),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{ // TODO: Miss the validation of the bus type if disk is detachable or not.
						stringvalidator.OneOf(diskparams.ListOfBusTypes...),
					},
				},
			},
			"bus_number": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The bus number of the disk controller.",
				},
				Resource: &schemaR.Int64Attribute{
					MarkdownDescription: "If the disk is attached to a VM and this attribute is not set, the disk will be attached to the first available bus.",
					Optional:            true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
					Validators: []validator.Int64{
						int64validator.Between(0, 3),
					},
				},
			},
			"unit_number": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The unit number of the disk controller.",
				},
				Resource: &schemaR.Int64Attribute{
					MarkdownDescription: "If the disk is attached to a VM and this attribute is not set, the disk will be attached to the first available unit.",
					Optional:            true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
					Validators: []validator.Int64{
						int64validator.Between(0, 15),
					},
				},
			},
		},
	}
}

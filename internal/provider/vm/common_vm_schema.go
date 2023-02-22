package vm

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/boolpm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

func vmSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The vapp_vm resource allows you to manage a ...",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the vapp_vm.",
				Computed:            true,
			},
			"vapp_name": schema.StringAttribute{
				MarkdownDescription: "The vApp this VM belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_name": schema.StringAttribute{
				MarkdownDescription: "The name of the VM. Unique within the vApp.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"computer_name": schema.StringAttribute{
				MarkdownDescription: "Computer name to assign to this virtual machine",
				Optional:            true,
				Computed:            true,
			},
			"vdc": schema.StringAttribute{
				MarkdownDescription: "The name of the VDC this VM belongs to. If not specified, VDC define in provider will be used.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vapp_template_id": schema.StringAttribute{
				MarkdownDescription: "The URN of the vApp template to use for this VM.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_name_in_template": schema.StringAttribute{
				MarkdownDescription: "The name of the VM in vApp Template to use. In cases when vApp template has more than one VM",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the VM.",
				Optional:            true,
				Computed:            true,
			},
			"resource": schema.SingleNestedAttribute{
				MarkdownDescription: "Resource configuration for the VM.",
				Optional:            true,
				Computed:            true,
			},
			"href": schema.StringAttribute{
				MarkdownDescription: "VM Hyper Reference",
				Computed:            true,
			},
			"accept_all_eulas": schema.BoolAttribute{
				MarkdownDescription: "Automatically accept EULA if OVA has it. Default is `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolpm.SetDefault(false),
				},
			},
			"power_on": schema.BoolAttribute{
				MarkdownDescription: "A boolean value stating if this VM should be powered on. Default is `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolpm.SetDefault(true),
				},
			},
			"storage_profile": schema.StringAttribute{
				MarkdownDescription: "Storage profile to override the default one.",
				Optional:            true,
				Computed:            true,
			},
			"os_type": schema.StringAttribute{
				MarkdownDescription: "Operating System type.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(vm.GetAllOsTypes()...),
				},
			},
			"boot_image_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the boot image to use for this VM.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_template_id"), path.MatchRoot("boot_image_id")),
				},
			},
			"network_dhcp_wait_seconds": schema.Int64Attribute{
				MarkdownDescription: "Optional number of seconds to try and wait for DHCP IP (valid for `network` block only)",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"networks": schema.ListNestedAttribute{
				MarkdownDescription: "A block to define network interface. Multiple can be used.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: vm.NetworkSchema(),
				},
			},
			"disks": schema.ListNestedAttribute{
				MarkdownDescription: "A block to define disk. Multiple can be used.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: vm.DiskSchema(),
				},
			},
			"override_template_disks": schema.SetNestedAttribute{
				MarkdownDescription: "A block to match internal_disk interface in template. Multiple can be used. Disk will be matched by bus_type, bus_number and unit_number.",
				Optional:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: vm.TemplateDiskSchema(),
				},
			},
			"internal_disks": schema.SingleNestedAttribute{
				MarkdownDescription: "A block will show internal disks details",
				Computed:            true,
				Attributes:          vm.InternalDiskSchema(),
			},
			"expose_hardware_virtualization": schema.BoolAttribute{
				MarkdownDescription: "Expose hardware-assisted CPU virtualization to guest OS.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolpm.SetDefault(false),
				},
			},
			"guest_properties": schema.MapAttribute{
				MarkdownDescription: "Key/Value settings for guest properties",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"customization": schema.SingleNestedAttribute{
				MarkdownDescription: "Guest customization block. ",
				Optional:            true,
				Computed:            true,
				Attributes:          vm.CustomizationSchema(),
			},
			"prevent_update_power_off": schema.BoolAttribute{
				MarkdownDescription: "`true` if the update of resource should fail when virtual machine power off needed. Default is `false`.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolpm.SetDefault(false),
				},
			},
			"sizing_policy_id": schema.StringAttribute{
				MarkdownDescription: "VM placement policy ID. Has to be assigned to Org VDC.",
				Optional:            true,
				Computed:            true, // As it can get populated automatically by VDC default policy
			},
			"placement_policy_id": schema.StringAttribute{
				MarkdownDescription: "VM placement policy ID. Has to be assigned to Org VDC.",
				Optional:            true,
				Computed:            true, // As it can get populated automatically by VDC default policy
			},
			"status_code": schema.Int64Attribute{
				MarkdownDescription: "Shows the status code of the VM",
				Computed:            true,
			},
			"status_text": schema.StringAttribute{
				MarkdownDescription: "Shows the status text of the VM",
				Computed:            true,
			},
		},
	}
}

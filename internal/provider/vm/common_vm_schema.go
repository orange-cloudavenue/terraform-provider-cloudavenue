package vm

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/boolpm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

// TODO
// - Add vapp_id

func vmSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The virtual machine (vm) resource allows you to manage a virtual machine in the CloudAvenue.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the VM.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vapp_name": schema.StringAttribute{
				MarkdownDescription: "The vApp this VM belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the VM. Unique within the vApp.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`), "must be between 1 and 80 characters long and can contain only letters, numbers and hyphen. It must not contain only digits."),
				},
			},
			"computer_name": schema.StringAttribute{
				MarkdownDescription: "Computer name to assign to this virtual machine. Default is `vm_name`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`), "must be between 1 and 80 characters long and can contain only letters, numbers and hyphen. It must not contain only digits."),
				},
			},
			"vdc": schema.StringAttribute{
				MarkdownDescription: "The name of the VDC this VM belongs to. If not specified, VDC define in provider will be used.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vapp_template_id": schema.StringAttribute{
				MarkdownDescription: "The URN of the vApp template to use for this VM. Required if `boot_image_id` is not set.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_name_in_template": schema.StringAttribute{
				MarkdownDescription: "The name of the VM in vApp Template to use. In cases when vApp template has more than one VM.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the VM.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource": schema.SingleNestedAttribute{
				MarkdownDescription: "Resource configuration for the VM.",
				Optional:            true,
				Computed:            true,
				Attributes:          vm.ResourceSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"href": schema.StringAttribute{
				MarkdownDescription: "VM Hyper Reference",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os_type": schema.StringAttribute{
				MarkdownDescription: "Operating System type.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(vm.GetAllOsTypes()...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"boot_image_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the boot image to use for this VM. Required if `vapp_template_id` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_template_id")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: vm.NetworkSchema(),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"disks": schema.SetNestedAttribute{
				MarkdownDescription: "A list of disks attached to this VM.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: diskInternalExternalSchema(),
				},
			},
			"expose_hardware_virtualization": schema.BoolAttribute{
				MarkdownDescription: "Expose hardware-assisted CPU virtualization to guest OS. Default is `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolpm.SetDefault(false),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"guest_properties": vm.GuestPropertiesSchema(),
			"customization": schema.SingleNestedAttribute{
				MarkdownDescription: "Guest customization block. ",
				Optional:            true,
				Computed:            true,
				Attributes:          vm.CustomizationSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"placement_policy_id": schema.StringAttribute{
				MarkdownDescription: "VM placement policy ID. Has to be assigned to Org VDC.",
				Optional:            true,
				Computed:            true, // As it can get populated automatically by VDC default policy
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

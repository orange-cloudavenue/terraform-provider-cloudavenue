package vm

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fint64validator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/int64validator"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/storageprofile"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

func vmSuperSchema(_ context.Context) superschema.Schema {
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
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the VM. Unique within the vApp.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`), "Must be between 1 and 80 characters long and can contain only letters, numbers and hyphen. It must not contain only digits."),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"vapp_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The vApp this VM belongs to.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
					},
				},
			},
			"vapp_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The vApp this VM belongs to.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
					},
				},
			},
			"vdc": vdc.SuperSchema(),
			"description": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the VM.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"deploy_os": superschema.SingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"vapp_template_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the vApp template to use for the VM.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							// Validators: []validator.String{
							// 	fstringvalidator.Not(stringvalidator.AlsoRequires(path.MatchRelative().AtName("boot_image_id"))),
							// },
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"vm_name_in_template": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the VM in the vApp template.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.AlsoRequires(path.MatchRelative().AtName("vapp_template_id")),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"boot_image_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the boot image to use for the VM.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							// Validators: []validator.String{
							// 	fstringvalidator.Not(stringvalidator.AlsoRequires(path.MatchRelative().AtName("vapp_template_id"))),
							// },
						},
					},
					"accept_all_eulas": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Automatically accept EULA if OVA has it",
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(true),
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
			"state": superschema.SingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The state of the VM.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"power_on": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether the VM should be powered on or not. `true` means powered on, `false` means powered off.",
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Computed: true,
							// TODO: BUG démarre pas la VM
							Default: booldefault.StaticBool(true),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
					"status": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The status of the VM.",
							Computed:            true,
						},
						// Resource: &schemaR.StringAttribute{
						// 	PlanModifiers: []planmodifier.String{
						// 		stringplanmodifier.UseStateForUnknown(),
						// 	},
						// },
					},
				},
			},
			"resource": superschema.SingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The resource of the VM.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					// Default:  objectdefault.StaticValue(resourceDefaultValue),
				},
				Attributes: map[string]superschema.Attribute{
					"cpus": superschema.Int64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of virtual CPUs to allocate to the VM.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Default:  int64default.StaticInt64(1),
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
							Validators: []validator.Int64{
								int64validator.AtMost(256),
							},
						},
					},
					"cpus_cores": superschema.Int64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of cores per virtual CPU to allocate to the VM.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Default:  int64default.StaticInt64(1),
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
							Validators: []validator.Int64{
								// TODO fix path
							},
						},
					},
					"cpu_hot_add_enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether CPU hot add is enabled or not.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							// TODO: BUG ne fonctionne pas
							Default: booldefault.StaticBool(true),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
					"memory": superschema.Int64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The amount of memory to allocate to the VM, in MB.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Default:  int64default.StaticInt64(1024),
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
							Validators: []validator.Int64{
								fint64validator.ZeroRemainder(4),
							},
						},
					},
					"memory_hot_add_enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether memory hot add is enabled or not.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							// TODO: BUG ne fonctionne pas
							Default: booldefault.StaticBool(true),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
					"networks": superschema.ListNestedAttribute{
						Common: &schemaR.ListNestedAttribute{
							MarkdownDescription: "The networks to attach to the VM.",
							Computed:            true,
						},
						Resource: &schemaR.ListNestedAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.List{
								listplanmodifier.UseStateForUnknown(),
							},
						},
						Attributes: map[string]superschema.Attribute{
							"type": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The type of network to attach to the VM.",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("vapp", "org", "none"),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"ip_allocation_mode": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The IP allocation mode for this network.",
									Computed:            true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Default:  stringdefault.StaticString("DHCP"),
									Validators: []validator.String{
										stringvalidator.OneOf("DHCP", "POOL", "MANUAL", "NONE"),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"name": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The name of the network to attach to the VM.",
								},
								Resource: &schemaR.StringAttribute{
									Optional:   true,
									Validators: []validator.String{
										// fstringvalidator.RequireIfAttributeIsOneOf(
										// 	path.MatchRelative().AtParent().AtName("type"),
										// 	[]attr.Value{
										// 		types.StringValue("vapp"),
										// 		types.StringValue("org"),
										// 	},
										// ),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"ip": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The IP address to assign to this VM on this network.",
								},
								Resource: &schemaR.StringAttribute{
									Optional:   true,
									Validators: []validator.String{
										// stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
										// fstringvalidator.RequireIfAttributeIsOneOf(
										// 	path.MatchRelative().AtParent().AtName("ip_allocation_mode"),
										// 	// path.MatchRelative().AtParent().AtName("ip_allocation_mode"),
										// 	[]attr.Value{
										// 		types.StringValue("MANUAL"),
										// 		types.StringValue("POOL"),
										// 	},
										// ),
										// TODO: error if ip_allocation_mode is not MANUAL or POOL
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"is_primary": superschema.BoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether this network is the primary network for the VM.",
									Computed:            true,
								},
								Resource: &schemaR.BoolAttribute{
									Optional: true,
									Default:  booldefault.StaticBool(false),
								},
							},
							"mac": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The MAC address to assign to this VM on this network.",
									Computed:            true,
								},
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "Autogenerated if not specified.",
									Optional:            true,
									Validators: []validator.String{
										fstringvalidator.IsMacAddress(),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"adapter_type": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The type of vNic to create on this interface.",
									Computed:            true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Default:  stringdefault.StaticString("VMXNET3"),
									Validators: []validator.String{
										stringvalidator.OneOf("VMXNET3", "E1000E", "VMXNET3VRDMA", "SRIOVETHERNETCARD"),
									},
								},
							},
							"connected": superschema.BoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether the network interface is connected or not.",
									Computed:            true,
								},
								Resource: &schemaR.BoolAttribute{
									Optional: true,
									Default:  booldefault.StaticBool(true),
								},
							},
						},
					},
				},
			},
			// End of resource attributes
			"settings": superschema.SingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The settings for the VM.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"expose_hardware_virtualization": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether to expose hardware CPU virtualization to the guest OS.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
					"os_type": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The type of OS installed on the VM.",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(vm.GetAllOsTypes()...),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
					"storage_profile":  storageprofile.SuperSchema(),
					"guest_properties": vm.GuestPropertiesSuperSchema(),
					"affinity_rule_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the affinity rule to apply to this VM.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							// Computed : As it can get populated automatically by VDC default policy
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
					"customization": superschema.SingleNestedAttribute{
						Resource: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The customization settings for the VM.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Object{
								objectplanmodifier.UseStateForUnknown(),
							},
						},
						Attributes: map[string]superschema.Attribute{
							"force": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "`true` value will cause the VM to reboot on every `apply` operation.",
									Optional:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"enabled": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether guest customization is enabled or not.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"change_sid": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to change the SID of the VM. Applicable only for Windows VMs.",
									Optional:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"allow_local_admin_password": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to allow the local admin password to be changed.",
									Optional:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"must_change_password_on_first_login": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									Optional: true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"auto_generate_password": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to auto-generate the password.",
									Optional:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
									Validators: []validator.Bool{
										boolvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("admin_password"), path.MatchRelative().AtParent().AtName("auto_generate_password")),
									},
								},
							},
							"admin_password": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The admin password for the VM.",
									Optional:            true,
									Sensitive:           true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
									Validators: []validator.String{
										stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("admin_password"), path.MatchRelative().AtParent().AtName("auto_generate_password")),
									},
								},
							},
							"number_of_auto_logons": superschema.Int64Attribute{
								Resource: &schemaR.Int64Attribute{
									MarkdownDescription: "The number of times the VM should auto-login.",
									Optional:            true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
							"join_domain": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Enable this VM to join a domain.",
									Optional:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"join_org_domain": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Use organization's domain for joining.",
									Optional:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"join_domain_name": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The domain name to join.",
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"join_domain_user": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The domain user to join.",
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"join_domain_password": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The domain password to join.",
									Sensitive:           true,
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"join_domain_account_ou": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The domain account OU to join.",
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"init_script": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "The init script to run.",
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
							"hostname": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "Computer name to assign to this virtual machine. Default is the value of attribute `name`.",
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
									Validators: []validator.String{
										stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`), "Must be between 1 and 80 characters long and can contain only letters, numbers and hyphen. It must not contain only digits."),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// func vmSchema() schema.Schema {
// 	return schema.Schema{
// 		MarkdownDescription: "The virtual machine (vm) resource allows you to manage a virtual machine in the CloudAvenue.",
// 		Attributes: map[string]schema.Attribute{
// 			"id": schema.StringAttribute{
// 				MarkdownDescription: "The ID of the VM.",
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"vapp_name": schema.StringAttribute{
// 				MarkdownDescription: "The vApp this VM belongs to.",
// 				Required:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.RequiresReplace(),
// 				},
// 			},
// 			"name": schema.StringAttribute{
// 				MarkdownDescription: "The name of the VM. Unique within the vApp.",
// 				Required:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.RequiresReplace(),
// 				},
// 				Validators: []validator.String{
// 					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`), "must be between 1 and 80 characters long and can contain only letters, numbers and hyphen. It must not contain only digits."),
// 				},
// 			},
// 			"computer_name": schema.StringAttribute{
// 				MarkdownDescription: "Computer name to assign to this virtual machine. Default is `vm_name`.",
// 				Optional:            true,
// 				Computed:            true,
// 				Validators: []validator.String{
// 					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`), "must be between 1 and 80 characters long and can contain only letters, numbers and hyphen. It must not contain only digits."),
// 				},
// 			},
// 			"vdc": schema.StringAttribute{
// 				MarkdownDescription: "The name of the VDC this VM belongs to. If not specified, VDC define in provider will be used.",
// 				Optional:            true,
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.RequiresReplace(),
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"vapp_template_id": schema.StringAttribute{
// 				MarkdownDescription: "The URN of the vApp template to use for this VM. Required if `boot_image_id` is not set.",
// 				Optional:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.RequiresReplace(),
// 				},
// 			},
// 			"vm_name_in_template": schema.StringAttribute{
// 				MarkdownDescription: "The name of the VM in vApp Template to use. In cases when vApp template has more than one VM.",
// 				Optional:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.RequiresReplace(),
// 				},
// 			},
// 			"description": schema.StringAttribute{
// 				MarkdownDescription: "Description of the VM.",
// 				Optional:            true,
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"resource": schema.SingleNestedAttribute{
// 				MarkdownDescription: "Resource configuration for the VM.",
// 				Optional:            true,
// 				Computed:            true,
// 				Attributes:          vm.ResourceSchema(),
// 				PlanModifiers: []planmodifier.Object{
// 					objectplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"href": schema.StringAttribute{
// 				MarkdownDescription: "VM Hyper Reference",
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"accept_all_eulas": schema.BoolAttribute{
// 				MarkdownDescription: "Automatically accept EULA if OVA has it. Default is `false`.",
// 				Optional:            true,
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.Bool{
// 					fboolplanmodifier.SetDefault(false),
// 				},
// 			},
// 			"power_on": schema.BoolAttribute{
// 				MarkdownDescription: "A boolean value stating if this VM should be powered on. Default is `true`.",
// 				Optional:            true,
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.Bool{
// 					fboolplanmodifier.SetDefault(true),
// 				},
// 			},
// 			"storage_profile": schema.StringAttribute{
// 				MarkdownDescription: "Storage profile to override the default one.",
// 				Optional:            true,
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"os_type": schema.StringAttribute{
// 				MarkdownDescription: "Operating System type.",
// 				Optional:            true,
// 				Computed:            true,
// 				Validators: []validator.String{
// 					stringvalidator.OneOf(vm.GetAllOsTypes()...),
// 				},
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"boot_image_id": schema.StringAttribute{
// 				MarkdownDescription: "The ID of the boot image to use for this VM. Required if `vapp_template_id` is not set.",
// 				Optional:            true,
// 				Validators: []validator.String{
// 					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_template_id")),
// 				},
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.RequiresReplace(),
// 				},
// 			},
// 			"network_dhcp_wait_seconds": schema.Int64Attribute{
// 				MarkdownDescription: "Optional number of seconds to try and wait for DHCP IP (valid for `network` block only)",
// 				Optional:            true,
// 				Validators: []validator.Int64{
// 					int64validator.AtLeast(0),
// 				},
// 			},
// 			"networks": schema.ListNestedAttribute{
// 				MarkdownDescription: "A block to define network interface. Multiple can be used.",
// 				Optional:            true,
// 				Computed:            true,
// 				NestedObject: schema.NestedAttributeObject{
// 					Attributes: vm.NetworkSchema(),
// 				},
// 				PlanModifiers: []planmodifier.List{
// 					listplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"disks": schema.SetNestedAttribute{
// 				MarkdownDescription: "A list of disks attached to this VM.",
// 				Computed:            true,
// 				NestedObject: schema.NestedAttributeObject{
// 					Attributes: diskInternalExternalSchema(),
// 				},
// 			},
// 			"expose_hardware_virtualization": schema.BoolAttribute{
// 				MarkdownDescription: "Expose hardware-assisted CPU virtualization to guest OS. Default is `false`.",
// 				Optional:            true,
// 				Computed:            true,
// 				PlanModifiers: []planmodifier.Bool{
// 					fboolplanmodifier.SetDefault(false),
// 					boolplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"guest_properties": vm.GuestPropertiesSchema(),
// 			"customization": schema.SingleNestedAttribute{
// 				MarkdownDescription: "Guest customization block. ",
// 				Optional:            true,
// 				Computed:            true,
// 				Attributes:          vm.CustomizationSchema(),
// 				PlanModifiers: []planmodifier.Object{
// 					objectplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"prevent_update_power_off": schema.BoolAttribute{
// 				MarkdownDescription: "`true` if the update of resource should fail when virtual machine power off needed. Default is `false`.",
// 				Optional:            true,
// 				PlanModifiers: []planmodifier.Bool{
// 					fboolplanmodifier.SetDefault(false),
// 				},
// 			},
// 			"sizing_policy_id": schema.StringAttribute{
// 				MarkdownDescription: "VM placement policy ID. Has to be assigned to Org VDC.",
// 				Optional:            true,
// 				Computed:            true, // As it can get populated automatically by VDC default policy
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"placement_policy_id": schema.StringAttribute{
// 				MarkdownDescription: "VM placement policy ID. Has to be assigned to Org VDC.",
// 				Optional:            true,
// 				Computed:            true, // As it can get populated automatically by VDC default policy
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"status_code": schema.Int64Attribute{
// 				MarkdownDescription: "Shows the status code of the VM",
// 				Computed:            true,
// 			},
// 			"status_text": schema.StringAttribute{
// 				MarkdownDescription: "Shows the status text of the VM",
// 				Computed:            true,
// 			},
// 		},
// 	}
// }

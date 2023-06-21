package vm

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	const (
		coldUpdate = `<a href="#restartrequired" style="color:red">(Restart Required)</a>`
	)

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
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
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
					MarkdownDescription: "The description of the VM",
				},
				Resource: &schemaR.StringAttribute{
					Optional:            true,
					MarkdownDescription: coldUpdate,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"deploy_os": superschema.SingleNestedAttribute{
				Resource: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Settings for deploying the operating system on the VM.",
					Optional:            true,
				},
				Attributes: map[string]superschema.Attribute{
					"vapp_template_id": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the vApp template to use for the VM.",
							Optional:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("boot_image_id")),
							},
						},
					},
					"vm_name_in_template": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the VM in the vApp template.",
							Optional:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("boot_image_id")),
							},
						},
					},
					"boot_image_id": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the boot image to use for the VM.",
							Optional:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("vapp_template_id")),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"accept_all_eulas": superschema.BoolAttribute{
						Resource: &schemaR.BoolAttribute{
							MarkdownDescription: "Automatically accept EULA if OVA has it",
							Optional:            true,
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
							Default:  booldefault.StaticBool(true),
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
				},
				Attributes: map[string]superschema.Attribute{
					"cpus": superschema.Int64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The number of virtual CPUs to allocate to the VM.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							MarkdownDescription: coldUpdate,
							Optional:            true,
							Default:             int64default.StaticInt64(1),
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
							MarkdownDescription: coldUpdate,
							Optional:            true,
							Default:             int64default.StaticInt64(1),
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
							Validators: []validator.Int64{
								fint64validator.AttributeIsDivisibleByAnInteger(path.MatchRelative().AtParent().AtName("cpus")),
							},
						},
					},
					"cpu_hot_add_enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether CPU hot add is enabled or not",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							MarkdownDescription: coldUpdate,
							Optional:            true,
							Default:             booldefault.StaticBool(true),
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
							MarkdownDescription: coldUpdate,
							Optional:            true,
							Default:             int64default.StaticInt64(1024),
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
							MarkdownDescription: coldUpdate,
							Optional:            true,
							Default:             booldefault.StaticBool(true),
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
							MarkdownDescription: coldUpdate,
							Optional:            true,
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
										fstringvalidator.OneOfWithDescription(
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "vapp",
												Description: "A vApp network. This network is only available in your vApp structure.",
											},
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "org",
												Description: "An organization network. This network can be a network isolated or routed in your Organization.",
											},
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "none",
												Description: "No network.",
											},
										),
										// stringvalidator.OneOf("vapp", "org", "none"),
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
										fstringvalidator.OneOfWithDescription(
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "DHCP",
												Description: "IP address is obtained from a DHCP service.",
											},
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "POOL",
												Description: "Static IP address is allocated automatically from defined static pool in network.",
											},
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "MANUAL",
												Description: "SIP address is assigned manually in the ip field. Must be valid IP address from static pool.",
											},
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "NONE",
												Description: "No IP address will be set because VM will have a NIC without network.",
											},
										),
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
									Optional: true,
									Validators: []validator.String{
										fstringvalidator.RequireIfAttributeIsOneOf(
											path.MatchRelative().AtParent().AtName("type"),
											[]attr.Value{
												types.StringValue("vapp"),
												types.StringValue("org"),
											},
										),
										fstringvalidator.NullIfAttributeIsOneOf(
											path.MatchRelative().AtParent().AtName("type"),
											[]attr.Value{
												types.StringValue("none"),
											},
										),
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
									Optional: true,
									Validators: []validator.String{
										fstringvalidator.IsIP(),
										fstringvalidator.RequireIfAttributeIsOneOf(
											path.MatchRelative().AtParent().AtName("ip_allocation_mode"),
											[]attr.Value{
												types.StringValue("MANUAL"),
												types.StringValue("POOL"),
											},
										),
										fstringvalidator.NullIfAttributeIsOneOf(
											path.MatchRelative().AtParent().AtName("ip_allocation_mode"),
											[]attr.Value{
												types.StringValue("DHCP"),
												types.StringValue("NONE"),
											},
										),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"is_primary": superschema.BoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether this network is the primary network for the VM",
									Computed:            true,
								},
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: coldUpdate,
									Optional:            true,
									Default:             booldefault.StaticBool(false),
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
										fstringvalidator.OneOfWithDescription(
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "VMXNET3",
												Description: "A VMXNET3 adapter it's an paravirtualized NIC designed for performance. VMXNET 3 offers all the features available in VMXNET 2 and adds several new features.",
											},
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "E1000E",
												Description: "A E1000E adapter it's an Emulated version of the Intel 82574 Gigabit Ethernet NIC. E1000E is the default adapter for Windows 8 and Windows Server 2012.",
											},
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "VMXNET3VRDMA",
												Description: "VMXNET3VRDMA adapter it's a paravirtualized NIC that supports remote direct memory access (RDMA) between virtual machines through the OFED verbs API.",
											},
											fstringvalidator.OneOfWithDescriptionValues{
												Value:       "SRIOVETHERNETCARD",
												Description: "SRIOVETHERNETCARD adapter it's a representation of a virtual function (VF) on a physical NIC with SR-IOV support.",
											},
										),
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
							MarkdownDescription: "Whether to expose hardware CPU virtualization to the guest OS",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							MarkdownDescription: coldUpdate,
							Optional:            true,
							Default:             booldefault.StaticBool(false),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
					"os_type": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The type of OS installed on the VM " + coldUpdate,
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								fstringvalidator.OneOfWithDescription(vm.GetAllOsTypesWithDescription()...),
								// TODO Validator field is require if attribute deploy_os.boot_image_id is set
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
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
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"force": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "`true` value will cause the VM to reboot on every `apply` operation.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.BoolAttribute{
									Computed: true,
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
								DataSource: &schemaD.BoolAttribute{
									Computed: true,
								},
							},
							"change_sid": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to change the SID of the VM. Applicable only for Windows VMs.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.BoolAttribute{
									Computed: true,
								},
							},
							"allow_local_admin_password": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to allow the local admin password to be changed.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.BoolAttribute{
									Computed: true,
								},
							},
							"must_change_password_on_first_login": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether the password must be changed on first login.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.BoolAttribute{
									Computed: true,
								},
							},
							"auto_generate_password": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to auto-generate the password.",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
									Validators: []validator.Bool{
										boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("admin_password")),
									},
								},
								DataSource: &schemaD.BoolAttribute{
									Computed: true,
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
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("auto_generate_password")),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"number_of_auto_logons": superschema.Int64Attribute{
								Resource: &schemaR.Int64Attribute{
									MarkdownDescription: "The number of times the VM should auto-login.",
									Optional:            true,
									Computed:            true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.Int64Attribute{
									Computed: true,
								},
							},
							"join_domain": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Enable this VM to join a domain.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.BoolAttribute{
									Computed: true,
								},
							},
							"join_org_domain": superschema.BoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "Use organization's domain for joining.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.BoolAttribute{
									Computed: true,
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
								DataSource: &schemaD.StringAttribute{
									Computed: true,
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
								DataSource: &schemaD.StringAttribute{
									Computed: true,
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
								DataSource: &schemaD.StringAttribute{
									Computed: true,
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
								DataSource: &schemaD.StringAttribute{
									Computed: true,
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
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"hostname": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									MarkdownDescription: "Computer name to assign to this virtual machine. Default is the value of attribute `name`.",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
									Validators: []validator.String{
										stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`), "Must be between 1 and 80 characters long and can contain only letters, numbers and hyphen. It must not contain only digits."),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

package vm

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/storageprofile"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fint64validator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/int64validator"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
)

// TODO : Remove unused imports.
// ! This is outside import block because golangci-lint remove commented import.
// * Hashicorp Validators
// "github.com/Hashicorp/terraform-plugin-framework-validators/stringvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/boolvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/int64validator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/float64validator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/listvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/mapvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/setvalidator"

// * Hashicorp Plan Modifiers Resource
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"

// * Hashicorp Plan Modifiers DataSource
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/stringplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/boolplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/int64planmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/listplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/mapplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/setplanmodifier"

// * Hashicorp Default Values
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

// * FrangipaneTeam Custom Validators
// fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
// fboolvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/boolvalidator"
// fint64validator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/int64validator"
// flistvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/listvalidator"
// fmapvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/mapvalidator"
// fsetvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/setvalidator"

// * FrangipaneTeam Custom Plan Modifiers
// fstringplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/stringplanmodifier"
// fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"
// fint64planmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/int64planmodifier"
// flistplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/listplanmodifier"
// fmapplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/mapplanmodifier"
// fsetplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/setplanmodifier"

// How to use types generator:
// 1. Define the schema in the file internal/provider/vm/v_2_schema.go
// 2. Add the resource or data source to the file internal/provider/provider_resources.go or internal/provider/provider_data_sources.go respectively
// 3. Launch the following command to generate golang structs for the schema:
// go run ./cmd/types-generator/*.go -file internal/provider/vm/v_2_schema.go -resource cloudavenue_vm_v_2 -is-resource
func vmV2Schema(_ context.Context) superschema.Schema {
	const (
		coldUpdate = `<a href="#restartrequired" style="color:red">(Restart Required)</a>`
	)

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vm_v2` resource allows you to manage ...",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vm_v2` data source allows you to retrieve information about an ...",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the VM.",
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the VM. Unique within the vApp.",
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`), "Must be between 1 and 80 characters long and can contain only letters, numbers and hyphen. It must not contain only digits."),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"vapp_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The vApp the VM belongs to.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"vapp_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The vApp the VM belongs to.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"vdc": vdc.SuperSchemaSuperType(),
			"description": superschema.SuperStringAttribute{
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
			"deploy_os": superschema.SuperSingleNestedAttributeOf[VMV2ModelDeployOS]{
				Resource: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Settings for deploying the operating system on the VM.",
					Optional:            true,
				},
				Attributes: superschema.Attributes{
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
					},
					"accept_all_eulas": superschema.BoolAttribute{
						Resource: &schemaR.BoolAttribute{
							MarkdownDescription: "Automatically accept EULA if OVA has it",
							Optional:            true,
						},
					},
				},
			},
			"state": superschema.SuperSingleNestedAttributeOf[VMV2ModelState]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The state of the VM.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: superschema.Attributes{
					"power_on": superschema.SuperBoolAttribute{
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
					"status": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The power status of the VM.",
							Computed:            true,
						},
					},
				},
			},
			"resource": superschema.SuperSingleNestedAttributeOf[VMV2ModelResource]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The resource of the VM.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: superschema.Attributes{
					"cpus": superschema.SuperInt64Attribute{
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
					"cpus_cores": superschema.SuperInt64Attribute{
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
					"cpu_hot_add_enabled": superschema.SuperBoolAttribute{
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
					"memory": superschema.SuperInt64Attribute{
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
					"memory_hot_add_enabled": superschema.SuperBoolAttribute{
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
					"networks": superschema.SuperListNestedAttributeOf[VMV2ModelResourceNetworks]{
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
						Attributes: superschema.Attributes{
							"type": superschema.SuperStringAttribute{
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
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"ip_allocation_mode": superschema.SuperStringAttribute{
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
							"name": superschema.SuperStringAttribute{
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
							"ip": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The IP address to assign to this VM on this network.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										fstringvalidator.IsIP(),
										fstringvalidator.RequireIfAttributeIsOneOf(
											path.MatchRelative().AtParent().AtName("ip_allocation_mode"),
											[]attr.Value{
												types.StringValue("MANUAL"),
											},
										),
										fstringvalidator.NullIfAttributeIsOneOf(
											path.MatchRelative().AtParent().AtName("ip_allocation_mode"),
											[]attr.Value{
												types.StringValue("DHCP"),
												types.StringValue("NONE"),
												types.StringValue("POOL"),
											},
										),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"is_primary": superschema.SuperBoolAttribute{
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
							"mac": superschema.SuperStringAttribute{
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
							"adapter_type": superschema.SuperStringAttribute{
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
							"connected": superschema.SuperBoolAttribute{
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
			"settings": superschema.SuperSingleNestedAttributeOf[VMV2ModelSettings]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The settings for the VM.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: superschema.Attributes{
					"expose_hardware_virtualization": superschema.SuperBoolAttribute{
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
					"os_type": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Operating System type",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "  to be installed on the VM." + coldUpdate,
							Optional:            true,
							Validators: []validator.String{
								fstringvalidator.OneOfWithDescription(vm.GetAllOsTypesWithDescription()...),
								// TODO Validator field is require if attribute deploy_os.boot_image_id is set
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: " installed on the VM.",
						},
					},
					// TODO
					"storage_profile": storageprofile.SuperSchema(),
					// TODO
					"guest_properties": vm.GuestPropertiesSuperSchema(),
					"affinity_rule_id": superschema.SuperStringAttribute{
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
					"customization": superschema.SuperSingleNestedAttributeOf[any]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The customization settings for the VM. To enable the customization, set the `enabled` attribute to `true`.",
							Computed:            true,
						},
						Resource: &schemaR.SingleNestedAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.Object{
								objectplanmodifier.UseStateForUnknown(),
							},
						},
						Attributes: superschema.Attributes{
							"force": superschema.SuperBoolAttribute{
								Resource: &schemaR.BoolAttribute{
									MarkdownDescription: "`true` value will cause the VM to reboot on every terraform `apply` operation.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								DataSource: &schemaD.BoolAttribute{
									MarkdownDescription: "This attributes is not set in the data source.",
									Computed:            true,
								},
							},
							"enabled": superschema.SuperBoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether guest customization is enabled or not. If enabled, the VM will be customized on first boot.",
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
							"change_sid": superschema.SuperBoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to change the SID of the VM. Applicable only for Windows VMs.",
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
							"allow_local_admin_password": superschema.SuperBoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to allow the local admin password to be changed.",
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
							"must_change_password_on_first_login": superschema.SuperBoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether the password must be changed on first login.",
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
							"auto_generate_password": superschema.SuperBoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Whether to auto-generate the password.",
									Computed:            true,
								},
								Resource: &schemaR.BoolAttribute{
									Optional: true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
									Validators: []validator.Bool{
										boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("admin_password")),
									},
								},
							},
							"admin_password": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The admin password for the VM.",
									Sensitive:           true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
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
							"number_of_auto_logons": superschema.SuperInt64Attribute{
								Common: &schemaR.Int64Attribute{
									MarkdownDescription: "The number of times the VM should auto-login.",
									Computed:            true,
								},
								Resource: &schemaR.Int64Attribute{
									Optional: true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
							"join_domain": superschema.SuperBoolAttribute{
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
									MarkdownDescription: "If true, this VM is join to a domain.",
									Computed:            true,
								},
							},
							"join_org_domain": superschema.SuperBoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Use organization's domain for joining.",
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
							"join_domain_name": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The domain name to join.",
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
							"join_domain_user": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The domain user to join.",
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
							"join_domain_password": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The domain password to join.",
									Sensitive:           true,
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
							"join_domain_account_ou": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The domain account OU to join.",
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
							"init_script": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The init script to run.",
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
							"hostname": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Computer name to assign to this virtual machine. Default is the value of attribute `name`.",
									Computed:            true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
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

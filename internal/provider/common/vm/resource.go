package vm

import (
	"context"
	"fmt"

	fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type Resource struct {
	// CPU
	CPUs             types.Int64 `tfsdk:"cpus"`
	CPUCores         types.Int64 `tfsdk:"cpu_cores"`
	CPUHotAddEnabled types.Bool  `tfsdk:"cpu_hot_add_enabled"`

	// Memory
	Memory              types.Int64 `tfsdk:"memory"`
	MemoryHotAddEnabled types.Bool  `tfsdk:"memory_hot_add_enabled"`
}

func ResourceAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"cpus":                types.Int64Type,
		"cpu_cores":           types.Int64Type,
		"cpu_hot_add_enabled": types.BoolType,

		"memory":                 types.Int64Type,
		"memory_hot_add_enabled": types.BoolType,
	}
}

// ToAttrValue converts the Customization struct to a map of attr.Value.
func (r *Resource) ToAttrValue() map[string]attr.Value {
	return map[string]attr.Value{
		"cpus":                r.CPUs,
		"cpu_cores":           r.CPUCores,
		"cpu_hot_add_enabled": r.CPUHotAddEnabled,

		"memory":                 r.Memory,
		"memory_hot_add_enabled": r.MemoryHotAddEnabled,
	}
}

// ObjectType returns the type of the resource object.
func (r *Resource) ObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: ResourceAttrType(),
	}
}

// ToPlan converts the resource struct to a plan.
func (r *Resource) ToPlan() basetypes.ObjectValue {
	if r == nil {
		return types.ObjectNull(ResourceAttrType())
	}

	return types.ObjectValueMust(ResourceAttrType(), r.ToAttrValue())
}

// ResourceFromPlan converts a plan to a resource struct.
func ResourceFromPlan(ctx context.Context, x types.Object) (*Resource, diag.Diagnostics) {
	if x.IsNull() || x.IsUnknown() {
		return &Resource{}, diag.Diagnostics{}
	}

	r := &Resource{}

	d := x.As(ctx, r, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})

	return r, d
}

// ResourceSchema returns the schema for the resource.
func ResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cpus": schema.Int64Attribute{
			MarkdownDescription: "The number of virtual CPUs to allocate to the VM.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"cpu_cores": schema.Int64Attribute{
			MarkdownDescription: "The number of cores per socket.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"cpu_hot_add_enabled": schema.BoolAttribute{
			MarkdownDescription: "`true` if the virtual machine supports addition of virtual CPUs while powered on. Default is `false`.",
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				fboolplanmodifier.SetDefault(false),
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"memory": schema.Int64Attribute{
			MarkdownDescription: "The amount of memory (in MB) to allocate to the VM.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			// TODO : Add validator to check if value is a multiple of 4
		},
		"memory_hot_add_enabled": schema.BoolAttribute{
			MarkdownDescription: "`true` if the virtual machine supports addition of memory resources while powered on. Default is `false`.",
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				fboolplanmodifier.SetDefault(false),
				boolplanmodifier.UseStateForUnknown(),
			},
		},
	}
}

// ResourceRead is the read function for the resource.
func ResourceRead(vm *govcd.VM) (Resource, error) {
	if vm == nil {
		return Resource{}, fmt.Errorf("vm is nil")
	}

	var resource Resource

	if vm.VM.VmSpecSection != nil {
		// CPU
		resource.CPUs = types.Int64Value(int64(*vm.VM.VmSpecSection.NumCpus))
		resource.CPUCores = types.Int64Value(int64(*vm.VM.VmSpecSection.NumCoresPerSocket))

		// Memory
		if vm.VM.VmSpecSection.MemoryResourceMb != nil {
			resource.Memory = types.Int64Value((vm.VM.VmSpecSection.MemoryResourceMb.Configured))
		}
	}

	if vm.VM.VMCapabilities != nil {
		// HotAddEnabled
		resource.CPUHotAddEnabled = types.BoolValue(vm.VM.VMCapabilities.CPUHotAddEnabled)
		resource.MemoryHotAddEnabled = types.BoolValue(vm.VM.VMCapabilities.MemoryHotAddEnabled)
	}

	return resource, nil
}

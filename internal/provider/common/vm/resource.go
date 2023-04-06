package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VMResourceModelResource struct { //nolint:revive
	CPUs                types.Int64 `tfsdk:"cpus"`
	CPUsCores           types.Int64 `tfsdk:"cpus_cores"`
	CPUHotAddEnabled    types.Bool  `tfsdk:"cpu_hot_add_enabled"`
	Memory              types.Int64 `tfsdk:"memory"`
	MemoryHotAddEnabled types.Bool  `tfsdk:"memory_hot_add_enabled"`
	Networks            types.List  `tfsdk:"networks"`
}

// attrTypes() returns the types of the attributes of the Resource attribute.
func (r *VMResourceModelResource) attrTypes(networks *VMResourceModelResourceNetworks) map[string]attr.Type {
	return map[string]attr.Type{
		"cpus":                   types.Int64Type,
		"cpus_cores":             types.Int64Type,
		"memory":                 types.Int64Type,
		"cpu_hot_add_enabled":    types.BoolType,
		"memory_hot_add_enabled": types.BoolType,
		"networks":               types.ListType{ElemType: types.ObjectType{networks.AttrTypes()}},
	}
}

// toAttrValues() returns the values of the attributes of the Resource attribute.
func (r *VMResourceModelResource) toAttrValues(ctx context.Context, networks *VMResourceModelResourceNetworks) map[string]attr.Value {
	net, err := networks.ToPlan(ctx)
	if err != nil {
		return nil
	}
	return map[string]attr.Value{
		"cpus":                   r.CPUs,
		"cpus_cores":             r.CPUsCores,
		"memory":                 r.Memory,
		"cpu_hot_add_enabled":    r.CPUHotAddEnabled,
		"memory_hot_add_enabled": r.MemoryHotAddEnabled,
		"networks":               net,
	}
}

// Equal returns true if the values of the attributes of the Resource attribute are equal.
func (r *VMResourceModelResource) Equal(other *VMResourceModelResource) bool {
	return r.CPUs.Equal(other.CPUs) &&
		r.CPUsCores.Equal(other.CPUsCores) &&
		r.Memory.Equal(other.Memory) &&
		r.CPUHotAddEnabled.Equal(other.CPUHotAddEnabled) &&
		r.MemoryHotAddEnabled.Equal(other.MemoryHotAddEnabled) &&
		r.Networks.Equal(other.Networks)
}

// ToPlan returns the value of the Resource attribute, if set, as a types.Object.
func (r *VMResourceModelResource) ToPlan(ctx context.Context, networks *VMResourceModelResourceNetworks) types.Object {
	if r == nil {
		return types.Object{}
	}

	return types.ObjectValueMust(r.attrTypes(networks), r.toAttrValues(ctx, networks))
}

// ResourceRead is the read function for the resource.
func (v VM) ResourceRead(_ context.Context) (resource *VMResourceModelResource) {
	resource = &VMResourceModelResource{
		CPUs:                types.Int64Null(),
		CPUsCores:           types.Int64Null(),
		Memory:              types.Int64Null(),
		CPUHotAddEnabled:    types.BoolNull(),
		MemoryHotAddEnabled: types.BoolNull(),
	}

	if v.CpusIsDefined() {
		resource.CPUs = types.Int64Value(int64(v.GetCpus()))
	}

	if v.CpusCoresIsDefined() {
		resource.CPUsCores = types.Int64Value(int64(v.GetCpusCores()))
	}

	if v.MemoryIsDefined() {
		resource.Memory = types.Int64Value(v.GetMemory())
	}

	if v.HotAddIsDefined() {
		resource.CPUHotAddEnabled = types.BoolValue(v.GetCpuHotAddEnabled())
		resource.MemoryHotAddEnabled = types.BoolValue(v.GetMemoryHotAddEnabled())
	}

	return
}

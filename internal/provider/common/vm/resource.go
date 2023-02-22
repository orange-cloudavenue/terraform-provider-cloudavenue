package vm

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/boolpm"
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

// ResourceSchema returns the schema for the resource.
func ResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cpus": schema.Int64Attribute{
			MarkdownDescription: "The number of virtual CPUs to allocate to the VM.",
			Optional:            true,
			Computed:            true,
		},
		"cpu_cores": schema.Int64Attribute{
			MarkdownDescription: "The number of cores per socket.",
			Optional:            true,
			Computed:            true,
		},
		"cpu_hot_add_enabled": schema.BoolAttribute{
			MarkdownDescription: "`true` if the virtual machine supports addition of virtual CPUs while powered on. Default is `false`.",
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				boolpm.SetDefault(false),
			},
		},
		"memory": schema.Int64Attribute{
			MarkdownDescription: "The amount of memory (in MB) to allocate to the VM.",
			Optional:            true,
			Computed:            true,
			// TODO : Add validator to check if value is a multiple of 4
		},
		"memory_hot_add_enabled": schema.BoolAttribute{
			MarkdownDescription: "`true` if the virtual machine supports addition of memory resources while powered on. Default is `false`.",
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				boolpm.SetDefault(false),
			},
		},
	}
}

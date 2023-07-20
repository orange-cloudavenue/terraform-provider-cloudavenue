package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type storageProfiles []profileDataSourceModel

func (s *storageProfiles) attrTypes(_ context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"name":                  types.StringType,
		"vdc":                   types.StringType,
		"limit":                 types.Int64Type,
		"used_storage":          types.Int64Type,
		"default":               types.BoolType,
		"enabled":               types.BoolType,
		"iops_allocated":        types.Int64Type,
		"units":                 types.StringType,
		"iops_limiting_enabled": types.BoolType,
		"maximum_disk_iops":     types.Int64Type,
		"default_disk_iops":     types.Int64Type,
		"disk_iops_per_gb_max":  types.Int64Type,
		"iops_limit":            types.Int64Type,
	}
}

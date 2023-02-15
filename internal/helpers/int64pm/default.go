package int64pm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// SetDefault returns a plan modifier that conditionally requires
// resource replacement if:
//
//   - The resource is planned for update.
//   - The plan and state values are not equal.
//   - The plan or state values are not null or known
func SetDefault(i int64) planmodifier.Int64 {
	return setDefaultFunc(
		func(ctx context.Context, req planmodifier.Int64Request, resp *DefaultFuncResponse) {
			resp.Value = i
		},
		"Set default value",
		"Set default value",
	)
}

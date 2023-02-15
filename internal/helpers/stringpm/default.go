// Package stringpm provides a plan modifier for string values.
package stringpm

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
func SetDefault(str string) planmodifier.String {
	return setDefaultFunc(
		func(_ context.Context, _ planmodifier.StringRequest, resp *DefaultFuncResponse) {
			resp.Value = str
		},
		"Set default value",
		"Set default value",
	)
}

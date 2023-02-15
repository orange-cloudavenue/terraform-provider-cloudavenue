// Package int64pm provides a plan modifier for int64 values.
package int64pm

import "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

// SetDefaultFunc returns a plan modifier that conditionally requires
// resource replacement if:
//
//   - The resource is planned for update.
//   - The plan and state values are not equal.
//   - The plan or state values are not null or known
func SetDefaultFunc(f DefaultFunc) planmodifier.Int64 {
	return setDefaultFunc(
		f,
		"Set default value",
		"Set default value",
	)
}

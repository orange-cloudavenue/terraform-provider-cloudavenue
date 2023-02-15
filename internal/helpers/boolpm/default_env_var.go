package boolpm

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// SetDefaultFunc returns a plan modifier that conditionally requires
// resource replacement if:
//
//   - The resource is planned for update.
//   - The plan and state values are not equal.
//   - The plan or state values are not null or known
func SetDefaultEnvVar(envVar string) planmodifier.Bool {
	return setDefaultFunc(
		func(_ context.Context, _ planmodifier.BoolRequest, resp *DefaultFuncResponse) {
			v := os.Getenv(envVar)
			if v != "" {
				// Parse string to bool
				b, err := strconv.ParseBool(v)
				if err != nil {
					resp.Diagnostics.AddError("Environment variable set but is not Boolean", fmt.Sprintf("The environment variable %s is set but is not a Boolean", envVar))
					return
				}
				resp.Value = b
			} else {
				resp.Diagnostics.AddError("Environment variable not set", fmt.Sprintf("The environment variable %s is not set", envVar))
			}
		},
		"Set default value from environment variable",
		"Set default value from environment variable",
	)
}

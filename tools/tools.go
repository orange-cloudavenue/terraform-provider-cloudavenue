//go:build tools

package tools

import (
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"

	_ "github.com/orange-cloudavenue/terraform-provider-cloudavenue/cmd/vdc-doc"
)

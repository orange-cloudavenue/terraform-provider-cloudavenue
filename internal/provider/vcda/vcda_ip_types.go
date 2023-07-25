package vcda

import "github.com/hashicorp/terraform-plugin-framework/types"

type vcdaIPResourceModel struct {
	ID        types.String `tfsdk:"id"`
	IPAddress types.String `tfsdk:"ip_address"`
}

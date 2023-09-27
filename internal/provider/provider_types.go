package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type cloudavenueProviderModel struct {
	URL               types.String `tfsdk:"url"`
	User              types.String `tfsdk:"user"`
	Password          types.String `tfsdk:"password"`
	Org               types.String `tfsdk:"org"`
	VDC               types.String `tfsdk:"vdc"`
	NetBackupURL      types.String `tfsdk:"netbackup_url"`
	NetBackupUser     types.String `tfsdk:"netbackup_user"`
	NetBackupPassword types.String `tfsdk:"netbackup_password"`
}

package netbackup

import "github.com/hashicorp/terraform-plugin-framework/types"

type BackupModel struct {
	ID types.String `tfsdk:"id"`
}

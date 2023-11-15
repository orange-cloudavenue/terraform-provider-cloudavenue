package s3

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type UserDataSourceModel struct {
	ID          supertypes.StringValue `tfsdk:"id"`
	Username    supertypes.StringValue `tfsdk:"user_name"`
	UserID      supertypes.StringValue `tfsdk:"user_id"`
	FullName    supertypes.StringValue `tfsdk:"full_name"`
	CanonicalID supertypes.StringValue `tfsdk:"canonical_id"`
}

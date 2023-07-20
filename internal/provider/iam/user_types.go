package iam

import "github.com/hashicorp/terraform-plugin-framework/types"

type userResourceModel struct {
	// Base
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoleName        types.String `tfsdk:"role_name"`
	FullName        types.String `tfsdk:"full_name"`
	Email           types.String `tfsdk:"email"`
	Telephone       types.String `tfsdk:"telephone"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	DeployedVMQuota types.Int64  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   types.Int64  `tfsdk:"stored_vm_quota"`

	// Specific
	Password      types.String `tfsdk:"password"`
	TakeOwnership types.Bool   `tfsdk:"take_ownership"`
}

type userDataSourceModel struct {
	// Base
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoleName        types.String `tfsdk:"role_name"`
	FullName        types.String `tfsdk:"full_name"`
	Email           types.String `tfsdk:"email"`
	Telephone       types.String `tfsdk:"telephone"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	DeployedVMQuota types.Int64  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   types.Int64  `tfsdk:"stored_vm_quota"`

	// Specific
	ProviderType types.String `tfsdk:"provider_type"`
}

// Package vapp provides a Terraform resource to manage vApps.
package vapp

import "github.com/hashicorp/terraform-plugin-framework/types"

type vappLeaseModel struct {
	RuntimeLeaseInSec types.Int64 `tfsdk:"runtime_lease_in_sec"`
	StorageLeaseInSec types.Int64 `tfsdk:"storage_lease_in_sec"`
}

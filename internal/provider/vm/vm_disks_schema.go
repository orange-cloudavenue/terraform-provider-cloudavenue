package vm

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func disksSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_vm_disks` data source allows you to retrieve information about an existing disks in vApp and VM.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Generated ID of the resource.",
				},
			},
			"vdc": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The name of vDC to use, optional if defined at provider level.",
					Optional:            true,
					Computed:            true,
				},
			},
			"vapp_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
						fstringvalidator.IsURN(),
						fstringvalidator.PrefixContains(urn.VAPP.String()),
					},
				},
			},
			"vapp_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
					},
				},
			},
			"vm_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vm_name"), path.MatchRoot("vm_id")),
						fstringvalidator.IsURN(),
						fstringvalidator.PrefixContains(urn.VM.String()),
					},
				},
			},
			"vm_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("vm_name"), path.MatchRoot("vm_id")),
					},
				},
			},
			"disks": superschema.SuperListNestedAttribute{
				DataSource: &schemaD.ListNestedAttribute{
					Computed:            true,
					MarkdownDescription: "List of disks in the vApp and attached to the VM.",
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the disk.",
						},
					},
					"name": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the disk.",
						},
					},
					"is_detachable": superschema.SuperBoolAttribute{
						DataSource: &schemaD.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "If set to true, the disk could be detached from the VM. If set to false, the disk canot detached to the VM.",
						},
					},
					"storage_profile": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the storage profile.",
						},
					},
					"size_in_mb": superschema.SuperInt64Attribute{
						DataSource: &schemaD.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The size of the disk in MB.",
						},
					},
				},
			},
		},
	}
}

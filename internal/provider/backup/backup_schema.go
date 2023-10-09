package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
)

// BackupSchema returns the schema for the backup resource.
func backupSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_backup` resource allows you to manage backup strategy for `vdc`, `vapp` and `vm` from NetBackup solution. [Please refer to the documentation for more information.](https://wiki.cloudavenue.orange-business.com/wiki/Backup)",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_backup` data source allows you to retrieve information about a backup of NetBackup solution.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					Computed:            true,
					MarkdownDescription: "The ID of the backup.",
				},
				Resource: &schemaR.Int64Attribute{
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.Int64Attribute{
					Optional: true,
				},
			},
			"type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Scope of the backup.",
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					Validators: []validator.String{
						stringvalidator.OneOf("vdc", "vapp", "vm"),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"target_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the target. A target can be a VDC, a VApp or a VM.",
					Optional:            true,
				},
				Resource: &schemaR.StringAttribute{
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("target_id"), path.MatchRoot("target_name")),
						fstringvalidator.IsUUID(),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"target_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the target. A target can be a VDC, a VApp or a VM.",
					Optional:            true,
				},
				Resource: &schemaR.StringAttribute{
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("target_id"), path.MatchRoot("target_name")),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"policies": superschema.SuperSetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "The backup policies of the target.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Required: true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"policy_id": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The ID of the backup policy.",
							Computed:            true,
						},
					},
					"policy_name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the backup policy. Each letter represent a strategy predefined: D = Daily, W = Weekly, M = Monthly, X = Replication, The number is the retention period. [Please refer to the documentation for more information.](https://wiki.cloudavenue.orange-business.com/wiki/Backup)",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("D6", "D30", "D30NQ", "D60", "W4", "M3", "M12", "XD6", "XD30", "XD60", "XW4", "XM3", "XM12"),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

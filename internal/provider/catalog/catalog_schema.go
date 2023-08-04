package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

/*
catalogSchema

This function is used to create the schema for the catalog resource and datasource.
*/
func catalogSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Catalog allows you to",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "manage a catalog in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "retrieve information about a catalog in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the catalog.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the catalog.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"created_at": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The creation date of the catalog.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"description": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the catalog.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"owner_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The owner name of the catalog.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"preserve_identity_information": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Include BIOS UUIDs and MAC addresses in the downloaded OVF package. Keep in mind that preserving this identity information reduces the package's portability, so only include it when necessary.",
					Computed:            true,
				},
			},
			"number_of_media": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "The number of media in the catalog.",
					Computed:            true,
				},
			},
			"media_item_list": superschema.ListAttribute{
				DataSource: &schemaD.ListAttribute{
					MarkdownDescription: "The list of media items in the catalog.",
					Computed:            true,
					ElementType:         types.StringType,
				},
			},
			"is_shared": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether the catalog is shared.",
					Computed:            true,
				},
			},
			"is_local": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether the catalog is local.",
					Computed:            true,
				},
			},
			"is_published": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether the catalog is published.",
					Computed:            true,
				},
			},
			"is_cached": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether the catalog is cached.",
					Computed:            true,
				},
			},
			"storage_profile": superschema.StringAttribute{
				// TODO - this is a reference to a storage profile, not a string
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "Storage profile to override the VM default one.",
					Optional:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"delete_force": superschema.BoolAttribute{
				Resource: &schemaR.BoolAttribute{
					Required:            true,
					MarkdownDescription: "When destroying a catalog, use `delete_force=True` along with `delete_recursive=True` to remove the catalog and any contained objects, regardless of their state.",
				},
			},
			"delete_recursive": superschema.BoolAttribute{
				Resource: &schemaR.BoolAttribute{
					Required:            true,
					MarkdownDescription: "When destroying a catalog, use `delete_recursive=True to remove the catalog and any contained objects that are in a state permitting removal.",
				},
			},
		},
	}
}

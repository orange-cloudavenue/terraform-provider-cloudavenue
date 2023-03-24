package catalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
)

type catalogDataSourceModel struct {
	// BASE
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	OwnerName   types.String `tfsdk:"owner_name"`

	// SPECIFIC DATA SOURCE
	PreserveIdentityInformation types.Bool  `tfsdk:"preserve_identity_information"`
	NumberOfMedia               types.Int64 `tfsdk:"number_of_media"`
	MediaItemList               types.List  `tfsdk:"media_item_list"`
	IsShared                    types.Bool  `tfsdk:"is_shared"`
	IsPublished                 types.Bool  `tfsdk:"is_published"`
	IsLocal                     types.Bool  `tfsdk:"is_local"`
	IsCached                    types.Bool  `tfsdk:"is_cached"`
}

type catalogResourceModel struct {
	// BASE
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	OwnerName   types.String `tfsdk:"owner_name"`

	// SPECIFIC RESOURCE
	StorageProfile  types.String `tfsdk:"storage_profile"`
	DeleteForce     types.Bool   `tfsdk:"delete_force"`
	DeleteRecursive types.Bool   `tfsdk:"delete_recursive"`
}

func catalogDatasourceAttributes() map[string]schemaD.Attribute {
	attr := catalogSchema().GetDataSource(context.Background()).Attributes

	for k, v := range attr {
		switch val := v.(type) {
		case schemaD.StringAttribute:
			val.Computed = true
			val.Optional = false
			val.Required = false
			attr[k] = val
		case schemaD.BoolAttribute:
			val.Computed = true
			val.Optional = false
			val.Required = false
			attr[k] = val
		case schemaD.Int64Attribute:
			val.Computed = true
			val.Optional = false
			val.Required = false
			attr[k] = val
		}
	}

	return attr
}

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
			MarkdownDescription: " manage a catalog in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: " retrieve information about a catalog in Cloud Avenue.",
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
					MarkdownDescription: " Required if `name` is not set.",
					Optional:            true,
					Computed:            true,
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
					MarkdownDescription: " Required if `id` is not set.",
					Optional:            true,
					Computed:            true,
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
					MarkdownDescription: "Include BIOS UUIDs and MAC addresses in the downloaded OVF package. Preserving the identity information limits the portability of the package and you should use it only when necessary.",
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
					MarkdownDescription: "When destroying use `delete_force=True` with `delete_recursive=True` to remove a catalog and any objects it contains, regardless of their state.",
				},
			},
			"delete_recursive": superschema.BoolAttribute{
				Resource: &schemaR.BoolAttribute{
					Required:            true,
					MarkdownDescription: "When destroying use `delete_recursive=True` to remove a catalog and any objects it contains that are in a state that normally allows removal.",
				},
			},
		},
	}
}

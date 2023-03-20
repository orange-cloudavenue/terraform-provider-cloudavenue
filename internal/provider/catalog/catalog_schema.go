package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/storageprofile"
)

type catalogDataSourceModel struct {
	// BASE
	ID          types.String `tfsdk:"id"`
	CatalogName types.String `tfsdk:"catalog_name"`
	CatalogID   types.String `tfsdk:"catalog_id"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	OwnerName   types.String `tfsdk:"owner_name"`
	Href        types.String `tfsdk:"href"`

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
	CatalogName types.String `tfsdk:"catalog_name"`
	CatalogID   types.String `tfsdk:"catalog_id"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	OwnerName   types.String `tfsdk:"owner_name"`
	Href        types.String `tfsdk:"href"`

	// SPECIFIC RESOURCE
	StorageProfile  types.String `tfsdk:"storage_profile"`
	DeleteForce     types.Bool   `tfsdk:"delete_force"`
	DeleteRecursive types.Bool   `tfsdk:"delete_recursive"`
}

type catalogSchemaOpts func(*catalogSchemaParams)

type catalogSchemaParams struct {
	resource   bool
	datasource bool
}

func withDataSource() catalogSchemaOpts {
	return func(params *catalogSchemaParams) {
		params.datasource = true
	}
}

/*
catalogSchema

This function is used to create the schema for the catalog resource and datasource.
Default is to create a resource schema.  If you want to create a datasource schema
you must pass in the withDataSource() option.
*/
func catalogSchema(opts ...catalogSchemaOpts) superschema.Schema {
	params := &catalogSchemaParams{}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(params)
		}
	} else {
		params.resource = true
	}

	_schema := superschema.Schema{}

	_schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The creation date of the catalog.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"description": schema.StringAttribute{
			Optional:            params.resource,
			Computed:            true,
			MarkdownDescription: "The description of the catalog.",
		},
		"href": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The catalog HREF.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"owner_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The name of the owner of the catalog.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}

	if params.datasource {
		_schema.MarkdownDescription = "The catalog data source show the details of the catalog."

		_schema.Attributes[schemaName] = schemaCatalogName(common.IsOptional())
		_schema.Attributes[schemaID] = schemaCatalogID(common.IsOptional())
		_schema.Attributes["preserve_identity_information"] = schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Include BIOS UUIDs and MAC addresses in the downloaded OVF package. Preserving the identity information limits the portability of the package and you should use it only when necessary.",
		}
		_schema.Attributes["number_of_media"] = schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "Number of Medias this catalog contains.",
		}
		_schema.Attributes["media_item_list"] = schema.ListAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			MarkdownDescription: "List of Media items in this catalog.",
		}
		_schema.Attributes["is_shared"] = schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this catalog is shared.",
		}
		_schema.Attributes["is_local"] = schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this catalog belongs to the current organization.",
		}
		_schema.Attributes["is_published"] = schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this catalog is shared to all organizations.",
		}
		_schema.Attributes["is_cached"] = schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "True if this catalog is cached.",
		}
	}

	if params.resource {
		_schema.MarkdownDescription = "The Catalog resource allows you to manage a catalog in CloudAvenue."

		_schema.Attributes[schemaName] = schemaCatalogName()
		_schema.Attributes[schemaID] = schemaCatalogID(common.IsComputed())
		_schema.Attributes[storageprofile.SchemaStorageProfile] = storageprofile.Schema()
		_schema.Attributes["delete_force"] = schema.BoolAttribute{
			Required:            true,
			MarkdownDescription: "When destroying use `delete_force=True` with `delete_recursive=True` to remove a catalog and any objects it contains, regardless of their state.",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		_schema.Attributes["delete_recursive"] = schema.BoolAttribute{
			Required:            true,
			MarkdownDescription: "When destroying use `delete_recursive=True` to remove the catalog and any objects it contains that are in a state that normally allows removal.",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
	}

	return _schema
}

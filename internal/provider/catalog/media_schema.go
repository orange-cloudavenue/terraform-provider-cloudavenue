package catalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
)

type catalogMediaDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	CatalogID      types.String `tfsdk:"catalog_id"`
	CatalogName    types.String `tfsdk:"catalog_name"`
	Description    types.String `tfsdk:"description"`
	IsISO          types.Bool   `tfsdk:"is_iso"`
	OwnerName      types.String `tfsdk:"owner_name"`
	IsPublished    types.Bool   `tfsdk:"is_published"`
	CreatedAt      types.String `tfsdk:"created_at"`
	Size           types.Int64  `tfsdk:"size"`
	Status         types.String `tfsdk:"status"`
	StorageProfile types.String `tfsdk:"storage_profile"`
}

func catalogMediaDataSourceModelType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"name":            types.StringType,
		"catalog_id":      types.StringType,
		"catalog_name":    types.StringType,
		"description":     types.StringType,
		"is_iso":          types.BoolType,
		"owner_name":      types.StringType,
		"is_published":    types.BoolType,
		"created_at":      types.StringType,
		"size":            types.Int64Type,
		"status":          types.StringType,
		"storage_profile": types.StringType,
	}
}

func mediaDatasourceAttributes() map[string]schemaD.Attribute {
	attr := mediaSchema().GetDataSource(context.Background()).Attributes

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

func mediaSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Catalog media allows you to",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: " manage a media in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: " retrieve information about a media in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the media.",
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
			"catalog_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the catalog.",
				},
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: " Required if `catalog_name` is not set.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_name"), path.MatchRoot("catalog_id")),
					},
				},
			},
			"catalog_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the catalog.",
				},
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: " Required if `catalog_id` is not set.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_name"), path.MatchRoot("catalog_id")),
					},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the media.",
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
			"description": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: " The description of the media.",
					Computed:            true,
				},
			},
			"is_iso": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "`True` if the media is an ISO.",
					Computed:            true,
				},
			},
			"owner_name": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: " The name of the owner of the media.",
					Computed:            true,
				},
			},
			"is_published": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "`True` if the media is published.",
					Computed:            true,
				},
			},
			"created_at": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: " The date and time when the media was created.",
					Computed:            true,
				},
			},
			"size": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: " The size of the media in bytes.",
					Computed:            true,
				},
			},
			"status": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: " The status of the media.",
					Computed:            true,
				},
			},
			"storage_profile": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: " The storage profile of the media.",
					Computed:            true,
				},
			},
		},
	}
}

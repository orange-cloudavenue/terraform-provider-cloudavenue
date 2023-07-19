package catalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func vappTemplateSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "The `catalog_vapp_template` datasource provides information about a vApp Template in a catalog.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the vApp Template",
				Computed:            true,
			},
			"template_name": schema.StringAttribute{
				MarkdownDescription: "Name of the vApp Template. Required if `template_id` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("template_name"), path.MatchRoot("template_id")),
				},
			},
			"template_id": schema.StringAttribute{
				MarkdownDescription: "ID of the vApp Template. Required if `template_name` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("template_name"), path.MatchRoot("template_id")),
				},
			},
			catalogID:   mediaSchema().GetDataSource(ctx).Attributes[catalogID],
			catalogName: mediaSchema().GetDataSource(ctx).Attributes[catalogName],
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the vApp Template",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation date of the vApp Template",
				Computed:            true,
			},
			"vm_names": schema.ListAttribute{
				MarkdownDescription: "Set of VM names within the vApp template",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

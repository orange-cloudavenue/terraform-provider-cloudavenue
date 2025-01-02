package catalog

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func vappTemplateSuperSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `catalog_vapp_template` datasource provides information about a vApp Template in a catalog.",
		},
		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "ID of the vApp Template",
					Computed:            true,
				},
			},
			"description": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Description of the vApp Template",
					Computed:            true,
				},
			},
			"template_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The Name of the vApp Template.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("template_name"), path.MatchRoot("template_id")),
					},
				},
			},
			"template_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The ID of the vApp Template.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("template_name"), path.MatchRoot("template_id")),
					},
				},
			},
			"catalog_id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The ID of the catalog.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_name"), path.MatchRoot("catalog_id")),
					},
				},
			},
			"catalog_name": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The Name of the catalog.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_name"), path.MatchRoot("catalog_id")),
					},
				},
			},
			"vm_names": superschema.SuperSetAttribute{
				DataSource: &schemaD.SetAttribute{
					MarkdownDescription: "Set of VM names within the vApp template",
					Computed:            true,
					ElementType:         supertypes.StringType{},
				},
			},
			"created_at": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Creation date of the vApp Template",
					Computed:            true,
				},
			},
		},
	}
}

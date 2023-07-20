package iam

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func iamRightSchema() schema.Schema {
	return schema.Schema{
		Description: "Provides a data source for available rights in Cloud Avenue.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The id of the right.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the right.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description for the right.",
			},
			"category_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The category id for the right.",
			},
			// * Remove
			"bundle_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The bundle key for the right.",
			},
			"right_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The right type for the right.",
			},
			"implied_rights": schema.SetNestedAttribute{
				MarkdownDescription: "The list of rights that are implied with this one.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the implied right.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "ID of the implied right.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

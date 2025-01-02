package vrf

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func tier0VrfsSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Tier-0 VRFs",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allow access to a list of Tier-0 that can be accessed by the user.",
		},

		Attributes: superschema.Attributes{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Tier-0 VRFs.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"names": superschema.SuperListAttributeOf[string]{
				Common: &schemaR.ListAttribute{
					MarkdownDescription: "List of Tier-0 VRFs names.",
				},
				DataSource: &schemaD.ListAttribute{
					ElementType: supertypes.StringType{},
					Computed:    true,
				},
			},
		},
	}
}

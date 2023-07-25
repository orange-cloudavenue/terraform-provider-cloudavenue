package vrf

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

func tier0VrfsSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Tier-0 VRFs",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allow access to a list of Tier-0 that can be accessed by the user.",
		},

		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Tier-0 VRFs.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"names": superschema.ListAttribute{
				Common: &schemaR.ListAttribute{
					MarkdownDescription: "List of Tier-0 VRFs names.",
				},
				DataSource: &schemaD.ListAttribute{
					ElementType: types.StringType,
					Computed:    true,
				},
			},
		},
	}
}

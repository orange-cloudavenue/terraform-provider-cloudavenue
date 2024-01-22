package vrf

import (
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

func tier0VrfSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Tier-0 VRF",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source retrieve informations about a Tier-0 VRF.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Tier-0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Tier-0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"tier0_provider": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Tier-0 provider info.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"class_service": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "List of Tags for the Tier-0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"services": superschema.SuperListNestedAttributeOf[segmentModel]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "Services list of the Tier-0 VRF.",
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"service": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Service of the segment.",
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"vlan_id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "VLAN ID of the segment.",
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

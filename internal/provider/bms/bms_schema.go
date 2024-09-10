package bms

import (
	"context"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

// TODO : Remove unused imports.
// ! This is outside import block because golangci-lint remove commented import.
// * Hashicorp Validators
// "github.com/Hashicorp/terraform-plugin-framework-validators/stringvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/boolvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/int64validator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/float64validator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/listvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/mapvalidator"
// "github.com/Hashicorp/terraform-plugin-framework-validators/setvalidator"

// * Hashicorp Plan Modifiers Resource
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"

// * Hashicorp Plan Modifiers DataSource
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/stringplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/boolplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/int64planmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/listplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/mapplanmodifier"
// "github.com/hashicorp/terraform-plugin-framework/datasource/schema/setplanmodifier"

// * Hashicorp Default Values
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

// * FrangipaneTeam Custom Validators
// fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
// fboolvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/boolvalidator"
// fint64validator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/int64validator"
// flistvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/listvalidator"
// fmapvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/mapvalidator"
// fsetvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/setvalidator"

// * FrangipaneTeam Custom Plan Modifiers
// fstringplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/stringplanmodifier"
// fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"
// fint64planmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/int64planmodifier"
// flistplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/listplanmodifier"
// fmapplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/mapplanmodifier"
// fsetplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/setplanmodifier"

// How to use types generator:
// 1. Define the schema in the file internal/provider/bms/bms_schema.go
// 2. Add the resource or data source to the file internal/provider/provider_resources.go or internal/provider/provider_data_sources.go respectively
// 3. Launch the following command to generate golang structs for the schema:
// go run ./cmd/types-generator/*.go -file internal/provider/bms/datasource_schema.go -resource cloudavenue_bms_datasource -is-resource.
func bmsSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_bms_datasource` resource allows you to manage ...",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_bms_datasource` data source allows you to retrieve information about an ...",
		},
		Attributes: map[string]superschema.Attribute{
			"timeouts": superschema.TimeoutAttribute{
				// Resource: &superschema.ResourceTimeoutAttribute{
				// 	Create: true,
				// 	Update: true,
				// 	Delete: true,
				// 	Read:   true,
				// },
				DataSource: &superschema.DatasourceTimeoutAttribute{
					Read: true,
				},
			},
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the datasource.",
				},
			},
			"network": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceNetwork]{
				DataSource: &schemaD.SetNestedAttribute{
					Computed:            true,
					MarkdownDescription: "The network of the BMS list.",
				},
				Attributes: superschema.Attributes{
					"vlan_id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The VLAN ID of the network.",
						},
					},
					"subnet": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The subnet of the network.",
						},
					},
					"prefix": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The prefix of the network.",
						},
					},
				},
			},

			"bms": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMS]{
				DataSource: &schemaD.SetNestedAttribute{
					Computed:            true,
					MarkdownDescription: "The BMS list.",
				},
				Attributes: superschema.Attributes{
					"hostname": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The hostname of the BMS.",
						},
					},
					"bms_type": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the BMS.",
						},
					},
					"os": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OS of the BMS.",
						},
					},
					"bios_configuration": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The BIOS configuration of the BMS.",
						},
					},
					"storage": superschema.SuperSingleNestedAttributeOf[bmsModelDatasourceBMSStorage]{
						DataSource: &schemaD.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The storage of the BMS.",
						},
						Attributes: superschema.Attributes{
							"local": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMSStorageGen]{
								DataSource: &schemaD.SetNestedAttribute{
									Computed:            true,
									MarkdownDescription: "The local storage of the BMS.",
								},
								Attributes: superschema.Attributes{
									"size": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The size of the local storage.",
										},
									},
									"storage_class": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The storage class of the local storage.",
										},
									},
								},
							},
							"system": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMSStorageGen]{
								DataSource: &schemaD.SetNestedAttribute{
									Computed:            true,
									MarkdownDescription: "The system storage of the BMS.",
								},
								Attributes: superschema.Attributes{
									"size": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The size of the system storage.",
										},
									},
									"storage_class": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The storage class of the system storage.",
										},
									},
								},
							},
							"data": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMSStorageGen]{
								DataSource: &schemaD.SetNestedAttribute{
									Computed:            true,
									MarkdownDescription: "The data storage of the BMS.",
								},
								Attributes: superschema.Attributes{
									"size": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The size of the data storage.",
										},
									},
									"storage_class": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The storage class of the data storage.",
										},
									},
								},
							},
							"shared": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMSStorageGen]{
								DataSource: &schemaD.SetNestedAttribute{
									Computed:            true,
									MarkdownDescription: "The shared storage of the BMS.",
								},
								Attributes: superschema.Attributes{
									"size": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The size of the shared storage.",
										},
									},
									"storage_class": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The storage class of the shared storage.",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

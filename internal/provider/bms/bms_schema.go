/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package bms

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func bmsSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_bms` data source allows you to retrieve information about your Bare Metal Server.",
		},
		Attributes: map[string]superschema.Attribute{
			"timeouts": superschema.TimeoutAttribute{
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
			"env": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceEnv]{
				DataSource: &schemaD.SetNestedAttribute{
					Computed:            true,
					MarkdownDescription: "Return the list of BMS environement.",
				},
				Attributes: map[string]superschema.Attribute{
					"network": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceNetwork]{
						DataSource: &schemaD.SetNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The network array for all BMS listed.",
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
							"type": superschema.SuperStringAttribute{
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
									"local": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMSStorageDetail]{
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
									"system": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMSStorageDetail]{
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
									"data": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMSStorageDetail]{
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
									"shared": superschema.SuperSetNestedAttributeOf[bmsModelDatasourceBMSStorageDetail]{
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
			},
		},
	}
}

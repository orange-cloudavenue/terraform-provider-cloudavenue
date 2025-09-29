/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"context"

	"github.com/orange-cloudavenue/common-go/urn"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/validators"
)

func servicesSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_edgegateway_services` resource allows you to manage Edge Gateway network services in CloudAvenue. It provides the foundational configuration for CloudAvenue services. For more details, refer to the [CloudAvenue documentation (FR)](https://cloud.orange-business.com/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/fiches-pratiques/zone-de-services/se-connecter-a-la-zone-de-services/presentation-de-la-zone-de-service-ngp/). To configure services, you can use the Terraform module [terraform-cloudavenue-edgegateway-services](https://registry.terraform.io/modules/orange-cloudavenue/edgegateway-services/cloudavenue/latest).",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_edgegateway_services` data source provides details about Edge Gateway network services in CloudAvenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Unique identifier of the services.",
					Computed:            true,
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Name of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
						validators.ResourceName("edgegateway"),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Identifier of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
						fstringvalidator.IsURN(),
						fstringvalidator.PrefixContains(urn.EdgeGateway.String()),
					},
				},
			},
			"network": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Allocated network for the CloudAvenue services.",
					Computed:            true,
				},
			},
			"ip_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Dedicated IP for the service (used for the NAT to connect to the service).",
					Computed:            true,
				},
			},
			"services": superschema.SuperMapNestedAttributeOf[ServicesModelCatalog]{
				Common: &schemaR.MapNestedAttribute{
					MarkdownDescription: "Collection of services.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"network": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Network for accessing this category of services.",
							Computed:            true,
						},
					},
					"category": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Category of services.",
							Computed:            true,
						},
					},
					"services": superschema.SuperMapNestedAttributeOf[ServicesModelCatalogService]{
						Common: &schemaR.MapNestedAttribute{
							MarkdownDescription: "Details of individual services.",
							Computed:            true,
						},
						Attributes: superschema.Attributes{
							"name": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Name of the service.",
									Computed:            true,
								},
							},
							"description": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Description of the service.",
									Computed:            true,
								},
							},
							"ips": superschema.SuperListAttributeOf[string]{
								Common: &schemaR.ListAttribute{
									MarkdownDescription: "List of IP addresses associated with the service.",
									Computed:            true,
								},
							},
							"fqdns": superschema.SuperListAttributeOf[string]{
								Common: &schemaR.ListAttribute{
									MarkdownDescription: "List of FQDNs associated with the service.",
									Computed:            true,
								},
							},
							"ports": superschema.SuperListNestedAttributeOf[ServicesModelCatalogServicePorts]{
								Common: &schemaR.ListNestedAttribute{
									MarkdownDescription: "List of ports used by the service.",
									Computed:            true,
								},
								Attributes: superschema.Attributes{
									"port": superschema.SuperInt32Attribute{
										Common: &schemaR.Int32Attribute{
											MarkdownDescription: "Port number used by the service.",
											Computed:            true,
										},
									},
									"protocol": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Protocol used by the service (e.g., TCP, UDP).",
											Computed:            true,
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

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package elb

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
)

func virtualServiceSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a resource to manage ELB Virtual services for particular Gateway. A virtual service advertises an IP address and ports to the external world and listens for client traffic. When a virtual service receives traffic, it directs it to members in ELB Pool.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a data source to read ELB Virtual services for particular Gateway. A virtual service advertises an IP address and ports to the external world and listens for client traffic. When a virtual service receives traffic, it directs it to members in ELB Pool.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the ELB virtual service.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the ELB Virtual Service.",
					Required:            true,
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the edge gateway on which the ELB Virtual Service is to be created.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the edge gateway on which the ELB Virtual Service is to be created.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the ELB Virtual Service.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Defines if the ELB Virtual Service is enabled.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(true),
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"pool_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the ELB Server Pool associated.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("pool_name"), path.MatchRoot("pool_id")),
					},
				},
			},
			"pool_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the ELB Server Pool associated.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("pool_name"), path.MatchRoot("pool_id")),
					},
				},
			},
			"service_engine_group_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the service Engine Group (Take the first one if not specified).",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
			},
			"virtual_ip": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The virtual IP address of the ELB Virtual Service.",
				},
				Resource: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "\n\n -> The `virtual_ip` accept a private IP in your network range or a public IP (Warning: the public IP must be reserved and not used by any other service).",
					Validators: []validator.String{
						fstringvalidator.IsNetwork([]fstringvalidator.NetworkValidatorType{fstringvalidator.IPV4}, false),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"service_type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of the service. The different modes that the ELB supports for handling TCP traffic and various parameters that can be tuned for optimization of the TCP traffic are also detailed here.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(func() []string {
							var values []string
							for _, v := range edgeloadbalancer.VirtualServiceApplicationProfiles {
								values = append(values, string(v))
							}
							return values
						}()...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"certificate_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the certificate. The certificate must be uploaded to your certificate library before it can be used. The certificate MUSTN'T be expired.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("service_type"), []attr.Value{types.StringValue("L4_TLS"), types.StringValue("HTTPS")}),
						fstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("service_type"), []attr.Value{types.StringValue("HTTP"), types.StringValue("L4_TCP"), types.StringValue("L4_UDP")}),
						fstringvalidator.PrefixContains(urn.CertificateLibraryItem.String()),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"service_ports": superschema.SuperListNestedAttributeOf[VirtualServiceModelServicePort]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "The service port of the ELB Virtual Service. The service port is the port on which the virtual service listens for client traffic.",
				},
				Resource: &schemaR.ListNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"start": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The start port of the service port range or exact port number if `end` is not set.",
						},
						Resource: &schemaR.Int64Attribute{
							Required: true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"end": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The end port of the service port range. If not specified, only the `start` value is used.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
					},
				},
			},
		},
	}
}

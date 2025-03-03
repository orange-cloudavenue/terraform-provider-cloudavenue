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
	"slices"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
)

func poolSchema(ctx context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_elb_pool` resource allows you to manage edgegateway load balancer pools. Pools maintain the list of servers assigned to them and perform health monitoring, load balancing, persistence. A pool may only be used or referenced by only one virtual service at a time.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_elb_pool` data source allows you to retrieve information about an existing edgegateway load balancer pool.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the pool.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the pool.",
					Required:            true,
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the pool.",
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
					MarkdownDescription: "Enable or disable the pool.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Default:  booldefault.StaticBool(true),
				},
			},
			"algorithm": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The heart of a load balancer is its ability to effectively distribute traffic across healthy servers. If persistence is enabled, only the first connection from a client is load balanced. While the persistence remains in effect, subsequent connections or requests from a client are directed to the same server.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Default:  stringdefault.StaticString(string(edgeloadbalancer.PoolAlgorithmLeastConnections)),
					Validators: []validator.String{
						fstringvalidator.OneOfWithDescription(func() (resp []fstringvalidator.OneOfWithDescriptionValues) {
							x := []string{}

							for k := range edgeloadbalancer.PoolAlgorithms {
								x = append(x, string(k))
							}

							slices.Sort(x)

							for _, v := range x {
								resp = append(resp, fstringvalidator.OneOfWithDescriptionValues{
									Value:       v,
									Description: edgeloadbalancer.PoolAlgorithms[edgeloadbalancer.PoolAlgorithm(v)],
								})
							}
							return
						}()...),
					},
				},
			},
			"default_port": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "DefaultPort defines destination server port used by the traffic sent to the member.",
				},
				Resource: &schemaR.Int64Attribute{
					Required: true,
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"members": superschema.SuperSingleNestedAttributeOf[PoolModelMembers]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The members of the pool.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"graceful_timeout_period": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Maximum time (in minutes) to gracefully disable a member. Virtual service waits for the specified time before terminating the existing connections to the members that are disabled. Special values: `0` represents `Immediate` and `-1` represents `Infinite`. The maximum value is `7200` minutes.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Default:  stringdefault.StaticString("1"),
						},
					},
					"target_group": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The group contains reference to the Edge Firewall Group representing destination servers which are used by the Load Balancer Pool to direct load balanced traffic. This permit to reference `IP Set` or `Static Group` ID.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								fstringvalidator.PrefixContains(urn.SecurityGroup.String()),
								stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("targets"), path.MatchRelative().AtParent().AtName("target_group")),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"targets": superschema.SuperListNestedAttributeOf[PoolModelMembersIPAddress]{
						Common: &schemaR.ListNestedAttribute{
							MarkdownDescription: "targets field defines list of destination servers which are used by the Load Balancer Pool to direct load balanced traffic.",
						},
						Resource: &schemaR.ListNestedAttribute{
							Optional: true,
							Validators: []validator.List{
								listvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("targets"), path.MatchRelative().AtParent().AtName("target_group")),
							},
						},
						DataSource: &schemaD.ListNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"enabled": superschema.SuperBoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Enable or disable the member.",
									Computed:            true,
								},
								Resource: &schemaR.BoolAttribute{
									Optional: true,
									Default:  booldefault.StaticBool(true),
								},
							},
							"ip_address": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The IP address of the member.",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										fstringvalidator.IsNetwork([]fstringvalidator.NetworkValidatorType{
											fstringvalidator.IPV4,
										},
											false,
										),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"port": superschema.SuperInt64Attribute{
								Common: &schemaR.Int64Attribute{
									MarkdownDescription: "The port of the member.",
								},
								Resource: &schemaR.Int64Attribute{
									Required: true,
								},
								DataSource: &schemaD.Int64Attribute{
									Computed: true,
								},
							},
							"ratio": superschema.SuperInt64Attribute{
								Common: &schemaR.Int64Attribute{
									MarkdownDescription: "The ratio of the member. The ratio of each pool member denotes the traffic that goes to each server pool member. A server with a ratio of 2 gets twice as much traffic as a server with a ratio of 1.",
									Computed:            true,
								},
								Resource: &schemaR.Int64Attribute{
									Optional: true,
									Default:  int64default.StaticInt64(1),
								},
							},
						},
					},
				},
			},
			"health": superschema.SuperSingleNestedAttributeOf[PoolModelHealth]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Health check member servers health. It can be monitored by using one or more health monitors. Active monitors generate synthetic traffic and mark a server up or down based on the response.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Default: objectdefault.StaticValue(supertypes.NewObjectValueOf[PoolModelHealth](ctx, &PoolModelHealth{
						PassiveMonitoringEnabled: supertypes.NewBoolValue(true),
						Monitors:                 supertypes.NewListValueOfNull[string](ctx),
					}).ObjectValue),
				},
				Attributes: superschema.Attributes{
					"passive_monitoring_enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "PassiveMonitoringEnabled sets if client traffic should be used to check if pool member is up or down.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(true),
						},
					},
					"monitors": superschema.SuperListAttributeOf[string]{
						Common: &schemaR.ListAttribute{
							MarkdownDescription: "The active health monitors.",
						},
						Resource: &schemaR.ListAttribute{
							Optional: true,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(
									stringvalidator.OneOf(func() (resp []string) {
										for _, v := range edgeloadbalancer.PoolHealthMonitorTypes {
											resp = append(resp, string(v))
										}
										slices.Sort(resp)
										return
									}()...),
								),
							},
						},
						DataSource: &schemaD.ListAttribute{
							Computed: true,
						},
					},
				},
			},
			"tls": superschema.SuperSingleNestedAttributeOf[PoolModelTLS]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The TLS configuration of the pool.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Default: objectdefault.StaticValue(supertypes.NewObjectValueOf[PoolModelTLS](ctx, &PoolModelTLS{
						Enabled:                supertypes.NewBoolValue(false),
						DomainNames:            supertypes.NewListValueOfNull[string](ctx),
						CaCertificateRefs:      supertypes.NewListValueOfNull[string](ctx),
						CommonNameCheckEnabled: supertypes.NewBoolValue(false),
					}).ObjectValue),
				},
				Attributes: superschema.Attributes{
					"enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Enable or disable the TLS.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
						},
					},
					"domain_names": superschema.SuperListAttributeOf[string]{
						Common: &schemaR.ListAttribute{
							MarkdownDescription: "The domain names of the TLS check. This attribute is taken into account if the `common_name_check_enabled` is set to `true`.",
						},
						Resource: &schemaR.ListAttribute{
							Optional: true,
							Validators: []validator.List{
								listvalidator.SizeBetween(0, 10),
							},
						},
						DataSource: &schemaD.ListAttribute{
							Computed: true,
						},
					},
					"ca_certificate_refs": superschema.SuperListAttributeOf[string]{
						Common: &schemaR.ListAttribute{
							MarkdownDescription: "The CA certificate references point to root certificates to use when validating certificates presented by the pool members.",
						},
						Resource: &schemaR.ListAttribute{
							MarkdownDescription: "Use `cloudavenue_org_certificate` resource to create a certificate and get the ID.",
							Optional:            true,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(
									fstringvalidator.PrefixContains(urn.CertificateLibraryItem.String()),
								),
							},
						},
						DataSource: &schemaD.ListAttribute{
							Computed: true,
						},
					},
					"common_name_check_enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Enable common name check for server certificate. If enabled and no explicit domain name is specified, the incoming host header will be used to do the match.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
						},
					},
				},
			},
			"persistence": superschema.SuperSingleNestedAttributeOf[PoolModelPersistence]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Persistence profile will ensure that the same user sticks to the same server for a desired duration of time. If the persistence profile is unmanaged by ELB, updates that leave the values unchanged will continue to use the same unmanaged profile. Any changes made to the persistence profile will cause ELB to switch the pool to a profile managed by ELB.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Default: objectdefault.StaticValue(supertypes.NewObjectValueOf[PoolModelPersistence](ctx, &PoolModelPersistence{
						Type:  supertypes.NewStringValue(string(edgeloadbalancer.PoolPersistenceProfileTypeClientIP)),
						Value: supertypes.NewStringNull(),
					}).ObjectValue),
				},
				Attributes: superschema.Attributes{
					"type": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the persistence.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Default:  stringdefault.StaticString(string(edgeloadbalancer.PoolPersistenceProfileTypeClientIP)),
							Validators: []validator.String{
								fstringvalidator.OneOfWithDescription(func() (resp []fstringvalidator.OneOfWithDescriptionValues) {
									x := []string{}

									for k := range edgeloadbalancer.PoolPersistenceProfileTypes {
										x = append(x, string(k))
									}

									slices.Sort(x)

									for _, v := range x {
										resp = append(resp, fstringvalidator.OneOfWithDescriptionValues{
											Value:       v,
											Description: edgeloadbalancer.PoolPersistenceProfileTypes[edgeloadbalancer.PoolPersistenceProfileType(v)],
										})
									}
									return
								}()...),
							},
						},
					},
					"value": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The value of the persistence.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("type"), func() (resp []attr.Value) {
									resp = make([]attr.Value, 0)
									resp = append(resp, types.StringValue(string(edgeloadbalancer.PoolPersistenceProfileTypeHTTPCookie)))
									resp = append(resp, types.StringValue(string(edgeloadbalancer.PoolPersistenceProfileTypeCustomHTTPHeader)))
									resp = append(resp, types.StringValue(string(edgeloadbalancer.PoolPersistenceProfileTypeAPPCookie)))
									return
								}()),
							},
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

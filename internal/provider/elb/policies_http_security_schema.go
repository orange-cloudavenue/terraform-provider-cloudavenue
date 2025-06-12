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
	fboolvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/boolvalidator"
	fintvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/int64validator"
	"github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	stringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
)

func policiesHTTPSecuritySchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_elb_policies_http_security` resource allows you to manage HTTP security policies. HTTP security rules modify requests before they are either forwarded to the application, used as a basis for content switching, or discarded.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_elb_policies_http_security` data source allows you to retrieve information about an existing HTTP security policies.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the policies http security.",
				},
			},
			"virtual_service_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the virtual service to which the policies http security belongs.",
					Required:            true,
					Validators: []validator.String{
						fstringvalidator.Formats(
							[]fstringvalidator.FormatsValidatorType{
								fstringvalidator.FormatsIsURN,
							},
							false,
						),
						fstringvalidator.PrefixContains(string(urn.LoadBalancerVirtualService)),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"policies": superschema.SuperListNestedAttributeOf[PoliciesHTTPSecurityModelPolicies]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "HTTP security policies.",
				},
				Resource: &schemaR.ListNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Policy name, it must be unique within the virtual service's HTTP security policies.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"active": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether the policy is enable or not.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(true),
						},
					},
					"logging": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether to enable logging with headers on rule match or not.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
						},
					},
					"criteria": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPSecurityMatchCriteria]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Match criteria for the HTTP security. The criteria is used to match the request and determine if the action should be applied.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"protocol": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Protocol to match. Only HTTP application layer protocol (OSI 7) are supported.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPMatchCriteriaProtocolsString...),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"client_ip": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPClientIPMatch]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Match the rule based on client IP address rules.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"criteria": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Criteria to match.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPClientIPMatchCriteriaString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"ip_addresses": superschema.SuperSetAttributeOf[string]{
										Common: &schemaR.SetAttribute{
											MarkdownDescription: "IP addresses to match.",
										},
										Resource: &schemaR.SetAttribute{
											Required: true,
											Validators: []validator.Set{
												setvalidator.ValueStringsAre(
													fstringvalidator.IsNetwork([]fstringvalidator.NetworkValidatorType{
														fstringvalidator.IPV4,
														fstringvalidator.IPV4WithCIDR,
														fstringvalidator.IPV4Range,
													}, true),
												),
											},
										},
										DataSource: &schemaD.SetAttribute{
											Computed: true,
										},
									},
								},
							}, // End client_ip
							"service_ports": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPServicePortMatch]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Match the rule based on service port rules.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"criteria": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Criteria to match.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPServicePortMatchCriteriaString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"ports": superschema.SuperSetAttributeOf[int64]{
										Common: &schemaR.SetAttribute{
											MarkdownDescription: "A port list allows you to define which service ports (e.g.: [80, 443] ) the HTTP security policy should match.",
										},
										Resource: &schemaR.SetAttribute{
											Required: true,
											Validators: []validator.Set{
												setvalidator.ValueInt64sAre(
													int64validator.Between(1, 65535),
												),
											},
										},
										DataSource: &schemaD.SetAttribute{
											Computed: true,
										},
									},
								},
							}, // End service_port
							"http_methods": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPMethodMatch]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Match the rule based on HTTP method rules.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"criteria": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Criteria to match.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPMethodMatchCriteriaString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"methods": superschema.SuperSetAttributeOf[string]{
										Common: &schemaR.SetAttribute{
											MarkdownDescription: "Methods to match.",
										},
										Resource: &schemaR.SetAttribute{
											Required: true,
											Validators: []validator.Set{
												setvalidator.ValueStringsAre(
													stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPMethodsMatchString...),
												),
											},
										},
										DataSource: &schemaD.SetAttribute{
											Computed: true,
										},
									},
								},
							}, // End method
							"path": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPPathMatch]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Match the rule based on path rules.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"criteria": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Criteria to match.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPPathMatchCriteriaString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"paths": superschema.SuperSetAttributeOf[string]{
										Common: &schemaR.SetAttribute{
											MarkdownDescription: "A set of paths to match given criteria.",
										},
										Resource: &schemaR.SetAttribute{
											Required: true,
										},
										DataSource: &schemaD.SetAttribute{
											Computed: true,
										},
									},
								},
							}, // End path
							"cookie": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPCookieMatch]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Match the rule based on cookie rules.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"criteria": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Criteria to match.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPCookieMatchCriteriaString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"name": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Name of the cookie to match.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"value": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Value of the cookie to match.",
										},
										Resource: &schemaR.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												fstringvalidator.NullIfAttributeIsOneOf(
													path.MatchRelative().AtParent().AtName("criteria"),
													[]attr.Value{
														types.StringValue(string(edgeloadbalancer.PoliciesHTTPMatchCriteriaCriteriaEXISTS)),
														types.StringValue(string(edgeloadbalancer.PoliciesHTTPMatchCriteriaCriteriaDOESNOTEXIST)),
													},
												),
												fstringvalidator.RequireIfAttributeIsOneOf(
													path.MatchRelative().AtParent().AtName("criteria"),
													func() []attr.Value {
														x := make([]attr.Value, 0)
														for _, v := range edgeloadbalancer.PoliciesHTTPCookieMatchCriteriaString {
															if v != string(edgeloadbalancer.PoliciesHTTPMatchCriteriaCriteriaEXISTS) && v != string(edgeloadbalancer.PoliciesHTTPMatchCriteriaCriteriaDOESNOTEXIST) {
																x = append(x, types.StringValue(v))
															}
														}
														return x
													}(),
												),
												stringvalidator.LengthAtMost(10240),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
								},
							}, // End cookie
							"request_headers": superschema.SuperSetNestedAttributeOf[PoliciesHTTPHeaderMatch]{
								Common: &schemaR.SetNestedAttribute{
									MarkdownDescription: "Match the rule based on request headers rules.",
								},
								Resource: &schemaR.SetNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SetNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"criteria": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Criteria to match.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												// Use the same criteria as for HTTP headers Request
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPHeaderMatchCriteriaString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"name": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Name of the HTTP header whose value is to be matched.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"values": superschema.SuperSetAttributeOf[string]{
										Common: &schemaR.SetAttribute{
											MarkdownDescription: "Values of the HTTP header to match.",
										},
										Resource: &schemaR.SetAttribute{
											Optional: true,
											// ref: https://github.com/orange-cloudavenue/terraform-plugin-framework-validators/issues/71
											// Validators: []validator.Set{
											// 	fsetvalidator.NullIfAttributeIsOneOf(
											// 		path.MatchRelative().AtParent().AtName("criteria"),
											// 		[]attr.Value{
											// 			types.StringValue(string(edgeloadbalancer.PoliciesHTTPMatchCriteriaCriteriaEXISTS)),
											// 			types.StringValue(string(edgeloadbalancer.PoliciesHTTPMatchCriteriaCriteriaDOESNOTEXIST)),
											// 		},
											// 	),
											// 	fsetvalidator.RequireIfAttributeIsOneOf(
											// 		path.MatchRelative().AtParent().AtName("criteria"),
											// 		func() []attr.Value {
											// 			x := make([]attr.Value, 0)
											// 			for _, v := range edgeloadbalancer.PoliciesHTTPHeaderMatchCriteriaString {
											// 				if v != string(edgeloadbalancer.PoliciesHTTPMatchCriteriaCriteriaEXISTS) && v != string(edgeloadbalancer.PoliciesHTTPMatchCriteriaCriteriaDOESNOTEXIST) {
											// 					x = append(x, types.StringValue(v))
											// 				}
											// 			}
											// 			return x
											// 		}(),
											// 	),
											// },
										},
										DataSource: &schemaD.SetAttribute{
											Computed: true,
										},
									},
								},
							}, // End security_headers
							"query": superschema.SuperSetAttributeOf[string]{
								Common: &schemaR.SetAttribute{
									MarkdownDescription: "Text contained in the query string",
								},
								Resource: &schemaR.SetAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SetAttribute{
									Computed: true,
								},
							},
						},
					}, // End criteria
					"actions": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPSecurityActions]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Actions to perform when the rule matches.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"connection": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Connection action to perform.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPConnectionActionsString...),
										fstringvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("redirect_to_https")),
										fstringvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("send_response")),
										fstringvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("rate_limit")),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"redirect_to_https": superschema.SuperInt64Attribute{
								Common: &schemaR.Int64Attribute{
									MarkdownDescription: "A port number, when set, configures the rule to redirect matching HTTP requests to HTTPS on the specified port.",
								},
								Resource: &schemaR.Int64Attribute{
									Optional: true,
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
										fintvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("connection")),
										fintvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("send_response")),
										fintvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("rate_limit")),
									},
								},
								DataSource: &schemaD.Int64Attribute{
									Computed: true,
								},
							},
							"send_response": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPActionSendResponse]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Send a customized response.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
									Validators: []validator.Object{
										objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("redirect_to_https")),
										objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("connection")),
										objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("rate_limit")),
									},
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"status_code": superschema.SuperInt64Attribute{
										Common: &schemaR.Int64Attribute{
											MarkdownDescription: "HTTP status code to return.",
											Computed:            true,
										},
										Resource: &schemaR.Int64Attribute{
											Optional: true,
											Validators: []validator.Int64{
												int64validator.OneOf(edgeloadbalancer.PoliciesHTTPActionResponseStatusCodes...),
											},
											Default: int64default.StaticInt64(200),
										},
									},
									"content": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Content of the response.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.LengthAtMost(10240),
												fstringvalidator.Formats(
													[]fstringvalidator.FormatsValidatorType{
														fstringvalidator.FormatsIsBase64,
													},
													false,
												),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"content_type": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Mime type of content.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPActionContentTypesString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
								},
							},
							"rate_limit": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPActionRateLimit]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "The rate_limit allows you to specify an action to take when the rate limit is reached. A rate limit defines the maximum number of requests permitted within a specific time frame.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
									Validators: []validator.Object{
										objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("redirect_to_https")),
										objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("connection")),
										objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("send_response")),
									},
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"count": superschema.SuperInt64Attribute{
										Common: &schemaR.Int64Attribute{
											MarkdownDescription: "Number of requests.",
											Computed:            true,
										},
										Resource: &schemaR.Int64Attribute{
											Optional: true,
											Validators: []validator.Int64{
												int64validator.Between(1, 1000),
											},
											Default: int64default.StaticInt64(100),
										},
									},
									"period": superschema.SuperInt64Attribute{
										Common: &schemaR.Int64Attribute{
											MarkdownDescription: "Period in seconds.",
											Computed:            true,
										},
										Resource: &schemaR.Int64Attribute{
											Optional: true,
											Validators: []validator.Int64{
												int64validator.Between(1, 1000000000),
											},
											Default: int64default.StaticInt64(60),
										},
									},
									"redirect": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPActionRedirect]{
										Common: &schemaR.SingleNestedAttribute{
											MarkdownDescription: "Redirects the request to different location when the rate limit is reached.",
										},
										Resource: &schemaR.SingleNestedAttribute{
											Optional: true,
											Validators: []validator.Object{
												objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("local_response")),
												objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("close_connection")),
											},
										},
										DataSource: &schemaD.SingleNestedAttribute{
											Computed: true,
										},
										Attributes: superschema.Attributes{
											"host": superschema.SuperStringAttribute{
												Common: &schemaR.StringAttribute{
													MarkdownDescription: "Host to which redirect the request. Default is the original host",
												},
												Resource: &schemaR.StringAttribute{
													Optional: true,
												},
												DataSource: &schemaD.StringAttribute{
													Computed: true,
												},
											},
											"keep_query": superschema.SuperBoolAttribute{
												Common: &schemaR.BoolAttribute{
													MarkdownDescription: "Keep or drop the query of the incoming request URI in the redirected URI",
													Computed:            true,
												},
												Resource: &schemaR.BoolAttribute{
													Optional: true,
													Default:  booldefault.StaticBool(true),
												},
											},
											"path": superschema.SuperStringAttribute{
												Common: &schemaR.StringAttribute{
													MarkdownDescription: "Path to which redirect the request. Default is the original path",
												},
												Resource: &schemaR.StringAttribute{
													Optional: true,
												},
												DataSource: &schemaD.StringAttribute{
													Computed: true,
												},
											},
											"port": superschema.SuperInt64Attribute{
												Common: &schemaR.Int64Attribute{
													MarkdownDescription: "Port to which redirect the request.",
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
											"protocol": superschema.SuperStringAttribute{
												Common: &schemaR.StringAttribute{
													MarkdownDescription: "HTTP protocol",
												},
												Resource: &schemaR.StringAttribute{
													Required: true,
													Validators: []validator.String{
														stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPProtocolsString...),
													},
												},
												DataSource: &schemaD.StringAttribute{
													Computed: true,
												},
											},
											"status_code": superschema.SuperInt64Attribute{
												Common: &schemaR.Int64Attribute{
													MarkdownDescription: "Redirect status code",
													Computed:            true,
												},
												Resource: &schemaR.Int64Attribute{
													Optional: true,
													Validators: []validator.Int64{
														int64validator.OneOf(edgeloadbalancer.PoliciesHTTPRedirectStatusCodes...),
													},
													Default: int64default.StaticInt64(302),
												},
											},
										},
									}, // End redirect
									"local_response": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPActionSendResponse]{
										Common: &schemaR.SingleNestedAttribute{
											MarkdownDescription: "Local response action can be used to send a customized response when the rate limit is reached.",
										},
										Resource: &schemaR.SingleNestedAttribute{
											Optional: true,
											Validators: []validator.Object{
												objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("redirect")),
												objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("close_connection")),
											},
										},
										DataSource: &schemaD.SingleNestedAttribute{
											Computed: true,
										},
										Attributes: superschema.Attributes{
											"status_code": superschema.SuperInt64Attribute{
												Common: &schemaR.Int64Attribute{
													MarkdownDescription: "HTTP status code to return.",
												},
												Resource: &schemaR.Int64Attribute{
													Required: true,
													Validators: []validator.Int64{
														int64validator.OneOf(edgeloadbalancer.PoliciesHTTPActionResponseStatusCodes...),
													},
												},
												DataSource: &schemaD.Int64Attribute{
													Computed: true,
												},
											},
											"content": superschema.SuperStringAttribute{
												Common: &schemaR.StringAttribute{
													MarkdownDescription: "Content of the response must be a base64 encoded string.",
												},
												Resource: &schemaR.StringAttribute{
													MarkdownDescription: "Example: result of `echo -n 'Hello World' | base64` will be `SGVsbG8gV29ybGQ=`.",
													Required:            true,
													Validators: []validator.String{
														stringvalidator.LengthAtMost(10240),
														fstringvalidator.Formats(
															[]fstringvalidator.FormatsValidatorType{
																fstringvalidator.FormatsIsBase64,
															},
															false,
														),
													},
												},
												DataSource: &schemaD.StringAttribute{
													Computed: true,
												},
											},
											"content_type": superschema.SuperStringAttribute{
												Common: &schemaR.StringAttribute{
													MarkdownDescription: "Mime type of content.",
												},
												Resource: &schemaR.StringAttribute{
													Required: true,
													Validators: []validator.String{
														stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPActionContentTypesString...),
													},
												},
												DataSource: &schemaD.StringAttribute{
													Computed: true,
												},
											},
										},
									}, // End local_response
									"close_connection": superschema.SuperBoolAttribute{
										Common: &schemaR.BoolAttribute{
											MarkdownDescription: "Close connection when the rate limit is reached",
										},
										Resource: &schemaR.BoolAttribute{
											Optional: true,
											Validators: []validator.Bool{
												fboolvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("redirect")),
												fboolvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("local_response")),
											},
										},
										DataSource: &schemaD.BoolAttribute{
											Computed: true,
										},
									},
								},
							}, // End rate_limit
						},
					}, // End actions
				},
			}, // End policies
		},
	}
}

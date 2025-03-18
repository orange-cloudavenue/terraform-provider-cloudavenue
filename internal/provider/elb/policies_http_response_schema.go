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
	stringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
)

func policiesHTTPResponseSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_elb_policies_http_response` resource allows you to manage HTTP response policies. HTTP response rules can be used to to evaluate and modify the response and response attributes that the application returns.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_elb_policies_http_response` data source allows you to retrieve information about an existing HTTP response policies.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the policies http response.",
				},
			},
			"virtual_service_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the virtual service to which the policies http response belongs.",
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"policies": superschema.SuperListNestedAttributeOf[PoliciesHTTPResponseModelPolicies]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "HTTP response policies.",
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
							MarkdownDescription: "The name of the policy.",
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
							MarkdownDescription: "Whether the policy is active or not.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(true),
						},
					},
					"logging": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Enable logging for this policy.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
						},
					},
					"criteria": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPResponseMatchCriteria]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Match criteria for the HTTP response.",
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
									MarkdownDescription: "Protocol to match.",
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
											MarkdownDescription: "Ports to match.",
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
							"location": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPLocationMatch]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Match the rule based on location rules.",
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
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPLocationMatchCriteriaString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"values": superschema.SuperSetAttributeOf[string]{
										Common: &schemaR.SetAttribute{
											MarkdownDescription: "A set of locations to match given criteria.",
										},
										Resource: &schemaR.SetAttribute{
											Required: true,
										},
										DataSource: &schemaD.SetAttribute{
											Computed: true,
										},
									},
								},
							}, // End location
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
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPRequestHeaderMatchCriteriaString...),
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
											// 			for _, v := range edgeloadbalancer.PoliciesHTTPRequestHeaderMatchCriteriaString {
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
							}, // End request_headers
							"response_headers": superschema.SuperSetNestedAttributeOf[PoliciesHTTPHeaderMatch]{
								Common: &schemaR.SetNestedAttribute{
									MarkdownDescription: "Match the rule based on response headers rules.",
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
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPRequestHeaderMatchCriteriaString...),
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
											// 			for _, v := range edgeloadbalancer.PoliciesHTTPRequestHeaderMatchCriteriaString {
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
							}, // End request_headers
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
							}, // End query
							"status_code": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPStatusCodeMatch]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Match the rule based on response HTTP status code.",
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
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPStatusCodeMatchCriteriaString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"codes": superschema.SuperSetAttributeOf[string]{
										Common: &schemaR.SetAttribute{
											MarkdownDescription: "HTTP Status codes or range to match. (Example: `200` or `301-304`) Warning: all ports must have valid HTTP return codes. `200-299` are invalid range because they are not a valid HTTP status code.",
										},
										Resource: &schemaR.SetAttribute{
											Required: true,
											// TODO add http status code
										},
										DataSource: &schemaD.SetAttribute{
											Computed: true,
										},
									},
								},
							}, // End status_code
						},
					}, // End criteria
					"actions": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPResponseActions]{
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
							"location_rewrite": superschema.SuperSingleNestedAttributeOf[PoliciesHTTPActionLocationRewrite]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Redirects the request to different location.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
									// Validators: []validator.Object{
									// 	objectvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("modify_headers")),
									// },
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
								},
							}, // End redirect
							"modify_headers": superschema.SuperSetNestedAttributeOf[PoliciesHTTPActionHeaderRewrite]{
								Common: &schemaR.SetNestedAttribute{
									MarkdownDescription: "Modify HTTP request headers.",
								},
								Resource: &schemaR.SetNestedAttribute{
									Optional: true,
									// Validators: []validator.Set{
									// 	setvalidator.SizeAtMost(10),
									// 	fsetvalidator.NullIfAttributeIsSet(path.MatchRelative().AtParent().AtName("redirect")),
									// },
								},
								DataSource: &schemaD.SetNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"action": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Action to perform on the header.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf(edgeloadbalancer.PoliciesHTTPActionHeaderRewriteActionsString...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"name": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Name of the HTTP header to modify.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.LengthAtMost(128),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"value": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Value of the HTTP header to modify.",
										},
										Resource: &schemaR.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												// fstringvalidator.RequireIfAttributeIsOneOf(
												// 	path.MatchRelative().AtParent().AtName("action"),
												// 	[]attr.Value{
												// 		types.StringValue(string(edgeloadbalancer.PoliciesHTTPActionHeaderRewriteActionADD)),
												// 		types.StringValue(string(edgeloadbalancer.PoliciesHTTPActionHeaderRewriteActionREPLACE)),
												// 	},
												// ),
												// fstringvalidator.NullIfAttributeIsOneOf(
												// 	path.MatchRelative().AtParent().AtName("action"),
												// 	[]attr.Value{
												// 		types.StringValue(string(edgeloadbalancer.PoliciesHTTPActionHeaderRewriteActionREMOVE)),
												// 	},
												// ),
												stringvalidator.LengthAtMost(128),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
								},
							}, // End modify_headers
						},
					}, // End actions
				},
			}, // End policies
		},
	}
}

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

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	policiesHTTPRequestPrivateModel struct {
		EdgeGatewayID string `json:"edgeGatewayId"`
	}

	PoliciesHTTPRequestModel struct {
		ID               supertypes.StringValue                                               `tfsdk:"id"`
		VirtualServiceID supertypes.StringValue                                               `tfsdk:"virtual_service_id"`
		Policies         supertypes.ListNestedObjectValueOf[PoliciesHTTPRequestModelPolicies] `tfsdk:"policies"`
	}

	PoliciesHTTPRequestModelPolicies struct {
		Name     supertypes.StringValue                                                 `tfsdk:"name"`
		Active   supertypes.BoolValue                                                   `tfsdk:"active"`
		Logging  supertypes.BoolValue                                                   `tfsdk:"logging"`
		Criteria supertypes.SingleNestedObjectValueOf[PoliciesHTTPRequestMatchCriteria] `tfsdk:"criteria"`
		Actions  supertypes.SingleNestedObjectValueOf[PoliciesHTTPRequestActions]       `tfsdk:"actions"`
	}

	PoliciesHTTPRequestMatchCriteria struct {
		Protocol       supertypes.StringValue                                             `tfsdk:"protocol"`
		ClientIP       supertypes.SingleNestedObjectValueOf[PoliciesHTTPClientIPMatch]    `tfsdk:"client_ip"`
		ServicePorts   supertypes.SingleNestedObjectValueOf[PoliciesHTTPServicePortMatch] `tfsdk:"service_ports"`
		HTTPMethods    supertypes.SingleNestedObjectValueOf[PoliciesHTTPMethodMatch]      `tfsdk:"http_methods"`
		Path           supertypes.SingleNestedObjectValueOf[PoliciesHTTPPathMatch]        `tfsdk:"path"`
		Cookie         supertypes.SingleNestedObjectValueOf[PoliciesHTTPCookieMatch]      `tfsdk:"cookie"`
		RequestHeaders supertypes.SetNestedObjectValueOf[PoliciesHTTPHeaderMatch]         `tfsdk:"request_headers"`
		Query          supertypes.SetValueOf[string]                                      `tfsdk:"query"`
	}

	PoliciesHTTPRequestActions struct {
		Redirect      supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionRedirect]   `tfsdk:"redirect"`
		RewriteURL    supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionURLRewrite] `tfsdk:"rewrite_url"`
		ModifyHeaders supertypes.SetNestedObjectValueOf[PoliciesHTTPActionHeaderRewrite] `tfsdk:"modify_headers"`
	}

	// * Match.
	PoliciesHTTPClientIPMatch struct {
		Criteria    supertypes.StringValue        `tfsdk:"criteria"`
		IPAddresses supertypes.SetValueOf[string] `tfsdk:"ip_addresses"`
	}
	PoliciesHTTPServicePortMatch struct {
		Criteria supertypes.StringValue       `tfsdk:"criteria"`
		Ports    supertypes.SetValueOf[int64] `tfsdk:"ports"`
	}
	PoliciesHTTPMethodMatch struct {
		Criteria supertypes.StringValue        `tfsdk:"criteria"`
		Methods  supertypes.SetValueOf[string] `tfsdk:"methods"`
	}
	PoliciesHTTPPathMatch struct {
		Criteria supertypes.StringValue        `tfsdk:"criteria"`
		Paths    supertypes.SetValueOf[string] `tfsdk:"paths"`
	}
	PoliciesHTTPHeaderMatch struct {
		Criteria supertypes.StringValue        `tfsdk:"criteria"`
		Name     supertypes.StringValue        `tfsdk:"name"`
		Values   supertypes.SetValueOf[string] `tfsdk:"values"`
	}
	PoliciesHTTPCookieMatch struct {
		Criteria supertypes.StringValue `tfsdk:"criteria"`
		Name     supertypes.StringValue `tfsdk:"name"`
		Value    supertypes.StringValue `tfsdk:"value"`
	}

	// * Action.
	PoliciesHTTPActionRedirect struct {
		Host       supertypes.StringValue `tfsdk:"host"`
		KeepQuery  supertypes.BoolValue   `tfsdk:"keep_query"`
		Path       supertypes.StringValue `tfsdk:"path"`
		Port       supertypes.Int64Value  `tfsdk:"port"`
		Protocol   supertypes.StringValue `tfsdk:"protocol"`
		StatusCode supertypes.Int64Value  `tfsdk:"status_code"`
	}
	PoliciesHTTPActionURLRewrite struct {
		Host      supertypes.StringValue `tfsdk:"host"`
		Path      supertypes.StringValue `tfsdk:"path"`
		Query     supertypes.StringValue `tfsdk:"query"`
		KeepQuery supertypes.BoolValue   `tfsdk:"keep_query"`
	}
	PoliciesHTTPActionHeaderRewrite struct {
		Action supertypes.StringValue `tfsdk:"action"`
		Name   supertypes.StringValue `tfsdk:"name"`
		Value  supertypes.StringValue `tfsdk:"value"`
	}
)

func (rm *PoliciesHTTPRequestModel) Copy() *PoliciesHTTPRequestModel {
	x := &PoliciesHTTPRequestModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *GoPoliciesHTTPRequestModel) ToSDKPoliciesHTTPRequestModel(ctx context.Context) (*edgeloadbalancer.PoliciesHTTPRequestModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &edgeloadbalancer.PoliciesHTTPRequestModel{
		VirtualServiceID: rm.VirtualServiceID,
		Policies:         make([]*edgeloadbalancer.PoliciesHTTPRequestModelPolicies, 0),
	}

	for _, policy := range rm.Policies {
		model.Policies = append(model.Policies, &edgeloadbalancer.PoliciesHTTPRequestModelPolicies{
			Name:    policy.Name,
			Active:  policy.Active,
			Logging: policy.Logging,
			MatchCriteria: edgeloadbalancer.PoliciesHTTPRequestMatchCriteria{
				Protocol: policy.Criteria.Protocol,
				ClientIPMatch: func() *edgeloadbalancer.PoliciesHTTPClientIPMatch {
					if policy.Criteria.ClientIP == nil {
						return nil
					}

					return &edgeloadbalancer.PoliciesHTTPClientIPMatch{
						Criteria:  policy.Criteria.ClientIP.Criteria,
						Addresses: policy.Criteria.ClientIP.IPAddresses,
					}
				}(),
				ServicePortMatch: func() *edgeloadbalancer.PoliciesHTTPServicePortMatch {
					if policy.Criteria.ServicePorts == nil {
						return nil
					}

					return &edgeloadbalancer.PoliciesHTTPServicePortMatch{
						Criteria: policy.Criteria.ServicePorts.Criteria,
						Ports: func(x []int64) (p []int) { // Convert int64 to int
							for _, v := range x {
								p = append(p, int(v))
							}
							return p
						}(policy.Criteria.ServicePorts.Ports),
					}
				}(),
				MethodMatch: func() *edgeloadbalancer.PoliciesHTTPMethodMatch {
					if policy.Criteria.HTTPMethods == nil {
						return nil
					}

					return &edgeloadbalancer.PoliciesHTTPMethodMatch{
						Criteria: policy.Criteria.HTTPMethods.Criteria,
						Methods:  policy.Criteria.HTTPMethods.Methods,
					}
				}(),
				PathMatch: func() *edgeloadbalancer.PoliciesHTTPPathMatch {
					if policy.Criteria.Path == nil {
						return nil
					}

					return &edgeloadbalancer.PoliciesHTTPPathMatch{
						Criteria:     policy.Criteria.Path.Criteria,
						MatchStrings: policy.Criteria.Path.Paths,
					}
				}(),
				CookieMatch: func() *edgeloadbalancer.PoliciesHTTPCookieMatch {
					if policy.Criteria.Cookie == nil {
						return nil
					}

					return &edgeloadbalancer.PoliciesHTTPCookieMatch{
						Criteria: policy.Criteria.Cookie.Criteria,
						Name:     policy.Criteria.Cookie.Name,
						Value:    policy.Criteria.Cookie.Value,
					}
				}(),
				HeaderMatch: func() edgeloadbalancer.PoliciesHTTPHeadersMatch {
					var headers edgeloadbalancer.PoliciesHTTPHeadersMatch
					for _, header := range policy.Criteria.RequestHeaders {
						headers = append(headers, edgeloadbalancer.PoliciesHTTPHeaderMatch{
							Criteria: header.Criteria,
							Name:     header.Name,
							Values:   header.Values,
						})
					}
					return headers
				}(),
				QueryMatch: policy.Criteria.Query,
			},
			HeaderRewriteActions: func() edgeloadbalancer.PoliciesHTTPActionHeadersRewrite {
				var headers edgeloadbalancer.PoliciesHTTPActionHeadersRewrite
				for _, header := range policy.Actions.ModifyHeaders {
					headers = append(headers, &edgeloadbalancer.PoliciesHTTPActionHeaderRewrite{
						Action: header.Action,
						Name:   header.Name,
						Value:  header.Value,
					})
				}
				return headers
			}(),
			URLRewriteAction: func() *edgeloadbalancer.PoliciesHTTPActionURLRewrite {
				if policy.Actions.RewriteURL == nil {
					return nil
				}

				return &edgeloadbalancer.PoliciesHTTPActionURLRewrite{
					HostHeader: policy.Actions.RewriteURL.Host,
					Path:       policy.Actions.RewriteURL.Path,
					Query:      policy.Actions.RewriteURL.Query,
					KeepQuery:  policy.Actions.RewriteURL.KeepQuery,
				}
			}(),
			RedirectAction: func() *edgeloadbalancer.PoliciesHTTPActionRedirect {
				if policy.Actions.Redirect == nil {
					return nil
				}

				return &edgeloadbalancer.PoliciesHTTPActionRedirect{
					Host:       policy.Actions.Redirect.Host,
					KeepQuery:  policy.Actions.Redirect.KeepQuery,
					Path:       policy.Actions.Redirect.Path,
					Port:       policy.Actions.Redirect.Port,
					Protocol:   policy.Actions.Redirect.Protocol,
					StatusCode: policy.Actions.Redirect.StatusCode,
				}
			}(),
		})
	}

	return model, diags
}

func isNotNil[T any](x *T) *T {
	if x == nil {
		return nil
	}

	return x
}

// ToSDKPoliciesHttpRequestGroupModel converts the model to the SDK model.
// func (rm *PoliciesHTTPRequestModel) ToSDKPoliciesHTTPRequestModel(ctx context.Context) (*edgeloadbalancer.PoliciesHTTPRequestModel, diag.Diagnostics) {
// 	var diags diag.Diagnostics

// 	model := &edgeloadbalancer.PoliciesHTTPRequestModel{
// 		VirtualServiceID: rm.VirtualServiceID.Get(),
// 		Policies:         make([]*edgeloadbalancer.PoliciesHTTPRequestModelPolicies, 0),
// 	}

// 	if rm.Policies.IsKnown() {
// 		x, d := rm.Policies.Get(ctx)
// 		diags.Append(d...)
// 		if diags.HasError() {
// 			return nil, diags
// 		}

// 		for _, v := range x {
// 			if v == nil {
// 				continue
// 			}

// 			// Criteria
// 			criteria, d := v.Criteria.Get(ctx)
// 			diags.Append(d...)
// 			if diags.HasError() {
// 				return nil, diags
// 			}

// 			clientIP, d := criteria.ClientIP.Get(ctx)
// 			diags.Append(d...)
// 			servicePorts, d := criteria.ServicePorts.Get(ctx)
// 			diags.Append(d...)
// 			httpMethods, d := criteria.HTTPMethods.Get(ctx)
// 			diags.Append(d...)
// 			path, d := criteria.Path.Get(ctx)
// 			diags.Append(d...)
// 			cookie, d := criteria.Cookie.Get(ctx)
// 			diags.Append(d...)
// 			requestHeaders, d := criteria.RequestHeaders.Get(ctx)
// 			diags.Append(d...)
// 			query, d := criteria.Query.Get(ctx)
// 			diags.Append(d...)
// 			if diags.HasError() {
// 				return nil, diags
// 			}

// 			computeRequestHeaders := func() edgeloadbalancer.PoliciesHTTPHeadersMatch {
// 				var headers edgeloadbalancer.PoliciesHTTPHeadersMatch
// 				for _, header := range requestHeaders {
// 					headers = append(headers, edgeloadbalancer.PoliciesHTTPHeaderMatch{
// 						Criteria: header.Criteria.Get(),
// 						Name:     header.Name.Get(),
// 						Values:   mustSet(ctx, header.Values, diags),
// 					})
// 				}
// 				return headers
// 			}

// 			// Actions

// 			actions, d := v.Actions.Get(ctx)
// 			diags.Append(d...)
// 			if diags.HasError() {
// 				return nil, diags
// 			}

// 			var (
// 				modifyHeaders []*PoliciesHTTPActionHeaderRewrite
// 				rewriteURL    *PoliciesHTTPActionURLRewrite
// 				redirect      *PoliciesHTTPActionRedirect
// 			)

// 			if actions.ModifyHeaders.IsKnown() {
// 				modifyHeaders, d = actions.ModifyHeaders.Get(ctx)
// 				diags.Append(d...)
// 			}
// 			if actions.RewriteURL.IsKnown() {
// 				rewriteURL, d = actions.RewriteURL.Get(ctx)
// 				diags.Append(d...)
// 			}
// 			if actions.Redirect.IsKnown() {
// 				redirect, d = actions.Redirect.Get(ctx)
// 				diags.Append(d...)
// 			}
// 			if diags.HasError() {
// 				return nil, diags
// 			}

// 			model.Policies = append(model.Policies, &edgeloadbalancer.PoliciesHTTPRequestModelPolicies{
// 				Name:    v.Name.Get(),
// 				Active:  v.Active.Get(),
// 				Logging: v.Logging.Get(),
// 				MatchCriteria: edgeloadbalancer.PoliciesHTTPRequestMatchCriteria{
// 					Protocol:      criteria.Protocol.Get(),
// 					ClientIPMatch: isKnown(clientIP, &edgeloadbalancer.PoliciesHTTPClientIPMatch{Criteria: clientIP.Criteria.Get(), Addresses: mustSet(ctx, clientIP.IPAddresses, diags)}),
// 					ServicePortMatch: isKnown(servicePorts, &edgeloadbalancer.PoliciesHTTPServicePortMatch{Criteria: servicePorts.Criteria.Get(), Ports: func(x []int64) (p []int) {
// 						// Convert int64 to int
// 						for _, v := range x {
// 							p = append(p, int(v))
// 						}

// 						return p
// 					}(mustSet(ctx, servicePorts.Ports, diags))}),
// 					MethodMatch: isKnown(httpMethods, &edgeloadbalancer.PoliciesHTTPMethodMatch{Criteria: httpMethods.Criteria.Get(), Methods: mustSet(ctx, httpMethods.Methods, diags)}),
// 					PathMatch:   isKnown(path, &edgeloadbalancer.PoliciesHTTPPathMatch{Criteria: path.Criteria.Get(), MatchStrings: mustSet(ctx, path.Paths, diags)}),
// 					CookieMatch: isKnown(cookie, &edgeloadbalancer.PoliciesHTTPCookieMatch{Criteria: cookie.Criteria.Get(), Name: cookie.Name.Get(), Value: cookie.Value.Get()}),
// 					HeaderMatch: computeRequestHeaders(),
// 					QueryMatch:  query,
// 				},
// 				RedirectAction: func() *edgeloadbalancer.PoliciesHTTPActionRedirect {
// 					if redirect == nil {
// 						return nil
// 					}

// 					return &edgeloadbalancer.PoliciesHTTPActionRedirect{
// 						Host:       redirect.Host.Get(),
// 						KeepQuery:  redirect.KeepQuery.Get(),
// 						Path:       redirect.Path.Get(),
// 						Port:       redirect.Port.GetIntPtr(),
// 						Protocol:   redirect.Protocol.Get(),
// 						StatusCode: redirect.StatusCode.GetInt(),
// 					}
// 				}(),
// 				URLRewriteAction: func() *edgeloadbalancer.PoliciesHTTPActionURLRewrite {
// 					if rewriteURL == nil {
// 						return nil
// 					}

// 					return &edgeloadbalancer.PoliciesHTTPActionURLRewrite{
// 						HostHeader: rewriteURL.Host.Get(),
// 						Path:       rewriteURL.Path.Get(),
// 						Query:      rewriteURL.Query.Get(),
// 						KeepQuery:  rewriteURL.KeepQuery.Get(),
// 					}
// 				}(),
// 				HeaderRewriteActions: func() edgeloadbalancer.PoliciesHTTPActionHeadersRewrite {
// 					var headers edgeloadbalancer.PoliciesHTTPActionHeadersRewrite
// 					for _, header := range modifyHeaders {
// 						headers = append(headers, &edgeloadbalancer.PoliciesHTTPActionHeaderRewrite{
// 							Action: header.Action.Get(),
// 							Name:   header.Name.Get(),
// 							Value:  header.Value.Get(),
// 						})
// 					}
// 					return headers
// 				}(),
// 			})
// 			if diags.HasError() {
// 				return nil, diags
// 			}
// 		}
// 	}

// 	return model, diags
// }

// func mustSet[T any](ctx context.Context, getter supertypes.SetValueOf[T], d diag.Diagnostics) []T {
// 	values, dd := getter.Get(ctx)
// 	d.Append(dd...)

// 	return values
// }

// func isKnown[T any](entry interface{}, x *T) *T {
// 	if entry == nil {
// 		return nil
// 	}

// 	return x
// }

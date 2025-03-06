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

func (rm *PoliciesHTTPRequestModel) ToSDKPoliciesHTTPRequestModel(ctx context.Context) (*edgeloadbalancer.PoliciesHTTPRequestModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &edgeloadbalancer.PoliciesHTTPRequestModel{
		VirtualServiceID: rm.VirtualServiceID.Get(),
		Policies:         make([]*edgeloadbalancer.PoliciesHTTPRequestModelPolicies, 0),
	}

	for _, policy := range rm.Policies.DiagsGet(ctx, diags) {
		actions := func() *PoliciesHTTPRequestActions {
			if !policy.Actions.IsKnown() {
				return nil
			}

			return policy.Actions.DiagsGet(ctx, diags)
		}()

		model.Policies = append(model.Policies, &edgeloadbalancer.PoliciesHTTPRequestModelPolicies{
			Name:    policy.Name.Get(),
			Active:  policy.Active.Get(),
			Logging: policy.Logging.Get(),
			MatchCriteria: func() edgeloadbalancer.PoliciesHTTPRequestMatchCriteria {
				if !policy.Criteria.IsKnown() {
					return edgeloadbalancer.PoliciesHTTPRequestMatchCriteria{}
				}

				criteria := policy.Criteria.DiagsGet(ctx, diags)
				return edgeloadbalancer.PoliciesHTTPRequestMatchCriteria{
					Protocol: criteria.Protocol.Get(),
					ClientIPMatch: func() *edgeloadbalancer.PoliciesHTTPClientIPMatch {
						if !criteria.ClientIP.IsKnown() {
							return nil
						}

						cIP := criteria.ClientIP.DiagsGet(ctx, diags)

						return &edgeloadbalancer.PoliciesHTTPClientIPMatch{
							Criteria:  cIP.Criteria.Get(),
							Addresses: cIP.IPAddresses.DiagsGet(ctx, diags),
						}
					}(),
					ServicePortMatch: func() *edgeloadbalancer.PoliciesHTTPServicePortMatch {
						if !criteria.ServicePorts.IsKnown() {
							return nil
						}

						sP := criteria.ServicePorts.DiagsGet(ctx, diags)
						return &edgeloadbalancer.PoliciesHTTPServicePortMatch{
							Criteria: sP.Criteria.Get(),
							Ports: func(x []int64) (p []int) { // Convert int64 to int
								for _, v := range x {
									p = append(p, int(v))
								}
								return p
							}(sP.Ports.DiagsGet(ctx, diags)),
						}
					}(),
					MethodMatch: func() *edgeloadbalancer.PoliciesHTTPMethodMatch {
						if !criteria.HTTPMethods.IsKnown() {
							return nil
						}

						m := criteria.HTTPMethods.DiagsGet(ctx, diags)
						return &edgeloadbalancer.PoliciesHTTPMethodMatch{
							Criteria: m.Criteria.Get(),
							Methods:  m.Methods.DiagsGet(ctx, diags),
						}
					}(),
					PathMatch: func() *edgeloadbalancer.PoliciesHTTPPathMatch {
						if !criteria.Path.IsKnown() {
							return nil
						}

						p := criteria.Path.DiagsGet(ctx, diags)

						return &edgeloadbalancer.PoliciesHTTPPathMatch{
							Criteria:     p.Criteria.Get(),
							MatchStrings: p.Paths.DiagsGet(ctx, diags),
						}
					}(),
					CookieMatch: func() *edgeloadbalancer.PoliciesHTTPCookieMatch {
						if !criteria.Cookie.IsKnown() {
							return nil
						}

						cM := criteria.Cookie.DiagsGet(ctx, diags)

						return &edgeloadbalancer.PoliciesHTTPCookieMatch{
							Criteria: cM.Criteria.Get(),
							Name:     cM.Name.Get(),
							Value:    cM.Value.Get(),
						}
					}(),
					HeaderMatch: func() edgeloadbalancer.PoliciesHTTPHeadersMatch {
						if !criteria.RequestHeaders.IsKnown() {
							return nil
						}

						h := criteria.RequestHeaders.DiagsGet(ctx, diags)

						var headers edgeloadbalancer.PoliciesHTTPHeadersMatch
						for _, header := range h {
							headers = append(headers, edgeloadbalancer.PoliciesHTTPHeaderMatch{
								Criteria: header.Criteria.Get(),
								Name:     header.Name.Get(),
								Values:   header.Values.DiagsGet(ctx, diags),
							})
						}
						return headers
					}(),
					QueryMatch: criteria.Query.DiagsGet(ctx, diags),
				}
			}(),
			HeaderRewriteActions: func() edgeloadbalancer.PoliciesHTTPActionHeadersRewrite {
				if actions == nil || !actions.ModifyHeaders.IsKnown() {
					return nil
				}

				hra := actions.ModifyHeaders.DiagsGet(ctx, diags)

				var headers edgeloadbalancer.PoliciesHTTPActionHeadersRewrite
				for _, header := range hra {
					headers = append(headers, &edgeloadbalancer.PoliciesHTTPActionHeaderRewrite{
						Action: header.Action.Get(),
						Name:   header.Name.Get(),
						Value:  header.Value.Get(),
					})
				}
				return headers
			}(),
			URLRewriteAction: func() *edgeloadbalancer.PoliciesHTTPActionURLRewrite {
				if actions == nil || !actions.RewriteURL.IsKnown() {
					return nil
				}

				rURL := actions.RewriteURL.DiagsGet(ctx, diags)

				return &edgeloadbalancer.PoliciesHTTPActionURLRewrite{
					HostHeader: rURL.Host.Get(),
					Path:       rURL.Path.Get(),
					Query:      rURL.Query.Get(),
					KeepQuery:  rURL.KeepQuery.Get(),
				}
			}(),
			RedirectAction: func() *edgeloadbalancer.PoliciesHTTPActionRedirect {
				if actions == nil || !actions.Redirect.IsKnown() {
					return nil
				}

				rA := actions.Redirect.DiagsGet(ctx, diags)

				return &edgeloadbalancer.PoliciesHTTPActionRedirect{
					Host:       rA.Host.Get(),
					KeepQuery:  rA.KeepQuery.Get(),
					Path:       rA.Path.Get(),
					Port:       rA.Port.GetIntPtr(),
					Protocol:   rA.Protocol.Get(),
					StatusCode: rA.StatusCode.GetInt(),
				}
			}(),
		})
	}

	return model, diags
}

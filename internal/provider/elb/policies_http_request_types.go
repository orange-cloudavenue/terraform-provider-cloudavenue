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

package elb

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
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
		Policies:         make([]*edgeloadbalancer.PoliciesHTTPRequestModelPolicy, 0),
	}

	for _, policy := range rm.Policies.DiagsGet(ctx, diags) {
		actions := func() *PoliciesHTTPRequestActions {
			if !policy.Actions.IsKnown() {
				return nil
			}

			return policy.Actions.DiagsGet(ctx, diags)
		}()

		model.Policies = append(model.Policies, &edgeloadbalancer.PoliciesHTTPRequestModelPolicy{
			Name:    policy.Name.Get(),
			Active:  policy.Active.Get(),
			Logging: policy.Logging.Get(),
			MatchCriteria: func() edgeloadbalancer.PoliciesHTTPRequestMatchCriteria {
				if !policy.Criteria.IsKnown() {
					return edgeloadbalancer.PoliciesHTTPRequestMatchCriteria{}
				}

				criteria := policy.Criteria.DiagsGet(ctx, diags)
				return edgeloadbalancer.PoliciesHTTPRequestMatchCriteria{
					Protocol:         criteria.Protocol.Get(),
					ClientIPMatch:    policiesHTTPClientIPMatchToSDK(ctx, diags, criteria.ClientIP),
					ServicePortMatch: policiesHTTPServicePortMatchToSDK(ctx, diags, criteria.ServicePorts),
					MethodMatch:      policiesHTTPMethodMatchToSDK(ctx, diags, criteria.HTTPMethods),
					PathMatch:        policiesHTTPPathMatchToSDK(ctx, diags, criteria.Path),
					CookieMatch:      policiesHTTPCookieMatchToSDK(ctx, diags, criteria.Cookie),
					HeaderMatch:      policiesHTTPHeadersMatchToSDK(ctx, diags, criteria.RequestHeaders),
					QueryMatch:       criteria.Query.DiagsGet(ctx, diags),
				}
			}(),
			HeaderRewriteActions: policiesHTTPActionHeadersRewriteToSDK(ctx, diags, actions.ModifyHeaders),
			URLRewriteAction:     policiesHTTPActionURLRewriteToSDK(ctx, diags, actions.RewriteURL),
			RedirectAction:       policiesHTTPActionRedirectToSDK(ctx, diags, actions.Redirect),
		})
	}

	return model, diags
}

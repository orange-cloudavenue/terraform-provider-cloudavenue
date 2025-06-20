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
	PoliciesHTTPSecurityModel struct {
		ID               supertypes.StringValue                                                `tfsdk:"id"`
		VirtualServiceID supertypes.StringValue                                                `tfsdk:"virtual_service_id"`
		Policies         supertypes.ListNestedObjectValueOf[PoliciesHTTPSecurityModelPolicies] `tfsdk:"policies"`
	}

	PoliciesHTTPSecurityModelPolicies struct {
		Name     supertypes.StringValue                                                  `tfsdk:"name"`
		Active   supertypes.BoolValue                                                    `tfsdk:"active"`
		Logging  supertypes.BoolValue                                                    `tfsdk:"logging"`
		Criteria supertypes.SingleNestedObjectValueOf[PoliciesHTTPSecurityMatchCriteria] `tfsdk:"criteria"`
		Actions  supertypes.SingleNestedObjectValueOf[PoliciesHTTPSecurityActions]       `tfsdk:"actions"`
	}

	PoliciesHTTPSecurityMatchCriteria struct {
		Protocol       supertypes.StringValue                                             `tfsdk:"protocol"`
		ClientIP       supertypes.SingleNestedObjectValueOf[PoliciesHTTPClientIPMatch]    `tfsdk:"client_ip"`
		ServicePorts   supertypes.SingleNestedObjectValueOf[PoliciesHTTPServicePortMatch] `tfsdk:"service_ports"`
		HTTPMethods    supertypes.SingleNestedObjectValueOf[PoliciesHTTPMethodMatch]      `tfsdk:"http_methods"`
		Path           supertypes.SingleNestedObjectValueOf[PoliciesHTTPPathMatch]        `tfsdk:"path"`
		Cookie         supertypes.SingleNestedObjectValueOf[PoliciesHTTPCookieMatch]      `tfsdk:"cookie"`
		RequestHeaders supertypes.SetNestedObjectValueOf[PoliciesHTTPHeaderMatch]         `tfsdk:"request_headers"`
		Query          supertypes.SetValueOf[string]                                      `tfsdk:"query"`
	}

	PoliciesHTTPSecurityActions struct {
		Connection      supertypes.StringValue                                               `tfsdk:"connection"`
		RedirectToHTTPS supertypes.Int64Value                                                `tfsdk:"redirect_to_https"`
		SendResponse    supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionSendResponse] `tfsdk:"send_response"`
		RateLimit       supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionRateLimit]    `tfsdk:"rate_limit"`
	}
)

func (rm *PoliciesHTTPSecurityModel) Copy() *PoliciesHTTPSecurityModel {
	x := &PoliciesHTTPSecurityModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *PoliciesHTTPSecurityModel) ToSDKPoliciesHTTPSecurityModel(ctx context.Context) (*edgeloadbalancer.PoliciesHTTPSecurityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &edgeloadbalancer.PoliciesHTTPSecurityModel{
		VirtualServiceID: rm.VirtualServiceID.Get(),
		Policies:         make([]*edgeloadbalancer.PoliciesHTTPSecurityModelPolicy, 0),
	}

	for _, policy := range rm.Policies.DiagsGet(ctx, diags) {
		actions := func() *PoliciesHTTPSecurityActions {
			if !policy.Actions.IsKnown() {
				return nil
			}

			return policy.Actions.DiagsGet(ctx, diags)
		}()

		model.Policies = append(model.Policies, &edgeloadbalancer.PoliciesHTTPSecurityModelPolicy{
			Name:    policy.Name.Get(),
			Active:  policy.Active.Get(),
			Logging: policy.Logging.Get(),
			MatchCriteria: func() edgeloadbalancer.PoliciesHTTPSecurityMatchCriteria {
				if !policy.Criteria.IsKnown() {
					return edgeloadbalancer.PoliciesHTTPSecurityMatchCriteria{}
				}

				criteria := policy.Criteria.DiagsGet(ctx, diags)
				return edgeloadbalancer.PoliciesHTTPSecurityMatchCriteria{
					Protocol:         edgeloadbalancer.PoliciesHTTPProtocol(criteria.Protocol.Get()),
					ClientIPMatch:    policiesHTTPClientIPMatchToSDK(ctx, diags, criteria.ClientIP),
					ServicePortMatch: policiesHTTPServicePortMatchToSDK(ctx, diags, criteria.ServicePorts),
					MethodMatch:      policiesHTTPMethodMatchToSDK(ctx, diags, criteria.HTTPMethods),
					PathMatch:        policiesHTTPPathMatchToSDK(ctx, diags, criteria.Path),
					CookieMatch:      policiesHTTPCookieMatchToSDK(ctx, diags, criteria.Cookie),
					HeaderMatch:      policiesHTTPHeadersMatchToSDK(ctx, diags, criteria.RequestHeaders),
					QueryMatch:       criteria.Query.DiagsGet(ctx, diags),
				}
			}(),
			ConnectionAction: edgeloadbalancer.PoliciesHTTPConnectionAction(actions.Connection.Get()),

			RedirectToHTTPSAction: actions.RedirectToHTTPS.GetIntPtr(),

			SendResponseAction: policiesHTTPActionSendResponseToSDK(ctx, diags, actions.SendResponse),
			RateLimitAction:    policiesHTTPActionRateLimitToSDK(ctx, diags, actions.RateLimit),
		})
	}

	return model, diags
}

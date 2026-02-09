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
	PoliciesHTTPResponseModel struct {
		ID               supertypes.StringValue                                                `tfsdk:"id"`
		VirtualServiceID supertypes.StringValue                                                `tfsdk:"virtual_service_id"`
		Policies         supertypes.ListNestedObjectValueOf[PoliciesHTTPResponseModelPolicies] `tfsdk:"policies"`
	}

	PoliciesHTTPResponseModelPolicies struct {
		Name     supertypes.StringValue                                                  `tfsdk:"name"`
		Active   supertypes.BoolValue                                                    `tfsdk:"active"`
		Logging  supertypes.BoolValue                                                    `tfsdk:"logging"`
		Criteria supertypes.SingleNestedObjectValueOf[PoliciesHTTPResponseMatchCriteria] `tfsdk:"criteria"`
		Actions  supertypes.SingleNestedObjectValueOf[PoliciesHTTPResponseActions]       `tfsdk:"actions"`
	}

	PoliciesHTTPResponseMatchCriteria struct {
		Protocol        supertypes.StringValue                                             `tfsdk:"protocol"`
		ClientIP        supertypes.SingleNestedObjectValueOf[PoliciesHTTPClientIPMatch]    `tfsdk:"client_ip"`
		ServicePorts    supertypes.SingleNestedObjectValueOf[PoliciesHTTPServicePortMatch] `tfsdk:"service_ports"`
		HTTPMethods     supertypes.SingleNestedObjectValueOf[PoliciesHTTPMethodMatch]      `tfsdk:"http_methods"`
		Path            supertypes.SingleNestedObjectValueOf[PoliciesHTTPPathMatch]        `tfsdk:"path"`
		Cookie          supertypes.SingleNestedObjectValueOf[PoliciesHTTPCookieMatch]      `tfsdk:"cookie"`
		Location        supertypes.SingleNestedObjectValueOf[PoliciesHTTPLocationMatch]    `tfsdk:"location"`
		RequestHeaders  supertypes.SetNestedObjectValueOf[PoliciesHTTPHeaderMatch]         `tfsdk:"request_headers"`
		ResponseHeaders supertypes.SetNestedObjectValueOf[PoliciesHTTPHeaderMatch]         `tfsdk:"response_headers"`
		StatusCode      supertypes.SingleNestedObjectValueOf[PoliciesHTTPStatusCodeMatch]  `tfsdk:"status_code"`
		Query           supertypes.SetValueOf[string]                                      `tfsdk:"query"`
	}

	PoliciesHTTPResponseActions struct {
		LocationRewrite supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionLocationRewrite] `tfsdk:"location_rewrite"`
		ModifyHeaders   supertypes.SetNestedObjectValueOf[PoliciesHTTPActionHeaderRewrite]      `tfsdk:"modify_headers"`
	}
)

func (rm *PoliciesHTTPResponseModel) Copy() *PoliciesHTTPResponseModel {
	x := &PoliciesHTTPResponseModel{}
	utils.ModelCopy(rm, x)
	return x
}

// ToSDKPoliciesHTTPResponseGroupModel converts the model to the SDK model.
func (rm *PoliciesHTTPResponseModel) ToSDKPoliciesHTTPResponseModel(ctx context.Context) (*edgeloadbalancer.PoliciesHTTPResponseModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &edgeloadbalancer.PoliciesHTTPResponseModel{
		VirtualServiceID: rm.VirtualServiceID.ValueString(),
		Policies:         make([]*edgeloadbalancer.PoliciesHTTPResponseModelPolicy, 0),
	}

	for _, policy := range rm.Policies.DiagsGet(ctx, diags) {
		actions := func() *PoliciesHTTPResponseActions {
			if !policy.Actions.IsKnown() {
				return nil
			}

			return policy.Actions.DiagsGet(ctx, diags)
		}()

		model.Policies = append(model.Policies, &edgeloadbalancer.PoliciesHTTPResponseModelPolicy{
			Name:    policy.Name.ValueString(),
			Active:  policy.Active.ValueBool(),
			Logging: policy.Logging.ValueBool(),
			MatchCriteria: func() edgeloadbalancer.PoliciesHTTPResponseMatchCriteria {
				if !policy.Criteria.IsKnown() {
					return edgeloadbalancer.PoliciesHTTPResponseMatchCriteria{}
				}

				criteria := policy.Criteria.DiagsGet(ctx, diags)
				return edgeloadbalancer.PoliciesHTTPResponseMatchCriteria{
					Protocol:            criteria.Protocol.Get(),
					ClientIPMatch:       policiesHTTPClientIPMatchToSDK(ctx, diags, criteria.ClientIP),
					ServicePortMatch:    policiesHTTPServicePortMatchToSDK(ctx, diags, criteria.ServicePorts),
					MethodMatch:         policiesHTTPMethodMatchToSDK(ctx, diags, criteria.HTTPMethods),
					PathMatch:           policiesHTTPPathMatchToSDK(ctx, diags, criteria.Path),
					CookieMatch:         policiesHTTPCookieMatchToSDK(ctx, diags, criteria.Cookie),
					LocationMatch:       policiesHTTPLocationMatchToSDK(ctx, diags, criteria.Location),
					RequestHeaderMatch:  policiesHTTPHeadersMatchToSDK(ctx, diags, criteria.RequestHeaders),
					ResponseHeaderMatch: policiesHTTPHeadersMatchToSDK(ctx, diags, criteria.ResponseHeaders),
					StatusCodeMatch:     policiesHTTPStatusCodeMatchToSDK(ctx, diags, criteria.StatusCode),
					QueryMatch:          criteria.Query.DiagsGet(ctx, diags),
				}
			}(),
			HeaderRewriteActions:  policiesHTTPActionHeadersRewriteToSDK(ctx, diags, actions.ModifyHeaders),
			LocationRewriteAction: policiesHTTPActionLocationRewriteToSDK(ctx, diags, actions.LocationRewrite),
		})
	}

	return model, diags
}

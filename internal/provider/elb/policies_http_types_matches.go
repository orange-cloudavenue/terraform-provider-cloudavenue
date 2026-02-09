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
)

type (

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
	PoliciesHTTPLocationMatch struct {
		Criteria supertypes.StringValue        `tfsdk:"criteria"`
		Values   supertypes.SetValueOf[string] `tfsdk:"values"`
	}
	PoliciesHTTPStatusCodeMatch struct {
		Criteria supertypes.StringValue        `tfsdk:"criteria"`
		Codes    supertypes.SetValueOf[string] `tfsdk:"codes"`
	}
)

// * ClientIPMatch

// policiesHTTPClientIPMatchToSDK converts the terraform model to the SDK model.
func policiesHTTPClientIPMatchToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPClientIPMatch]) *edgeloadbalancer.PoliciesHTTPClientIPMatch {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPClientIPMatch{
		Criteria:  v.Criteria.Get(),
		Addresses: v.IPAddresses.DiagsGet(ctx, diags),
	}
}

// policiesHTTPClientIPMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPClientIPMatchFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPClientIPMatch) supertypes.SingleNestedObjectValueOf[PoliciesHTTPClientIPMatch] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPClientIPMatch](ctx)
	}

	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPClientIPMatch{
		Criteria:    supertypes.NewStringValueOrNull(v.Criteria),
		IPAddresses: supertypes.NewSetValueOfSlice(ctx, v.Addresses),
	})
}

// * ServicePortMatch

// policiesHTTPServicePortMatchToSDK converts the terraform model to the SDK model.
func policiesHTTPServicePortMatchToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPServicePortMatch]) *edgeloadbalancer.PoliciesHTTPServicePortMatch {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPServicePortMatch{
		Criteria: v.Criteria.Get(),
		Ports: func(x []int64) (p []int) { // Convert int64 to int
			for _, v := range x {
				p = append(p, int(v))
			}
			return p
		}(v.Ports.DiagsGet(ctx, diags)),
	}
}

// policiesHTTPServicePortMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPServicePortMatchFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPServicePortMatch) supertypes.SingleNestedObjectValueOf[PoliciesHTTPServicePortMatch] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPServicePortMatch](ctx)
	}
	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPServicePortMatch{
		Criteria: supertypes.NewStringValueOrNull(v.Criteria),
		// TODO Wait new Int type in supertypes
		Ports: func() supertypes.SetValueOf[int64] {
			ports := []int64{}
			for _, port := range v.Ports {
				ports = append(ports, int64(port))
			}
			return supertypes.NewSetValueOfSlice(ctx, ports)
		}(),
	})
}

// * MethodMatch

// policiesHTTPMethodMatchToSDK converts the terraform model to the SDK model.
func policiesHTTPMethodMatchToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPMethodMatch]) *edgeloadbalancer.PoliciesHTTPMethodMatch {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPMethodMatch{
		Criteria: v.Criteria.Get(),
		Methods:  v.Methods.DiagsGet(ctx, diags),
	}
}

// policiesHTTPMethodMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPMethodMatchFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPMethodMatch) supertypes.SingleNestedObjectValueOf[PoliciesHTTPMethodMatch] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPMethodMatch](ctx)
	}

	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPMethodMatch{
		Criteria: supertypes.NewStringValueOrNull(v.Criteria),
		Methods:  supertypes.NewSetValueOfSlice(ctx, v.Methods),
	})
}

// * PathMatch

// policiesHTTPPathMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPPathMatchToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPPathMatch]) *edgeloadbalancer.PoliciesHTTPPathMatch {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPPathMatch{
		Criteria:     v.Criteria.Get(),
		MatchStrings: v.Paths.DiagsGet(ctx, diags),
	}
}

// policiesHTTPPathMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPPathMatchFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPPathMatch) supertypes.SingleNestedObjectValueOf[PoliciesHTTPPathMatch] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPPathMatch](ctx)
	}

	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPPathMatch{
		Criteria: supertypes.NewStringValueOrNull(v.Criteria),
		Paths:    supertypes.NewSetValueOfSlice(ctx, v.MatchStrings),
	})
}

// * CookieMatch

// policiesHTTPCookieMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPCookieMatchToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPCookieMatch]) *edgeloadbalancer.PoliciesHTTPCookieMatch {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPCookieMatch{
		Criteria: v.Criteria.Get(),
		Name:     v.Name.Get(),
		Value:    v.Value.Get(),
	}
}

// policiesHTTPCookieMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPCookieMatchFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPCookieMatch) supertypes.SingleNestedObjectValueOf[PoliciesHTTPCookieMatch] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPCookieMatch](ctx)
	}

	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPCookieMatch{
		Criteria: supertypes.NewStringValueOrNull(v.Criteria),
		Name:     supertypes.NewStringValueOrNull(v.Name),
		Value:    supertypes.NewStringValueOrNull(v.Value),
	})
}

// * LocationMatch

// policiesHTTPLocationMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPLocationMatchToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPLocationMatch]) *edgeloadbalancer.PoliciesHTTPLocationMatch {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPLocationMatch{
		Criteria: v.Criteria.Get(),
		Values:   v.Values.DiagsGet(ctx, diags),
	}
}

// policiesHTTPLocationMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPLocationMatchFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPLocationMatch) supertypes.SingleNestedObjectValueOf[PoliciesHTTPLocationMatch] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPLocationMatch](ctx)
	}

	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPLocationMatch{
		Criteria: supertypes.NewStringValueOrNull(v.Criteria),
		Values:   supertypes.NewSetValueOfSlice(ctx, v.Values),
	})
}

// * StatusCodeMatch

// policiesHTTPStatusCodeMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPStatusCodeMatchToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPStatusCodeMatch]) *edgeloadbalancer.PoliciesHTTPStatusCodeMatch {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPStatusCodeMatch{
		Criteria:    v.Criteria.Get(),
		StatusCodes: v.Codes.DiagsGet(ctx, diags),
	}
}

// policiesHTTPStatusCodeMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPStatusCodeMatchFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPStatusCodeMatch) supertypes.SingleNestedObjectValueOf[PoliciesHTTPStatusCodeMatch] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPStatusCodeMatch](ctx)
	}

	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPStatusCodeMatch{
		Criteria: supertypes.NewStringValueOrNull(v.Criteria),
		Codes:    supertypes.NewSetValueOfSlice(ctx, v.StatusCodes),
	})
}

// * HeadersMatch

// policiesHTTPHeadersMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPHeadersMatchToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SetNestedObjectValueOf[PoliciesHTTPHeaderMatch]) edgeloadbalancer.PoliciesHTTPHeadersMatch {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	var headers edgeloadbalancer.PoliciesHTTPHeadersMatch
	for _, header := range v {
		headers = append(headers, edgeloadbalancer.PoliciesHTTPHeaderMatch{
			Criteria: header.Criteria.Get(),
			Name:     header.Name.Get(),
			Values:   header.Values.DiagsGet(ctx, diags),
		})
	}
	return headers
}

// policiesHTTPHeadersMatchFromSDK converts the SDK model to the terraform model.
func policiesHTTPHeadersMatchFromSDK(ctx context.Context, v edgeloadbalancer.PoliciesHTTPHeadersMatch) supertypes.SetNestedObjectValueOf[PoliciesHTTPHeaderMatch] {
	headers := []*PoliciesHTTPHeaderMatch{}
	for _, header := range v {
		headers = append(headers, &PoliciesHTTPHeaderMatch{
			Criteria: supertypes.NewStringValueOrNull(header.Criteria),
			Name:     supertypes.NewStringValueOrNull(header.Name),
			Values:   supertypes.NewSetValueOfSlice(ctx, header.Values),
		})
	}
	if len(headers) == 0 {
		return supertypes.NewSetNestedObjectValueOfNull[PoliciesHTTPHeaderMatch](ctx)
	}
	return supertypes.NewSetNestedObjectValueOfSlice(ctx, headers)
}

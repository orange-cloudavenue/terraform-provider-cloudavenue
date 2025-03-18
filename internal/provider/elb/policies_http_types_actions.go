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
	PoliciesHTTPActionLocationRewrite struct {
		Protocol  supertypes.StringValue `tfsdk:"protocol"`
		Host      supertypes.StringValue `tfsdk:"host"`
		Port      supertypes.Int64Value  `tfsdk:"port"`
		Path      supertypes.StringValue `tfsdk:"path"`
		KeepQuery supertypes.BoolValue   `tfsdk:"keep_query"`
	}
)

// * ActionRedirect

// policiesHTTPActionRedirectToSDK converts the terraform model to the SDK model.
func policiesHTTPActionRedirectToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionRedirect]) *edgeloadbalancer.PoliciesHTTPActionRedirect {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPActionRedirect{
		Host:       v.Host.Get(),
		KeepQuery:  v.KeepQuery.Get(),
		Path:       v.Path.Get(),
		Port:       v.Port.GetIntPtr(),
		Protocol:   v.Protocol.Get(),
		StatusCode: v.StatusCode.GetInt(),
	}
}

// policiesHTTPActionRedirectFromSDK converts the SDK model to the terraform model.
func policiesHTTPActionRedirectFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPActionRedirect) supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionRedirect] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPActionRedirect](ctx)
	}

	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPActionRedirect{
		Host:      supertypes.NewStringValueOrNull(v.Host),
		KeepQuery: supertypes.NewBoolValue(v.KeepQuery),
		Path:      supertypes.NewStringValueOrNull(v.Path),
		// TODO wait new Int type in supertypes
		Port: func() supertypes.Int64Value {
			if v.Port == nil {
				return supertypes.NewInt64Null()
			}
			return supertypes.NewInt64Value(int64(*v.Port))
		}(),
		Protocol: supertypes.NewStringValueOrNull(v.Protocol),
		// TODO wait new Int type in supertypes
		StatusCode: func() supertypes.Int64Value {
			return supertypes.NewInt64Value(int64(v.StatusCode))
		}(),
	})
}

// * ActionURLRewrite

// policiesHTTPActionURLRewriteToSDK converts the terraform model to the SDK model.
func policiesHTTPActionURLRewriteToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionURLRewrite]) *edgeloadbalancer.PoliciesHTTPActionURLRewrite {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPActionURLRewrite{
		HostHeader: v.Host.Get(),
		Path:       v.Path.Get(),
		Query:      v.Query.Get(),
		KeepQuery:  v.KeepQuery.Get(),
	}
}

// policiesHTTPActionURLRewriteFromSDK converts the SDK model to the terraform model.
func policiesHTTPActionURLRewriteFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPActionURLRewrite) supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionURLRewrite] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPActionURLRewrite](ctx)
	}
	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPActionURLRewrite{
		Host:      supertypes.NewStringValueOrNull(v.HostHeader),
		Path:      supertypes.NewStringValueOrNull(v.Path),
		Query:     supertypes.NewStringValueOrNull(v.Query),
		KeepQuery: supertypes.NewBoolValue(v.KeepQuery),
	})
}

// * ActionHeadersRewrite

// policiesHTTPActionHeadersRewriteToSDK converts the terraform model to the SDK model.
func policiesHTTPActionHeadersRewriteToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SetNestedObjectValueOf[PoliciesHTTPActionHeaderRewrite]) edgeloadbalancer.PoliciesHTTPActionHeadersRewrite {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	var headers edgeloadbalancer.PoliciesHTTPActionHeadersRewrite
	for _, header := range v {
		headers = append(headers, &edgeloadbalancer.PoliciesHTTPActionHeaderRewrite{
			Action: header.Action.Get(),
			Name:   header.Name.Get(),
			Value:  header.Value.Get(),
		})
	}
	return headers
}

// policiesHTTPActionHeadersRewriteFromSDK converts the SDK model to the terraform model.
func policiesHTTPActionHeadersRewriteFromSDK(ctx context.Context, v edgeloadbalancer.PoliciesHTTPActionHeadersRewrite) supertypes.SetNestedObjectValueOf[PoliciesHTTPActionHeaderRewrite] {
	headers := []*PoliciesHTTPActionHeaderRewrite{}
	for _, header := range v {
		headers = append(headers, &PoliciesHTTPActionHeaderRewrite{
			Action: supertypes.NewStringValueOrNull(header.Action),
			Name:   supertypes.NewStringValueOrNull(header.Name),
			Value:  supertypes.NewStringValueOrNull(header.Value),
		})
	}
	if len(headers) == 0 {
		return supertypes.NewSetNestedObjectValueOfNull[PoliciesHTTPActionHeaderRewrite](ctx)
	}
	return supertypes.NewSetNestedObjectValueOfSlice(ctx, headers)
}

// * ActionLocationRewrite

// policiesHTTPActionLocationRewriteToSDK converts the terraform model to the SDK model.
func policiesHTTPActionLocationRewriteToSDK(ctx context.Context, diags diag.Diagnostics, s supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionLocationRewrite]) *edgeloadbalancer.PoliciesHTTPActionLocationRewrite {
	if !s.IsKnown() {
		return nil
	}

	v := s.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil
	}

	return &edgeloadbalancer.PoliciesHTTPActionLocationRewrite{
		Protocol:  v.Protocol.Get(),
		Host:      v.Host.Get(),
		Port:      v.Port.GetIntPtr(),
		Path:      v.Path.Get(),
		KeepQuery: v.KeepQuery.Get(),
	}
}

// policiesHTTPActionLocationRewriteFromSDK converts the SDK model to the terraform model.
func policiesHTTPActionLocationRewriteFromSDK(ctx context.Context, v *edgeloadbalancer.PoliciesHTTPActionLocationRewrite) supertypes.SingleNestedObjectValueOf[PoliciesHTTPActionLocationRewrite] {
	if v == nil {
		return supertypes.NewSingleNestedObjectValueOfNull[PoliciesHTTPActionLocationRewrite](ctx)
	}

	return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPActionLocationRewrite{
		Protocol: supertypes.NewStringValueOrNull(v.Protocol),
		Host:     supertypes.NewStringValueOrNull(v.Host),
		Port: func() supertypes.Int64Value {
			if v.Port == nil {
				return supertypes.NewInt64Null()
			}
			return supertypes.NewInt64Value(int64(*v.Port))
		}(),
		Path:      supertypes.NewStringValueOrNull(v.Path),
		KeepQuery: supertypes.NewBoolValue(v.KeepQuery),
	})
}

/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"testing"

	sdkv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

func TestAttributesFromSDKProfileAppIDSubAttributeNullWhenEmpty(t *testing.T) {
	t.Parallel()

	appID, _, diags := attributesFromSDKProfile(t.Context(), &sdkv1.NetworkContextProfile{
		Attributes: []sdkv1.NetworkContextProfileAttribute{
			{
				Type:   sdkv1.NetworkContextProfileAttributeTypeAppID,
				Values: []string{"SSH", "DNS"},
			},
		},
	})
	if diags.HasError() {
		t.Fatalf("attributesFromSDKProfile() diagnostics: %+v", diags)
	}
	if appID == nil {
		t.Fatal("expected app_id block, got nil")
	}
	if !appID.SubAttribute.IsNull() {
		t.Fatal("expected app_id.sub_attributes to be null when API sub-attributes are absent")
	}
}

func TestAttributesFromSDKProfileAppIDSubAttributeSetWhenPresent(t *testing.T) {
	t.Parallel()

	appID, _, diags := attributesFromSDKProfile(t.Context(), &sdkv1.NetworkContextProfile{
		Attributes: []sdkv1.NetworkContextProfileAttribute{
			{
				Type:   sdkv1.NetworkContextProfileAttributeTypeAppID,
				Values: []string{"SSL"},
				SubAttributes: []sdkv1.NetworkContextProfileSubAttribute{
					{
						Type:   sdkv1.NetworkContextProfileSubAttributeTypeTLSVersion,
						Values: []string{"TLS_V12", "TLS_V13"},
					},
				},
			},
		},
	})
	if diags.HasError() {
		t.Fatalf("attributesFromSDKProfile() diagnostics: %+v", diags)
	}
	if appID == nil {
		t.Fatal("expected app_id block, got nil")
	}
	if appID.SubAttribute.IsNull() {
		t.Fatal("expected app_id.sub_attributes to be set when API sub-attributes are present")
	}
}

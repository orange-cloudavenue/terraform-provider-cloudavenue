// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package iam

const (
	attrName        = "name"
	attrDescription = "description"
	attrRights      = "rights"
)

// User provider types as returned by the SDK v2 CloudAPI leaves.
// CloudAPI reports local users as LOCAL while the legacy VCD user API
// (and the provider contract) exposes them as INTEGRATED.
const (
	sdkUserProviderTypeLocal  = "LOCAL"
	sdkUserProviderTypeSAML   = "SAML"
	legacyUserProviderTypeInt = "INTEGRATED"
)

// quotaValue dereferences a quota pointer, treating unset as zero.
func quotaValue(quota *int) int {
	if quota == nil {
		return 0
	}
	return *quota
}

// legacyProviderType maps an SDK v2 CloudAPI user provider type to the
// legacy VCD user API value exposed by the provider (LOCAL -> INTEGRATED).
func legacyProviderType(providerType string) string {
	if providerType == sdkUserProviderTypeLocal {
		return legacyUserProviderTypeInt
	}
	return providerType
}

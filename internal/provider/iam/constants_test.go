// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package iam

import "testing"

func TestMapLegacyProviderType(t *testing.T) {
	for in, want := range map[string]string{
		"LOCAL":      "INTEGRATED",
		"SAML":       "SAML",
		"LDAP":       "LDAP",
		"INTEGRATED": "INTEGRATED",
	} {
		if got := mapLegacyProviderType(in); got != want {
			t.Errorf("mapLegacyProviderType(%q) = %q, want %q", in, got, want)
		}
	}
}

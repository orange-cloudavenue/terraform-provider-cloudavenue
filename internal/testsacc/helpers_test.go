/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// SPDX-FileCopyrightText: Copyright (c) 2025 Orange
// SPDX-License-Identifier: Mozilla Public License 2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package testsacc

import (
	"testing"

	"github.com/google/uuid"
)

func TestToValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		validator string
		wantErr   bool
	}{
		// uuid4
		{"valid uuid v4", uuid.New().String(), "uuid4", false},
		{"empty string", "", "uuid4", true},
		{"invalid uuid (wrong version)", func() string {
			v, _ := uuid.NewV7()
			return v.String()
		}(), "uuid4", true},
		{"invalid uuid (not uuid)", "not-a-uuid", "uuid4", true},
		{"invalid uuid (v4 but wrong chars)", "zzzzzzzz-zzzz-4zzz-8zzz-zzzzzzzzzzzz", "uuid4", true},
		// email
		{"valid email", "test@example.com", "email", false},
		{"invalid email", "not-an-email", "email", true},
		// urn
		{"valid urn", "urn:example:resource", "urn_rfc2141", false},
		{"invalid urn", "not-a-urn", "urn_rfc2141", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ToValidate(tt.validator)(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToValidate(%q) error = %v, wantErr %v", tt.validator, err, tt.wantErr)
			}
		})
	}
}

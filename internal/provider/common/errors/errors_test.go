// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0

/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package errors

import (
	"fmt"
	"testing"

	v2errors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

func TestIsNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error is not not-found",
			err:  nil,
			want: false,
		},
		{
			name: "unrelated error is not not-found",
			err:  fmt.Errorf("boom"),
			want: false,
		},
		{
			name: "bare SDK v2 not-found sentinel is not-found",
			err:  v2errors.ErrNotFound,
			want: true,
		},
		{
			name: "wrapped SDK v2 not-found sentinel is not-found",
			err:  fmt.Errorf("wrap: %w", v2errors.ErrNotFound),
			want: true,
		},
		{
			name: "doubly wrapped SDK v2 not-found sentinel is not-found",
			err:  fmt.Errorf("outer: %w", fmt.Errorf("wrap: %w", v2errors.ErrNotFound)),
			want: true,
		},
		{
			name: "wrapped unrelated SDK v2 error is not not-found",
			err:  fmt.Errorf("wrap: %w", v2errors.ErrBadRequest),
			want: false,
		},
		{
			// SDK v2 lookups by id-or-name (e.g. GetSecurityGroup,
			// GetVDCNetwork, GetCertificate) build the APIError directly with
			// StatusCode 404 and leave Err unset, so errors.Is(err,
			// v2errors.ErrNotFound) cannot match: Unwrap() returns nil.
			name: "SDK v2 APIError with status 404 and unset Err is not-found",
			err:  &v2errors.APIError{Operation: "GetSecurityGroup", StatusCode: 404, Message: `security group "x" not found`},
			want: true,
		},
		{
			name: "wrapped SDK v2 APIError with status 404 and unset Err is not-found",
			err: fmt.Errorf("reading security group: %w", &v2errors.APIError{
				Operation:  "GetSecurityGroup",
				StatusCode: 404,
				Message:    `security group "x" not found`,
			}),
			want: true,
		},
		{
			name: "doubly wrapped SDK v2 APIError with status 404 and unset Err is not-found",
			err: fmt.Errorf("outer: %w", fmt.Errorf("wrap: %w", &v2errors.APIError{
				Operation:  "GetVDCNetwork",
				StatusCode: 404,
				Message:    `vdc network "x" not found`,
			})),
			want: true,
		},
		{
			name: "SDK v2 APIError with status 400 and unset Err is not not-found",
			err:  &v2errors.APIError{Operation: "GetSecurityGroup", StatusCode: 400, Message: "id or name is required"},
			want: false,
		},
		{
			name: "SDK v2 APIError with status 500 and ErrNotFound sentinel is not-found",
			err: &v2errors.APIError{
				Operation:  "GetSecurityGroup",
				StatusCode: 500,
				Err:        v2errors.ErrNotFound,
			},
			want: true,
		},
		{
			// Known gap, documented on purpose: ClientError embeds APIError by
			// value and has no Unwrap, so it is opaque to errors.As and no
			// SDK request path can return one (nothing outside the SDK's own
			// tests constructs a *ClientError).
			name: "SDK v2 ClientError carrying a 404 APIError is not detected (opaque, no Unwrap)",
			err: &v2errors.ClientError{
				Message: "client side failure",
				APIErr:  v2errors.APIError{Operation: "GetSecurityGroup", StatusCode: 404},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IsNotFound(tt.err); got != tt.want {
				t.Fatalf("IsNotFound(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

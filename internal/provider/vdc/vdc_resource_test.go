// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vdc

import (
	"errors"
	"testing"
)

func TestIsMissingJobIDCreateError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "matches sdk create error",
			err:  errors.New("error on create VDC: cannot refresh job status: job ID is empty"),
			want: true,
		},
		{
			name: "case insensitive match",
			err:  errors.New("JOB ID IS EMPTY"),
			want: true,
		},
		{
			name: "other create error",
			err:  errors.New("error on create VDC: api rejected request"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isMissingJobIDCreateError(tt.err); got != tt.want {
				t.Fatalf("isMissingJobIDCreateError() = %v, want %v", got, tt.want)
			}
		})
	}
}

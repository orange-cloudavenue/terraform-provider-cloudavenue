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

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func newInt64Unknown() supertypes.Int64Value {
	return supertypes.NewInt64Unknown()
}

func newInt64Value(v int64) supertypes.Int64Value {
	return supertypes.NewInt64Value(v)
}

func TestDetermineModifyPlanBandwidthAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		plan  edgeGatewayResourceModel
		state edgeGatewayResourceModel
		want  modifyPlanBandwidthAction
	}{
		{
			name: "unknown plan with known state is update",
			plan: edgeGatewayResourceModel{
				Bandwidth: newInt64Unknown(),
			},
			state: edgeGatewayResourceModel{
				Bandwidth: newInt64Value(250),
			},
			want: modifyPlanBandwidthActionUpdate,
		},
		{
			name: "unknown plan with unknown state is create unknown",
			plan: edgeGatewayResourceModel{
				Bandwidth: newInt64Unknown(),
			},
			state: edgeGatewayResourceModel{
				Bandwidth: newInt64Unknown(),
			},
			want: modifyPlanBandwidthActionCreateUnknown,
		},
		{
			name: "known plan with unknown state is create known",
			plan: edgeGatewayResourceModel{
				Bandwidth: newInt64Value(100),
			},
			state: edgeGatewayResourceModel{
				Bandwidth: newInt64Unknown(),
			},
			want: modifyPlanBandwidthActionCreateKnown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := determineModifyPlanBandwidthAction(&tt.plan, &tt.state)
			if got != tt.want {
				t.Fatalf("determineModifyPlanBandwidthAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetermineModifyPlanBandwidthActionNilState(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("determineModifyPlanBandwidthAction() panicked: %v", r)
		}
	}()

	got := determineModifyPlanBandwidthAction(&edgeGatewayResourceModel{Bandwidth: newInt64Value(100)}, nil)
	if got != modifyPlanBandwidthActionCreateKnown {
		t.Fatalf("determineModifyPlanBandwidthAction() = %v, want %v", got, modifyPlanBandwidthActionCreateKnown)
	}
}

func TestModifyPlanUpdateUsesStateBandwidthWhenPlanUnknown(t *testing.T) {
	t.Parallel()

	plan := &edgeGatewayResourceModel{Bandwidth: newInt64Unknown()}
	state := &edgeGatewayResourceModel{Bandwidth: newInt64Value(250)}

	if got := determineModifyPlanBandwidthAction(plan, state); got != modifyPlanBandwidthActionUpdate {
		t.Fatalf("determineModifyPlanBandwidthAction() = %v, want update", got)
	}
}

func TestMaxAllowedBandwidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values []int
		want   int
		ok     bool
	}{
		{
			name:   "sorted input",
			values: []int{5, 25, 100, 1000},
			want:   1000,
			ok:     true,
		},
		{
			name:   "unsorted input",
			values: []int{1000, 5, 250, 50},
			want:   1000,
			ok:     true,
		},
		{
			name:   "empty input",
			values: nil,
			want:   0,
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := maxAllowedBandwidth(tt.values)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("maxAllowedBandwidth(%v) = (%d, %t), want (%d, %t)", tt.values, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestAllowedBandwidthAtMostUpdate(t *testing.T) {
	t.Parallel()

	got := allowedBandwidthAtMostUpdate(250, 100, []int{1000, 25, 500, 250})
	if got != 250 {
		t.Fatalf("allowedBandwidthAtMostUpdate() = %d, want %d", got, 250)
	}
}

func TestBestValueAtMostOrErrorNoFit(t *testing.T) {
	t.Parallel()

	got, err := bestValueAtMostOrError(4, []int{5, 25, 100})
	if err == nil {
		t.Fatalf("bestValueAtMostOrError() err = nil, want error")
	}

	if got != 0 {
		t.Fatalf("bestValueAtMostOrError() = %d, want 0", got)
	}

	if err.Error() != "no allowed bandwidth value fits current available capacity" {
		t.Fatalf("bestValueAtMostOrError() err = %q, want %q", err.Error(), "no allowed bandwidth value fits current available capacity")
	}
}

func TestDedicatedT0UpdateBandwidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		state   int
		allowed []int
		want    int
		wantErr bool
	}{
		{
			name:    "state bandwidth kept when allowed",
			state:   250,
			allowed: []int{25, 50, 250, 500},
			want:    250,
		},
		{
			name:    "fallback uses deterministic max allowed",
			state:   123,
			allowed: []int{25, 50, 250, 500},
			want:    500,
		},
		{
			name:    "empty allowed errors",
			state:   123,
			allowed: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := dedicatedT0UpdateBandwidth(tt.state, tt.allowed)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("dedicatedT0UpdateBandwidth() err = nil, want error")
				}
				return
			}

			if err != nil {
				t.Fatalf("dedicatedT0UpdateBandwidth() err = %v, want nil", err)
			}

			if got != tt.want {
				t.Fatalf("dedicatedT0UpdateBandwidth() = %d, want %d", got, tt.want)
			}
			if got == 0 {
				t.Fatalf("dedicatedT0UpdateBandwidth() = 0, want non-zero")
			}
		})
	}
}

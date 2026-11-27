// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vdc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	tfdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi/rules"
)

func TestValidateStorageProfilesDisponibilityCompatibility(t *testing.T) {
	t.Parallel()

	diags := validateStorageProfilesDisponibilityCompatibility(
		rules.ServiceClass("STD"),
		rules.DisponibilityClass("ONE-ROOM"),
		[]*vdcResourceModelVDCStorageProfile{{
			Class: supertypes.NewStringValue("gold_r1"),
			Limit: supertypes.NewInt64Value(500),
		}},
	)

	if !diags.HasError() {
		t.Fatal("expected compatibility validation error")
	}

	diagnostic := diags[0]
	if diagnostic.Severity() != tfdiag.SeverityError {
		t.Fatalf("expected error severity, got %v", diagnostic.Severity())
	}

	if got, want := diagnostic.Summary(), "Storage profile class is not compatible with disponibility class"; got != want {
		t.Fatalf("unexpected summary: got %q want %q", got, want)
	}

	detail := diagnostic.Detail()
	for _, expected := range []string{"gold_r1", "ONE-ROOM", "STD", "DUAL-ROOM"} {
		if !strings.Contains(detail, expected) {
			t.Fatalf("expected detail to contain %q, got %q", expected, detail)
		}
	}

	attrDiag, ok := diagnostic.(interface{ Path() path.Path })
	if !ok {
		t.Fatalf("expected diagnostic with attribute path, got %T", diagnostic)
	}

	if got, want := attrDiag.Path().String(), "storage_profiles[0].class"; got != want {
		t.Fatalf("unexpected path: got %q want %q", got, want)
	}
}

func TestValidateStorageProfilesDisponibilityCompatibilityCustomSuffixes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		disponibilityClass rules.DisponibilityClass
		storageClass       string
		wantError          bool
	}{
		{name: "custom one-room accepted", disponibilityClass: rules.DisponibilityClass("ONE-ROOM"), storageClass: "custom", wantError: false},
		{name: "custom dual-room r1 accepted", disponibilityClass: rules.DisponibilityClass("DUAL-ROOM"), storageClass: "custom_r1", wantError: false},
		{name: "custom dual-room r2 accepted", disponibilityClass: rules.DisponibilityClass("DUAL-ROOM"), storageClass: "custom_r2", wantError: false},
		{name: "custom ha-dual-room accepted", disponibilityClass: rules.DisponibilityClass("HA-DUAL-ROOM"), storageClass: "custom_hm", wantError: false},
		{name: "custom r1 rejected for one-room", disponibilityClass: rules.DisponibilityClass("ONE-ROOM"), storageClass: "custom_r1", wantError: true},
		{name: "custom no suffix rejected for dual-room", disponibilityClass: rules.DisponibilityClass("DUAL-ROOM"), storageClass: "custom", wantError: true},
		{name: "custom hm rejected for dual-room", disponibilityClass: rules.DisponibilityClass("DUAL-ROOM"), storageClass: "custom_hm", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := validateStorageProfilesDisponibilityCompatibility(
				rules.ServiceClass("STD"),
				tt.disponibilityClass,
				[]*vdcResourceModelVDCStorageProfile{{
					Class: supertypes.NewStringValue(tt.storageClass),
					Limit: supertypes.NewInt64Value(500),
				}},
			)

			if got := diags.HasError(); got != tt.wantError {
				t.Fatalf("validateStorageProfilesDisponibilityCompatibility(%q, %q) error=%v want %v; diags=%v", tt.disponibilityClass, tt.storageClass, got, tt.wantError, diags)
			}
		})
	}
}

func TestValidateConfigStorageProfilesDisponibilityCompatibility(t *testing.T) {
	t.Parallel()

	resource := &vdcResource{}
	storageProfiles := supertypes.NewSetNestedObjectValueOfValueSlice(context.Background(), []vdcResourceModelVDCStorageProfile{
		{
			Class:   supertypes.NewStringValue("gold_r1"),
			Limit:   supertypes.NewInt64Value(500),
			Default: supertypes.NewBoolValue(true),
		},
	})

	config := &vdcResourceModel{
		ServiceClass:        supertypes.NewStringValue("STD"),
		DisponibilityClass:  supertypes.NewStringValue("ONE-ROOM"),
		BillingModel:        supertypes.NewStringValue("PAYG"),
		VCPUInMhz:           supertypes.NewInt64Value(2200),
		CPUAllocated:        supertypes.NewInt64Value(22000),
		MemoryAllocated:     supertypes.NewInt64Value(30),
		StorageBillingModel: supertypes.NewStringValue("PAYG"),
		StorageProfiles:     storageProfiles,
	}

	validationDiags := resource.validateConfig(context.Background(), config)
	if !validationDiags.HasError() {
		t.Fatal("expected config validation error")
	}

	if got, want := validationDiags[0].Summary(), "Storage profile class is not compatible with disponibility class"; got != want {
		t.Fatalf("unexpected summary: got %q want %q", got, want)
	}
}

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

func TestCompatibilityDiagnosticPathMatchesTerraformSetTraversal(t *testing.T) {
	t.Parallel()

	storageProfiles := supertypes.NewSetNestedObjectValueOfValueSlice(context.Background(), []vdcResourceModelVDCStorageProfile{
		{
			Class:   supertypes.NewStringValue("gold_r1"),
			Limit:   supertypes.NewInt64Value(500),
			Default: supertypes.NewBoolValue(true),
		},
	})

	profiles, diags := storageProfiles.Get(context.Background())
	if diags.HasError() {
		t.Fatalf("unexpected storage profile conversion diagnostics: %v", diags)
	}

	compatDiags := validateStorageProfilesDisponibilityCompatibility(
		rules.ServiceClass("STD"),
		rules.DisponibilityClass("ONE-ROOM"),
		profiles,
	)
	if !compatDiags.HasError() {
		t.Fatal("expected compatibility diagnostics")
	}

	attrDiag, ok := compatDiags[0].(interface{ Path() path.Path })
	if !ok {
		t.Fatalf("expected diagnostic with attribute path, got %T", compatDiags[0])
	}

	if got, want := attrDiag.Path().String(), fmt.Sprintf("storage_profiles[%d].class", 0); got != want {
		t.Fatalf("unexpected path: got %q want %q", got, want)
	}
}

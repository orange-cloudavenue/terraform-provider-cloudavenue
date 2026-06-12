/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg_test

import (
	"context"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vdcg"
)

// networkContextProfileVDCGResourceSchemaAttrs are the top-level attributes expected in the resource schema.
var networkContextProfileVDCGResourceSchemaAttrs = []string{
	"id", "name", "description", "scope", "app_id", "vdc_group_id", "vdc_group_name",
}

// networkContextProfileVDCGDatasourceSchemaAttrs are the top-level attributes expected in the datasource schema.
// domain_name is Computed-only in the datasource (read from existing SYSTEM profiles).
var networkContextProfileVDCGDatasourceSchemaAttrs = []string{
	"id", "name", "description", "scope", "app_id", "domain_name", "vdc_group_id", "vdc_group_name",
}

func TestNetworkContextProfileVdcgSchemas(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	tests := map[string]struct {
		validate func(ctx context.Context, t *testing.T)
	}{
		"resource": {
			validate: func(ctx context.Context, t *testing.T) {
				t.Helper()
				resp := &fwresource.SchemaResponse{}
				vdcg.NewNetworkContextProfileResource().Schema(ctx, fwresource.SchemaRequest{}, resp)
				if resp.Diagnostics.HasError() {
					t.Fatalf("Schema() diagnostics: %+v", resp.Diagnostics)
				}
				if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
					t.Fatalf("ValidateImplementation() diagnostics: %+v", diags)
				}
				for _, attr := range networkContextProfileVDCGResourceSchemaAttrs {
					if _, ok := resp.Schema.Attributes[attr]; !ok {
						t.Errorf("expected attribute %q in resource schema", attr)
					}
				}
			},
		},
		"datasource": {
			validate: func(ctx context.Context, t *testing.T) {
				t.Helper()
				resp := &fwdatasource.SchemaResponse{}
				vdcg.NewNetworkContextProfileDataSource().Schema(ctx, fwdatasource.SchemaRequest{}, resp)
				if resp.Diagnostics.HasError() {
					t.Fatalf("Schema() diagnostics: %+v", resp.Diagnostics)
				}
				if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
					t.Fatalf("ValidateImplementation() diagnostics: %+v", diags)
				}
				for _, attr := range networkContextProfileVDCGDatasourceSchemaAttrs {
					if _, ok := resp.Schema.Attributes[attr]; !ok {
						t.Errorf("expected attribute %q in datasource schema", attr)
					}
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tt.validate(ctx, t)
		})
	}
}

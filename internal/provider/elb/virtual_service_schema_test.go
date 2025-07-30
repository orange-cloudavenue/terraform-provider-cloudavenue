/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package elb_test

import (
	"testing"

	// fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource".
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/elb"
)

// Unit test for the schema of the resource cloudavenue_elb_virtual_service.
func TestVirtualServiceResourceSchema(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	schemaResponse := &fwresource.SchemaResponse{}

	// Instantiate the resource.Resource and call its Schema method
	elb.NewVirtualServiceResource().Schema(ctx, fwresource.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

// Unit test for the schema of the datasource cloudavenue_elb_virtual_service.
func TestVirtualServiceDataSourceSchema(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	schemaResponse := &fwdatasource.SchemaResponse{}

	// Instantiate the datasource.Datasource and call its Schema method
	elb.NewVirtualServiceDataSource().Schema(ctx, fwdatasource.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

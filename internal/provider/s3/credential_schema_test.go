/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package s3_test

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/s3"
)

// Unit test for the schema of the resource cloudavenue_s3_Credential.
func TestCredentialResourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaResponse := &fwresource.SchemaResponse{}

	// Instantiate the resource.Resource and call its Schema method
	s3.NewCredentialResource().Schema(ctx, fwresource.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package provider

import (
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	fwschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestProviderSchema(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	schemaResponse := &fwprovider.SchemaResponse{}

	New("test")().Schema(ctx, fwprovider.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diagnostics)
	}

	attribute, ok := schemaResponse.Schema.Attributes["core_api"]
	if !ok {
		t.Fatal("expected core_api attribute in provider schema")
	}

	coreAPI, ok := attribute.(fwschema.StringAttribute)
	if !ok {
		t.Fatalf("expected core_api to be a string attribute, got %T", attribute)
	}

	if !coreAPI.IsOptional() {
		t.Fatal("expected core_api to be optional")
	}

	if coreAPI.GetMarkdownDescription() == "" {
		t.Fatal("expected core_api markdown description")
	}

	urlAttribute, ok := schemaResponse.Schema.Attributes["url"]
	if !ok {
		t.Fatal("expected url attribute to remain in provider schema")
	}

	url, ok := urlAttribute.(fwschema.StringAttribute)
	if !ok {
		t.Fatalf("expected url to be a string attribute, got %T", urlAttribute)
	}

	if !url.IsOptional() {
		t.Fatal("expected url to remain optional")
	}
}

func TestProviderClientOptsMapsCoreAPI(t *testing.T) {
	t.Parallel()

	config := cloudavenueProviderModel{
		URL:      types.StringValue("https://vcd.example.com"),
		CoreAPI:  types.StringValue("https://core-api.example.com"),
		User:     types.StringValue("test-user"),
		Password: types.StringValue("test-password"),
		Org:      types.StringValue("cav01ev01ocb0001234"),
		VDC:      types.StringNull(),
	}

	opts := providerClientOpts(config)

	if opts == nil || opts.CloudAvenue == nil {
		t.Fatal("expected CloudAvenue client options to be configured")
	}

	if opts.CloudAvenue.URL != "https://vcd.example.com" {
		t.Fatalf("expected url to map to VMware/VCD endpoint, got %q", opts.CloudAvenue.URL)
	}

	if opts.CloudAvenue.CoreAPI != "https://core-api.example.com" {
		t.Fatalf("expected core_api to map to SDK CoreAPI, got %q", opts.CloudAvenue.CoreAPI)
	}

	if opts.CloudAvenue.Username != "test-user" {
		t.Fatalf("expected user to map to SDK username, got %q", opts.CloudAvenue.Username)
	}

	if opts.CloudAvenue.Password != "test-password" {
		t.Fatalf("expected password to map to SDK password, got %q", opts.CloudAvenue.Password)
	}

	if opts.CloudAvenue.Org != "cav01ev01ocb0001234" {
		t.Fatalf("expected org to map to SDK org, got %q", opts.CloudAvenue.Org)
	}
}

func TestProviderClientOptsKeepsCoreAPIEmptyWhenUnset(t *testing.T) {
	t.Parallel()

	opts := providerClientOpts(cloudavenueProviderModel{
		URL:      types.StringValue("https://vcd.example.com"),
		CoreAPI:  types.StringNull(),
		User:     types.StringValue("test-user"),
		Password: types.StringValue("test-password"),
		Org:      types.StringValue("cav01ev01ocb0001234"),
	})

	if opts.CloudAvenue.CoreAPI != "" {
		t.Fatalf("expected empty core_api when unset, got %q", opts.CloudAvenue.CoreAPI)
	}

	if opts.CloudAvenue.URL != "https://vcd.example.com" {
		t.Fatalf("expected url to remain unchanged, got %q", opts.CloudAvenue.URL)
	}
}

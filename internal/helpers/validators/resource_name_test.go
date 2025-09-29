/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package validators

import (
	"context"
	"regexp"
	"testing"

	"github.com/orange-cloudavenue/common-go/regex"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestResourceNameValidator_ValidateString_Valid(t *testing.T) {
	// Setup a fake regex resource for testing
	regex.ListCavResourceNames = []regex.CavResourceName{
		{
			Key:         "test",
			Description: "Test resource name",
			RegexString: "^[a-z]{3,10}$",
			RegexP:      regexp.MustCompile("^[a-z]{3,10}$"),
		},
	}

	v := ResourceName("test")
	req := validator.StringRequest{
		ConfigValue: types.StringValue("validname"),
		Path:        path.Root("name"),
	}
	resp := &validator.StringResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateString(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

func TestResourceNameValidator_ValidateString_Invalid(t *testing.T) {
	regex.ListCavResourceNames = []regex.CavResourceName{
		{
			Key:         "test",
			Description: "Test resource name",
			RegexString: "^[a-z]{3,10}$",
			RegexP:      regexp.MustCompile("^[a-z]{3,10}$"),
		},
	}

	v := ResourceName("test")
	req := validator.StringRequest{
		ConfigValue: types.StringValue("INVALID123"),
		Path:        path.Root("name"),
	}
	resp := &validator.StringResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateString(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
}

func TestResourceNameValidator_ValidateString_UnknownOrNull(t *testing.T) {
	v := ResourceName("test")
	req := validator.StringRequest{
		ConfigValue: types.StringNull(),
		Path:        path.Root("name"),
	}
	resp := &validator.StringResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateString(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError())

	req.ConfigValue = types.StringUnknown()
	resp.Diagnostics = diag.Diagnostics{}
	v.ValidateString(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

func TestResourceNameValidator_ValidateString_UnknownResourceType(t *testing.T) {
	regex.ListCavResourceNames = []regex.CavResourceName{}
	v := ResourceName("unknown")
	req := validator.StringRequest{
		ConfigValue: types.StringValue("anyvalue"),
		Path:        path.Root("name"),
	}
	resp := &validator.StringResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateString(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

func TestResourceNameValidator_Description(t *testing.T) {
	regex.ListCavResourceNames = []regex.CavResourceName{
		{
			Key:         "test",
			Description: "Test description",
			RegexString: ".*",
			RegexP:      regexp.MustCompile(".*"),
		},
	}
	v := ResourceName("test")
	desc := v.Description(context.Background())
	assert.Equal(t, "Test description", desc)
}

func TestResourceNameValidator_Description_Default(t *testing.T) {
	regex.ListCavResourceNames = []regex.CavResourceName{
		{
			Key:         "test",
			Description: "",
			RegexString: ".*",
			RegexP:      regexp.MustCompile(".*"),
		},
	}
	validator := ResourceName("test")
	desc := validator.Description(context.Background())
	assert.Equal(t, "Ensures the string is a valid resource name.", desc)
}

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

	"github.com/orange-cloudavenue/common-go/regex"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = ResourceNameValidator{}

type ResourceNameValidator struct {
	ressourceType string
}

func ResourceName(ressourceType string) ResourceNameValidator {
	return ResourceNameValidator{ressourceType: ressourceType}
}

func (v ResourceNameValidator) findRegexResourceName() *regex.CavResourceName {
	for _, r := range regex.ListCavResourceNames {
		if r.Key == v.ressourceType {
			return &r
		}
	}
	return nil
}

func (v ResourceNameValidator) Description(_ context.Context) string {
	desc := v.findRegexResourceName().Description
	if desc == "" {
		return "Ensures the string is a valid resource name."
	}
	return desc
}

func (v ResourceNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ResourceNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	r := v.findRegexResourceName()
	if r == nil {
		return
	}

	str := req.ConfigValue.ValueString()
	if !r.RegexP.MatchString(str) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Resource Name",
			"The provided value is invalid.\n\n"+r.Description+"\n\n"+r.RegexString,
		)
	}
}

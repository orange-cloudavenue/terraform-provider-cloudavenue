/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package errors provides shared error-handling helpers for the provider:
// unified not-found detection across backends and harmonized diagnostic summaries.
package errors

import (
	"errors"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	caverrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

// IsNotFound reports whether err (or any wrapped error) represents a
// not-found condition from any backend used by the provider:
// the Cloud Avenue SDK (caverrors.ErrNotFound, commoncloudavenue.IsNotFound)
// or govcd (govcd.ErrorEntityNotFound / ContainsNotFound / IsNotFound).
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return caverrors.IsNotFound(err) ||
		commoncloudavenue.IsNotFound(err) ||
		errors.Is(err, govcd.ErrorEntityNotFound) ||
		govcd.IsNotFound(err) ||
		govcd.ContainsNotFound(err)
}

// Action verbs used in harmonized diagnostic summaries ("Error <verb> <resource>").
const (
	ActionCreate = "creating"
	ActionRead   = "reading"
	ActionUpdate = "updating"
	ActionDelete = "deleting"
)

// AddError appends a diagnostic with the harmonized summary
// "Error <verb> <resource>" (e.g. "Error reading vApp") and the raw error as detail.
func AddError(diags *diag.Diagnostics, verb, resource string, err error) {
	diags.AddError(fmt.Sprintf("Error %s %s", verb, resource), err.Error())
}

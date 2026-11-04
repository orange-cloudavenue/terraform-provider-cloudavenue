// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// Package errors provides shared error-handling helpers for the provider:
// unified not-found detection across backends and harmonized diagnostic summaries.
package errors

import (
	"errors"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	v2errors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	caverrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

// IsNotFound reports whether err (or any wrapped error) represents a
// not-found condition from any backend used by the provider:
//   - CloudAvenue SDK v1 (caverrors.IsNotFound / caverrors.ErrNotFound,
//     commoncloudavenue.IsNotFound)
//   - CloudAvenue SDK v2, either an *v2errors.APIError classified as not-found
//     by the SDK itself (HTTP 404, or wrapping v2errors.ErrNotFound), or the
//     bare v2errors.ErrNotFound sentinel
//   - govcd (govcd.ErrorEntityNotFound / govcd.IsNotFound / govcd.ContainsNotFound)
//
// This is the single not-found entry point for the provider: call sites must
// keep calling it. Matching a v2 sentinel directly at a v1 or govcd call site
// silently never matches.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return caverrors.IsNotFound(err) ||
		commoncloudavenue.IsNotFound(err) ||
		isSDKV2NotFound(err) ||
		errors.Is(err, govcd.ErrorEntityNotFound) ||
		govcd.IsNotFound(err) ||
		govcd.ContainsNotFound(err)
}

// isSDKV2NotFound detects not-found from the CloudAvenue SDK v2 error model.
// It prefers *v2errors.APIError.IsNotFound, the SDK's own predicate, because a
// lookup by id-or-name can return an APIError with StatusCode 404 and Err unset:
// such an error is not reachable through errors.Is(err, v2errors.ErrNotFound)
// because APIError.Unwrap() then returns nil, while the SDK still classifies
// it as not-found. The sentinel term is kept for errors that carry no
// APIError at all.
func isSDKV2NotFound(err error) bool {
	var apiErr *v2errors.APIError
	if errors.As(err, &apiErr) && apiErr.IsNotFound() {
		return true
	}
	return errors.Is(err, v2errors.ErrNotFound)
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

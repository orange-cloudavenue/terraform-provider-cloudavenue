/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package metrics

import (
	"time"
)

// version that can be overwritten by a release process.
var version = "dev"

// token can be overwritten by a release process.
var token = "dev"

// target can be overwritten by a release process.
var target = "https://localhost"

// GlobalExecutionID is the execution ID of the current Terraform run.
var GlobalExecutionID = ""

type (
	analyticRequest struct {
		*terraformRequest

		// OrganizationID is the id of the cloudavenue organization
		OrganizationID string `json:"organizationId"`
		// ResourceName is the name of the resource
		ResourceName string `json:"resourceName"`
		// Action is the action performed on the resource (create, update, delete, read or import)
		Action string `json:"action"`

		// ExecutionTime is the time in ms to execute the action
		ExecutionTime int64 `json:"executionTime"`

		// Data is the interface containing extra data
		Data map[string]any `json:"data,omitempty"`
	}

	terraformRequest struct {
		// TerraformExecutionID is the uniq id generated at the beginning of the terraform execution
		TerraformExecutionID string `json:"terraformExecutionId"`
		// ClientToken is the key used to identify the client
		ClientToken string `json:"clientToken"`
		// ClientVersion is the version of the provider
		ClientVersion string `json:"version"`
	}
)

func New(resourceName, organizationID string, action Action) func() {
	if everyThingIsOK() {
		start := time.Now()
		return func() {
			timeElapsed := time.Since(start)
			send(
				analyticRequest{
					terraformRequest: &terraformRequest{
						TerraformExecutionID: GlobalExecutionID,
						ClientVersion:        "terraform-cloudavenue/" + version,
						ClientToken:          token,
					},
					ResourceName:   resourceName,
					OrganizationID: organizationID,
					Action:         action.String(),
					ExecutionTime:  timeElapsed.Milliseconds(),
				})
		}
	}
	return func() {}
}

// everyThingIsOK Check if all variables are set.
func everyThingIsOK() bool {
	if version == "" || token == "" || target == "" {
		return false
	}
	return true
}

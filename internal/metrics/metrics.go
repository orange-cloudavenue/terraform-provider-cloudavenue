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

	tat "github.com/FrangipaneTeam/terraform-analytic-tool/api"
)

// version that can be overwritten by a release process.
var version = "dev"

// token can be overwritten by a release process.
var token = "dev"

// target can be overwritten by a release process.
var target = "https://localhost"

// GlobalExecutionID is the execution ID of the current Terraform run.
var GlobalExecutionID = ""

func New(resourceName, organizationID string, action Action) func() {
	if everyThingIsOK() {
		start := time.Now()
		return func() {
			timeElapsed := time.Since(start)
			send(
				tat.AnalyticRequest{
					TerraformRequest: &tat.TerraformRequest{
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

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

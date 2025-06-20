/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ELBPoliciesHTTPSecurityDataSource{}

const (
	ELBPoliciesHTTPSecurityDataSourceName = testsacc.ResourceName("data.cloudavenue_elb_policies_http_security")
)

type ELBPoliciesHTTPSecurityDataSource struct{}

func NewELBPoliciesHTTPSecurityDataSourceTest() testsacc.TestACC {
	return &ELBPoliciesHTTPSecurityDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBPoliciesHTTPSecurityDataSource) GetResourceName() string {
	return ELBPoliciesHTTPSecurityDataSourceName.String()
}

func (r *ELBPoliciesHTTPSecurityDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[ELBPoliciesHTTPSecurityResourceName]().GetDefaultConfig)
	return
}

func (r *ELBPoliciesHTTPSecurityDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_elb_policies_http_security" "example" {
						virtual_service_id = cloudavenue_elb_virtual_service.example.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[ELBPoliciesHTTPSecurityResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccELBPoliciesHTTPSecurityDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBPoliciesHTTPSecurityDataSource{}),
	})
}

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

var _ testsacc.TestACC = &ELBPoliciesHTTPResponseDataSource{}

const (
	ELBPoliciesHTTPResponseDataSourceName = testsacc.ResourceName("data.cloudavenue_elb_policies_http_response")
)

type ELBPoliciesHTTPResponseDataSource struct{}

func NewELBPoliciesHTTPResponseDataSourceTest() testsacc.TestACC {
	return &ELBPoliciesHTTPResponseDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *ELBPoliciesHTTPResponseDataSource) GetResourceName() string {
	return ELBPoliciesHTTPResponseDataSourceName.String()
}

func (r *ELBPoliciesHTTPResponseDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	// Add dependencies config to the resource
	resp.Append(GetResourceConfig()[ELBPoliciesHTTPResponseResourceName]().GetDefaultConfig)
	return
}

func (r *ELBPoliciesHTTPResponseDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, _ string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_elb_policies_http_response" "example" {
						virtual_service_id = cloudavenue_elb_virtual_service.example.id
					}`,
					// Here use resource config test to test the data source
					// the field example is the name of the test
					Checks: GetResourceConfig()[ELBPoliciesHTTPResponseResourceName]().GetDefaultChecks(),
				},
			}
		},
	}
}

func TestAccELBPoliciesHTTPResponseDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ELBPoliciesHTTPResponseDataSource{}),
	})
}

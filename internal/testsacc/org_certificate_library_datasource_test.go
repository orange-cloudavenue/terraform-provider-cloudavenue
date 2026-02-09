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

package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &OrgCertificateLibraryDatasource{}

const (
	OrgCertificateLibraryDatasourceName = testsacc.ResourceName("data.cloudavenue_org_certificate_library")
)

type OrgCertificateLibraryDatasource struct{}

func NewOrgCertificateLibraryDatasourceTest() testsacc.TestACC {
	return &OrgCertificateLibraryDatasource{}
}

// GetResourceName returns the name of the resource.
func (r *OrgCertificateLibraryDatasource) GetResourceName() string {
	return OrgCertificateLibraryDatasourceName.String()
}

func (r *OrgCertificateLibraryDatasource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[ORGCertificateLibraryResourceName]().GetDefaultConfig)
	return resp
}

func (r *OrgCertificateLibraryDatasource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_org_certificate_library" "example" {
						name = cloudavenue_org_certificate_library.example.name
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.CertificateLibraryItem)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrSet(resourceName, "certificate"),
					},
				},
			}
		},
		"example_id": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_org_certificate_library" "example_id" {
						id = cloudavenue_org_certificate_library.example.id
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.CertificateLibraryItem)),
						resource.TestCheckResourceAttrSet(resourceName, "name"),
						resource.TestCheckResourceAttrSet(resourceName, "certificate"),
					},
				},
			}
		},
	}
}

func TestAccOrgCertificateLibraryDatasource(t *testing.T) {
	cleanup := orgCertificateLibraryResourcePreCheck()
	defer cleanup()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&OrgCertificateLibraryDatasource{}),
	})
}

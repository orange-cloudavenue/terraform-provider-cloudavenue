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
	"os"
	"testing"

	"github.com/madflojo/testcerts"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &ORGCertificateLibraryResource{}

const (
	ORGCertificateLibraryResourceName = testsacc.ResourceName("cloudavenue_org_certificate_library")
)

type ORGCertificateLibraryResource struct{}

func NewORGCertificateLibraryResourceTest() testsacc.TestACC {
	return &ORGCertificateLibraryResource{}
}

// GetResourceName returns the name of the resource.
func (r *ORGCertificateLibraryResource) GetResourceName() string {
	return ORGCertificateLibraryResourceName.String()
}

func (r *ORGCertificateLibraryResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *ORGCertificateLibraryResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.CertificateLibraryItem)),
					resource.TestCheckResourceAttrSet(resourceName, "certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_org_certificate_library" "example" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						certificate = chomp(file("/tmp/cert.pem"))
						private_key = chomp(file("/tmp/key.pem"))
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckNoResourceAttr(resourceName, "passphrase"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_org_certificate_library" "example" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							certificate = chomp(file("/tmp/cert.pem"))
							private_key = chomp(file("/tmp/key.pem"))
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckNoResourceAttr(resourceName, "passphrase"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_org_certificate_library" "example" {
							name = {{ generate . "name" }}
							description = {{ generate . "description" }}
							certificate = chomp(file("/tmp/cert.pem"))
							private_key = chomp(file("/tmp/key.pem"))
							passphrase = "password"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrSet(resourceName, "passphrase"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder:    []string{"id"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"passphrase", "private_key"},
					},
					{
						ImportStateIDBuilder:    []string{"name"},
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"passphrase", "private_key"},
					},
				},
			}
		},
	}
}

const (
	orgCertificateLibraryCertFile = "/tmp/cert.pem"
	orgCertificateLibraryKeyFile  = "/tmp/key.pem"
)

func orgCertificateLibraryResourcePreCheck() (cleanup func()) {
	if err := testcerts.GenerateCertsToFile(
		orgCertificateLibraryCertFile,
		orgCertificateLibraryKeyFile,
	); err != nil {
		panic(err)
	}

	return func() {
		os.Remove(orgCertificateLibraryCertFile)
		os.Remove(orgCertificateLibraryKeyFile)
	}
}

func TestAccORGCertificateLibraryResource(t *testing.T) {
	cleanup := orgCertificateLibraryResourcePreCheck()
	defer cleanup()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&ORGCertificateLibraryResource{}),
	})
}

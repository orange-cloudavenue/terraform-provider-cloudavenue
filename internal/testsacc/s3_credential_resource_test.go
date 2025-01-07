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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3CredentialResource{}

const (
	S3CredentialResourceName = testsacc.ResourceName("cloudavenue_s3_credential")
)

type S3CredentialResource struct{}

func NewS3CredentialResourceTest() testsacc.TestACC {
	return &S3CredentialResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3CredentialResource) GetResourceName() string {
	return S3CredentialResourceName.String()
}

func (r *S3CredentialResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *S3CredentialResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`[a-z.-]+-\w{4}$`)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_credential" "example" {
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "username"),
						resource.TestCheckResourceAttr(resourceName, "file_name", "token.json"),
						resource.TestCheckResourceAttr(resourceName, "save_in_file", "false"),
						resource.TestCheckResourceAttr(resourceName, "print_token", "false"),
						resource.TestCheckResourceAttr(resourceName, "save_in_tfstate", "false"),
						resource.TestCheckResourceAttrSet(resourceName, "access_key"),
						resource.TestCheckNoResourceAttr(resourceName, "secret_key"),
						testCheckFileNotExists("token.json"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{},
				// ! Imports testing
				Imports: []testsacc.TFImport{},
				Destroy: true,
			}
		},
		"example_save_in_file": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`[a-z.-]+-\w{4}$`)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_credential" "example_save_in_file" {
					  save_in_file	  = true
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "username"),
						resource.TestCheckResourceAttr(resourceName, "file_name", "token.json"),
						resource.TestCheckResourceAttr(resourceName, "save_in_file", "true"),
						resource.TestCheckResourceAttr(resourceName, "print_token", "false"),
						resource.TestCheckResourceAttr(resourceName, "save_in_tfstate", "false"),
						resource.TestCheckResourceAttrSet(resourceName, "access_key"),
						resource.TestCheckNoResourceAttr(resourceName, "secret_key"),
						testCheckFileExists("token.json"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{},
				// ! Imports testing
				Imports: []testsacc.TFImport{},
				Destroy: true,
			}
		},
		"example_save_in_custom_file": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`[a-z.-]+-\w{4}$`)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_credential" "example_save_in_custom_file" {
					  save_in_file	  = true
					  file_name		  = "custom_token.json"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "username"),
						resource.TestCheckResourceAttr(resourceName, "file_name", "custom_token.json"),
						resource.TestCheckResourceAttr(resourceName, "save_in_file", "true"),
						resource.TestCheckResourceAttr(resourceName, "print_token", "false"),
						resource.TestCheckResourceAttr(resourceName, "save_in_tfstate", "false"),
						resource.TestCheckResourceAttrSet(resourceName, "access_key"),
						resource.TestCheckNoResourceAttr(resourceName, "secret_key"),
						testCheckFileExists("custom_token.json"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{},
				// ! Imports testing
				Imports: []testsacc.TFImport{},
				Destroy: true,
			}
		},
		"example_print_token": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`[a-z.-]+-\w{4}$`)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_credential" "example_print_token" {
					  print_token	  = true
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "username"),
						resource.TestCheckResourceAttr(resourceName, "file_name", "token.json"),
						resource.TestCheckResourceAttr(resourceName, "save_in_file", "false"),
						resource.TestCheckResourceAttr(resourceName, "print_token", "true"),
						resource.TestCheckResourceAttr(resourceName, "save_in_tfstate", "false"),
						resource.TestCheckResourceAttrSet(resourceName, "access_key"),
						resource.TestCheckNoResourceAttr(resourceName, "secret_key"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{},
				// ! Imports testing
				Imports: []testsacc.TFImport{},
				Destroy: true,
			}
		},
		"example_save_in_tfstate": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`[a-z.-]+-\w{4}$`)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_credential" "example_save_in_tfstate" {
					  save_in_tfstate	  = true
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "username"),
						resource.TestCheckResourceAttr(resourceName, "file_name", "token.json"),
						resource.TestCheckResourceAttr(resourceName, "save_in_file", "false"),
						resource.TestCheckResourceAttr(resourceName, "print_token", "false"),
						resource.TestCheckResourceAttr(resourceName, "save_in_tfstate", "true"),
						resource.TestCheckResourceAttrSet(resourceName, "access_key"),
						resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{},
				// ! Imports testing
				Imports: []testsacc.TFImport{},
				Destroy: true,
			}
		},
	}
}

func TestAccS3CredentialResource(t *testing.T) {
	t.Cleanup(deleteFile("token.json", t))
	t.Cleanup(deleteFile("custom_token.json", t))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3CredentialResource{}),
	})
}

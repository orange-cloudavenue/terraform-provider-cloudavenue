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

// Package client is the main client for the CloudAvenue provider.
package client

import (
	"errors"
	"os"
	"testing"

	clientca "github.com/orange-cloudavenue/cloudavenue-sdk-go"
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	caverror "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

func TestCloudAvenueClient(t *testing.T) {
	listOfEnvSet := []string{
		"TEST_CLOUDAVENUE_ORG",
		"TEST_CLOUDAVENUE_USERNAME",
		"TEST_CLOUDAVENUE_PASSWORD",
		"TEST_CLOUDAVENUE_VDC",
	}

	listOfEnvUnset := []string{
		"CLOUDAVENUE_ORG",
		"CLOUDAVENUE_USERNAME",
		"CLOUDAVENUE_PASSWORD",
		"CLOUDAVENUE_VDC",
	}

	for _, env := range listOfEnvSet {
		if os.Getenv(env) == "" {
			t.Fatalf("the environment variable %s is not set", env)
		}
	}

	for _, env := range listOfEnvUnset {
		if os.Getenv(env) != "" {
			t.Fatalf("the environment variable %s is set", env)
		}
	}

	t.Run("NewClient", func(t *testing.T) {
		tests := []struct {
			name        string
			opts        *clientca.ClientOpts
			wantErr     bool
			expectedErr error
		}{
			{
				name: "Bad Org",
				opts: &clientca.ClientOpts{
					CloudAvenue: &clientcloudavenue.Opts{
						Org:      "bad",
						Username: "user",
						Password: "password",
					},
				},
				wantErr:     true,
				expectedErr: caverror.ErrInvalidFormat,
			},
			{
				name: "Valid Org - Username not set",
				opts: &clientca.ClientOpts{
					CloudAvenue: &clientcloudavenue.Opts{
						Org:      "bad",
						Password: "password",
					},
				},
				wantErr:     true,
				expectedErr: caverror.ErrEmpty,
			},
			{
				name: "Valid Org - Password not set",
				opts: &clientca.ClientOpts{
					CloudAvenue: &clientcloudavenue.Opts{
						Org:      "bad",
						Username: "user",
					},
				},
				wantErr:     true,
				expectedErr: caverror.ErrEmpty,
			},
			{
				name: "Valid Org - Bad credential",
				opts: &clientca.ClientOpts{
					CloudAvenue: &clientcloudavenue.Opts{
						Org:      "cav01ev01ocb0001234",
						Username: "user",
						Password: "password",
					},
				},
				wantErr: true,
				// expectedErr: errors.New("ErrorMessage:Unauthorized"), - TODO: Catch error in SDK
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				c := CloudAvenue{
					CAVSDKOpts: tt.opts,
				}

				_, err := c.New()
				if tt.wantErr {
					if err == nil {
						t.Errorf("expected error: %v", tt.expectedErr)
						return
					}

					if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
						t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
					}

					return
				}
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			})
		}
	})

	t.Run("NewClientFromEnv", func(t *testing.T) {
		tests := []struct {
			name        string
			opts        *clientca.ClientOpts
			wantErr     bool
			expectedErr error
		}{
			{
				name:    "Valid Org - Credential provided by env",
				opts:    nil,
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := CloudAvenue{}
				// Set environment variables for the test
				t.Setenv("CLOUDAVENUE_ORG", os.Getenv("TEST_CLOUDAVENUE_ORG"))
				t.Setenv("CLOUDAVENUE_USERNAME", os.Getenv("TEST_CLOUDAVENUE_USERNAME"))
				t.Setenv("CLOUDAVENUE_PASSWORD", os.Getenv("TEST_CLOUDAVENUE_PASSWORD"))
				t.Setenv("CLOUDAVENUE_VDC", os.Getenv("TEST_CLOUDAVENUE_VDC"))

				client, err := c.New()
				if tt.wantErr {
					if err == nil {
						t.Fatalf("expected error: %v", tt.expectedErr)
					}

					if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
						t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
					}

					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if !client.DefaultVDCExist() {
					t.Fatalf("expected default VDC to exist")
				}

				if client.GetDefaultVDC() != os.Getenv("TEST_CLOUDAVENUE_VDC") {
					t.Fatalf("expected default VDC to be %s, got %s", os.Getenv("TEST_CLOUDAVENUE_VDC"), client.GetDefaultVDC())
				}

				if client.GetURL() == "" {
					t.Fatalf("expected URL to be set")
				}

				if client.GetOrgName() != os.Getenv("TEST_CLOUDAVENUE_ORG") {
					t.Fatalf("expected organization to be %s, got %s", os.Getenv("TEST_CLOUDAVENUE_ORG"), client.GetOrgName())
				}

				if client.GetUserName() != os.Getenv("TEST_CLOUDAVENUE_USERNAME") {
					t.Fatalf("expected username to be %s, got %s", os.Getenv("TEST_CLOUDAVENUE_USERNAME"), client.GetUserName())
				}
			})
		}
	})
}

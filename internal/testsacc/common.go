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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	regexpTier0VRFName = `prvrf[0-9]{2}eocb[0-9]{7}allsp[0-9]{2}`
)

func testCheckFileExists(filename string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		filename = filepath.Clean(filename)
		_, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		return nil
	}
}

func testCheckFileNotExists(filename string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		filename = filepath.Clean(filename)
		// Check if file exists
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("file %s exists", filename)
		}

		return nil
	}
}

// This is a helper function that attempts to remove file creating during unit test.
func deleteFile(filename string, t *testing.T) func() {
	return func() {
		if _, err := os.Stat(filename); err == nil {
			err := os.Remove(filename)
			if err != nil {
				t.Errorf("Failed to delete file: %s", err)
			}
		}
	}
}

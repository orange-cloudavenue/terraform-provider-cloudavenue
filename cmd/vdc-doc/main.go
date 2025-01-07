/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// vdc doc is a tool to generate documentation for the vdc provider.

package main

import (
	"log"
	"os"
	"strings"

	rules "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi/rules"
)

const (
	vdcMarker     = "<!-- TABLE VDC ATTRIBUTES PARAMETERS -->"
	storageMarker = "<!-- TABLE STORAGE PROFILES ATTRIBUTES PARAMETERS -->"
)

func main() {
	// Read the content of the file into a string
	filePath := "docs/resources/vdc.md"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Default().Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	// Get the content before and after the markers
	before := strings.Split(string(content), vdcMarker)[0]
	after := strings.Split(string(content), storageMarker)[1]

	// Generate the documentation for the VDC attributes
	vdcAttributes := rules.GetRulesDetails()

	// Generate the documentation for the Storage Profiles attributes
	storageProfilesAttributes := rules.GetStorageProfilesDetails()

	// Generate the content of the file
	newContent := before + vdcMarker + "\n" + vdcAttributes + "\n" + storageMarker + "\n" + storageProfilesAttributes + "\n" + after

	// Write the content to the file
	err = os.WriteFile(filePath, []byte(newContent), 0o600)
	if err != nil {
		log.Default().Printf("Failed to write file: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

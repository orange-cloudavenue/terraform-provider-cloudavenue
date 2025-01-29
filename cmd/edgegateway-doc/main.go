/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// edgegateway doc is a tool to generate documentation for the edgegateway resource.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

const (
	bandwidthMarker = "<!-- TABLE BANDWIDTH VALUES -->"
)

func main() {
	// Read the content of the file into a string
	filePath := "docs/resources/edgegateway.md"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Default().Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	// Get the content before and after the markers
	before := strings.Split(string(content), bandwidthMarker)[0]
	after := strings.Split(string(content), bandwidthMarker)[1]

	// * Retrieve the rules for the edgegateway and construct a markdown table

	rules := []string{}

	for key, value := range v1.EdgeGatewayAllowedBandwidth {
		rules = append(rules, fmt.Sprintf("* `%s` %s\n", key, strings.Trim(strings.ReplaceAll(fmt.Sprint(value.T1AllowedBandwidth), " ", ", "), "[]")))
	}

	// Generate the content of the file
	newContent := before + bandwidthMarker + "\n" + strings.Join(rules, "") + "\n" + after

	// Write the content to the file
	err = os.WriteFile(filePath, []byte(newContent), 0o600)
	if err != nil {
		log.Default().Printf("Failed to write file: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

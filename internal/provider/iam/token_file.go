// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package iam

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// apiTokenFileType is the token_type marker understood by govcd-based
// token file readers (govcd.GetTokenFromFile / SetApiTokenFromFile).
const apiTokenFileType = "API Token"

// apiTokenFileUpdatedBy identifies this provider as the writer of a
// token file. The field is informational: readers ignore it.
const apiTokenFileUpdatedBy = "terraform-provider-cloudavenue" //nolint:gosec // writer marker, not a credential

// apiTokenFile mirrors the govcd ApiTokenRefresh JSON document written by
// govcd.SaveApiTokenToFile, so files produced here stay readable by any
// govcd-based tooling. Field names and JSON keys are dictated by that
// external format contract.
type apiTokenFile struct {
	AccessToken  string `json:"access_token,omitempty"`  //nolint:tagliatelle // JSON key fixed by govcd token file format
	TokenType    string `json:"token_type,omitempty"`    //nolint:tagliatelle // JSON key fixed by govcd token file format
	ExpiresIn    int    `json:"expires_in,omitempty"`    //nolint:tagliatelle // JSON key fixed by govcd token file format
	RefreshToken string `json:"refresh_token,omitempty"` //nolint:tagliatelle // JSON key fixed by govcd token file format
	UpdatedBy    string `json:"updated_by,omitempty"`    //nolint:tagliatelle // JSON key fixed by govcd token file format
	UpdatedOn    string `json:"updated_on,omitempty"`    //nolint:tagliatelle // JSON key fixed by govcd token file format
}

// saveAPITokenToFile writes the API token file used to hand the durable
// refresh token over to external tooling. Ported from
// govcd.SaveApiTokenToFile to keep the on-disk format byte-compatible.
func saveAPITokenToFile(filename, refreshToken string) error {
	content := apiTokenFile{
		TokenType:    apiTokenFileType,
		RefreshToken: refreshToken,
		UpdatedBy:    apiTokenFileUpdatedBy,
		UpdatedOn:    time.Now().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(content, " ", " ") //nolint:gosec // the refresh token is marshaled on purpose
	if err != nil {
		return fmt.Errorf("error marshalling token to JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0o600); err != nil {
		return fmt.Errorf("error writing token file: %w", err)
	}

	return nil
}

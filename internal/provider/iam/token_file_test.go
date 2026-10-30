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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v56 "github.com/vmware/go-vcloud-director/v2/types/v56"
)

// TestSaveAPITokenToFile checks the written document stays compatible with
// the govcd token file readers (govcd.GetTokenFromFile).
func TestSaveAPITokenToFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "token.json")

	require.NoError(t, saveAPITokenToFile(filename, "refresh-secret"))

	info, err := os.Stat(filename)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	raw, err := os.ReadFile(filename)
	require.NoError(t, err)

	var doc map[string]any
	require.NoError(t, json.Unmarshal(raw, &doc))
	assert.Equal(t, "API Token", doc["token_type"])
	assert.Equal(t, "refresh-secret", doc["refresh_token"])
	assert.Equal(t, apiTokenFileUpdatedBy, doc["updated_by"])
	assert.NotEmpty(t, doc["updated_on"])
	assert.NotContains(t, doc, "access_token")

	var asGovcd v56.ApiTokenRefresh
	require.NoError(t, json.Unmarshal(raw, &asGovcd))
	assert.Equal(t, "refresh-secret", asGovcd.RefreshToken)
	assert.Equal(t, "API Token", asGovcd.TokenType)
}

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

package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Send is a function to send an event
// with a given configuration and client.
// This not return an error because it's not critical.
func send(event analyticRequest) {
	// Serialize and pack event
	eventPkg, err := json.Marshal(event)
	if err != nil {
		return
	}

	// Context with 1 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Compose request
	req, err := http.NewRequestWithContext(ctx, "POST", target+"/api/v1/send", bytes.NewReader(eventPkg))
	if err != nil {
		return
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	// Add Header Authorization if token is set
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	// Send request
	res, err := http.DefaultClient.Do(req)
	// Check error
	if err != nil {
		return
	}
	res.Body.Close()
}

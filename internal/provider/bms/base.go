/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package bms

import "time"

const (
	categoryName = "bms"

	// defaultReadTimeout is the default timeout for read operations.
	defaultReadTimeout = 5 * time.Minute
	// defaultCreateTimeout is the default timeout for create operations.
	// defaultCreateTimeout = 5 * time.Minute
	// defaultUpdateTimeout is the default timeout for update operations.
	// defaultUpdateTimeout = 5 * time.Minute
	// defaultDeleteTimeout is the default timeout for delete operations.
	// defaultDeleteTimeout = 5 * time.Minute.
)

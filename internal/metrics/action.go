/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package metrics

type Action string

const (
	Create Action = "Create"
	Read   Action = "Read"
	Update Action = "Update"
	Delete Action = "Delete"
	Import Action = "Import"
)

// String returns the string representation of the action.
func (a Action) String() string {
	return string(a)
}

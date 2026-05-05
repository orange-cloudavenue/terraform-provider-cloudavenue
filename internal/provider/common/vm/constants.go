/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vm

// Attribute name constants used across the vm package to avoid goconst violations.
const (
	attrBusType        = "bus_type"
	attrBusNumber      = "bus_number"
	attrSizeInMB       = "size_in_mb"
	attrUnitNumber     = "unit_number"
	attrStorageProfile = "storage_profile"
	attrName           = "name"
)

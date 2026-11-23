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

package vapp

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

const (
	attrVDC               = "vdc"
	attrName              = "name"
	attrVappName          = attrName
	attrRuntimeLeaseInSec = "runtime_lease_in_sec"
	attrStorageLeaseInSec = "storage_lease_in_sec"
)

// vApp statuses, see govcd types.VAppStatuses.
const (
	resolved            = types.VAppStatusResolved
	suspended           = types.VAppStatusSuspended
	poweredON           = types.VAppStatusPoweredOn
	partiallyDeployed   = types.VAppStatusVAppPartiallyDeployed
	partiallyPoweredOff = types.VAppStatusPartiallyPoweredOff
	partiallySuspended  = types.VAppStatusPartiallySuspended
)

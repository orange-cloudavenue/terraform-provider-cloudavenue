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

package vdc

import (
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

// DiskExist checks if a disk exists in a VDC.
func (v VDC) DiskExist(diskName string) (bool, error) {
	existingDisk, err := v.QueryDisk(diskName)
	if err != nil {
		if strings.Contains(err.Error(), "found results ") {
			return false, nil
		}
	}
	return existingDisk != (govcd.DiskRecord{}), err
}

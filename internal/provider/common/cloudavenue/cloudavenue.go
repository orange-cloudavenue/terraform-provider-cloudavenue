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

package cloudavenue

import (
	"context"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

var cAMutexKV = mutex.NewKV()

const keyLock = "cloudavenue:customer:api:lock"

// Lock
// lock call to the cloudavenue customer API.
func Lock(ctx context.Context) {
	cAMutexKV.KvLock(ctx, keyLock)
}

// Unlock
// unlock call to the cloudavenue customer API.
func Unlock(ctx context.Context) {
	cAMutexKV.KvUnlock(ctx, keyLock)
}

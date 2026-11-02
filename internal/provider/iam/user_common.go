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
	"context"

	iam "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/iam/v1"
)

// getUserQuotas reads the deployed/stored VM quotas of a user.
//
// CloudAPI accepts quotas on user create/update but never returns them on
// the user resource, so they are read back from the legacy admin user
// record (GET /api/admin/user/{uuid}). This is the only remaining consumer
// of the deprecated admin XML API for user reads; every other field of the
// user schema comes from CloudAPI.
//
// TODO: replace the legacy admin read with the CloudAPI entityQuotas
// endpoints (VCD 10.5+, GET/PUT /cloudapi/1.0.0/users/{userUrn}/quotas).
// Blocked on two open questions: they expose effective quota-policy values
// (possibly inherited from the org, with -1 = unlimited) instead of the raw
// assignment echoed by the legacy record, and it must be confirmed that
// quota writes through the user update body and quota policies are backed
// by a single store. See https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/entity-quotas/.
func getUserQuotas(ctx context.Context, client *iam.Client, id string) (deployed, stored int, err error) {
	quotas, err := client.GetUserQuotas(ctx, iam.ParamsGetUserQuotas{ID: id})
	if err != nil {
		return 0, 0, err
	}
	return quotas.DeployedVMQuota, quotas.StoredVMQuota, nil
}

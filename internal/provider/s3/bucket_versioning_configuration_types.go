/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package s3

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type BucketVersioningConfigurationModel struct {
	Timeouts timeoutsR.Value        `tfsdk:"timeouts"`
	ID       supertypes.StringValue `tfsdk:"id"`
	Bucket   supertypes.StringValue `tfsdk:"bucket"`
	Status   supertypes.StringValue `tfsdk:"status"`
}

type BucketVersioningConfigurationDatasourceModel struct {
	Timeouts timeoutsD.Value        `tfsdk:"timeouts"`
	ID       supertypes.StringValue `tfsdk:"id"`
	Bucket   supertypes.StringValue `tfsdk:"bucket"`
	Status   supertypes.StringValue `tfsdk:"status"`
}

func (rm *BucketVersioningConfigurationModel) Copy() *BucketVersioningConfigurationModel {
	x := &BucketVersioningConfigurationModel{}
	utils.ModelCopy(rm, x)
	return x
}

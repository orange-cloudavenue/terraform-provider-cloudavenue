/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package alb

import (
	"errors"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var ErrPersistenceProfileIsEmpty = errors.New("persistence profile is empty")

type albPool interface {
	GetID() string
	GetName() string
	GetAlbPool() (*govcd.NsxtAlbPool, error)
}

func processMembers(poolMembers []govcdtypes.NsxtAlbPoolMember) (members []member) {
	for _, poolMember := range poolMembers {
		members = append(members, member{
			Enabled:   types.BoolValue(poolMember.Enabled),
			IPAddress: types.StringValue(poolMember.IpAddress),
			Port:      types.Int64Value(int64(poolMember.Port)),
			Ratio:     types.Int64Value(int64(*poolMember.Ratio)),
		})
	}
	return
}

func processHealthMonitors(poolHealthMonitors []govcdtypes.NsxtAlbPoolHealthMonitor) (healthMonitors []string) {
	for _, poolHealthMonitor := range poolHealthMonitors {
		healthMonitors = append(healthMonitors, poolHealthMonitor.Type)
	}

	return
}

func processPersistenceProfile(poolPersistenceProfile *govcdtypes.NsxtAlbPoolPersistenceProfile) persistenceProfile {
	if poolPersistenceProfile == nil {
		return persistenceProfile{}
	}

	return persistenceProfile{
		Type:  types.StringValue(poolPersistenceProfile.Type),
		Value: utils.StringValueOrNull(poolPersistenceProfile.Value),
	}
}

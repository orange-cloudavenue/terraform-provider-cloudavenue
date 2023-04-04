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

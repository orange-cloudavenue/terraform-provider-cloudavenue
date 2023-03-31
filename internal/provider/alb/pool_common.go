package alb

import (
	"errors"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ErrPersistenceProfileIsEmpty = errors.New("persistence profile is empty")

type albPool interface {
	GetID() string
	GetName() string
	GetAlbPool() (*govcd.NsxtAlbPool, error)
}

func processMembers(poolMembers []govcdtypes.NsxtAlbPoolMember) []member {
	members := []member{}
	if len(poolMembers) > 0 {
		for _, poolMember := range poolMembers {
			members = append(members, member{
				Enabled:   types.BoolValue(poolMember.Enabled),
				IPAddress: types.StringValue(poolMember.IpAddress),
				Port:      types.Int64Value(int64(poolMember.Port)),
				Ratio:     types.Int64Value(int64(*poolMember.Ratio)),
			})
		}
	}

	return members
}

func processHealthMonitors(poolHealthMonitors []govcdtypes.NsxtAlbPoolHealthMonitor) []string {
	var healtMonitors []string
	if len(poolHealthMonitors) > 0 {
		for _, poolHealthMonitor := range poolHealthMonitors {
			healtMonitors = append(healtMonitors, poolHealthMonitor.Type)
		}
	}

	return healtMonitors
}

func processPersistenceProfile(poolPersistenceProfile *govcdtypes.NsxtAlbPoolPersistenceProfile) persistenceProfile {
	p := persistenceProfile{}
	if poolPersistenceProfile != nil {
		p.Type = types.StringValue(poolPersistenceProfile.Type)
		if poolPersistenceProfile.Value == "" {
			p.Value = types.StringNull()
		} else {
			p.Value = types.StringValue(poolPersistenceProfile.Value)
		}
	}

	return p
}

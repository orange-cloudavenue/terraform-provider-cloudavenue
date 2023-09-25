package testsacc

import "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"

func GetDataSourceConfig() map[testsacc.ResourceName]func() resourceConfig {
	return map[testsacc.ResourceName]func() resourceConfig{
		Tier0VRFDataSourceName: NewResourceConfig(NewTier0VRFDataSourceTest()),
		VDCDataSourceName:      NewResourceConfig(NewVDCDataSourceTest()),
	}
}

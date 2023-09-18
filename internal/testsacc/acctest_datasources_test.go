package testsacc

func GetDataSourceConfig() map[ResourceName]func() resourceConfig {
	return map[ResourceName]func() resourceConfig{
		Tier0VRFACLDataSourceName: NewResourceConfig(NewTier0VRFACLDataSourceTest()),
	}
}

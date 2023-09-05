package testsacc

func GetResourceConfig() map[ResourceName]func() resourceConfig {
	return map[ResourceName]func() resourceConfig{
		CatalogResourceName:    NewResourceConfig(NewCatalogResourceTest()),
		CatalogACLResourceName: NewResourceConfig(NewCatalogACLResourceTest()),
	}
}

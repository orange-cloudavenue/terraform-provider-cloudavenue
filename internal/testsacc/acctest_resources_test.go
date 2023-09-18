package testsacc

func GetResourceConfig() map[ResourceName]func() resourceConfig {
	return map[ResourceName]func() resourceConfig{
		// * Catalog
		CatalogResourceName:    NewResourceConfig(NewCatalogResourceTest()),
		CatalogACLResourceName: NewResourceConfig(NewCatalogACLResourceTest()),

		// * VDC
		VDCResourceName: NewResourceConfig(NewVDCResourceTest()),
	}
}

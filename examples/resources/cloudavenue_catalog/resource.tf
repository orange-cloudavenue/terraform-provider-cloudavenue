
resource "cloudavenue_catalog" "example" {
	catalog_name     = "test-catalog"
	description      = "catalog for ISO"
	delete_recursive = true
	delete_force     = true
}
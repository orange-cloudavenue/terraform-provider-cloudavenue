data "cloudavenue_catalogs" "test" {}

data "cloudavenue_catalog_media" "example" {
  catalog_name = data.cloudavenue_catalogs.test.catalogs["catalog-example"].catalog_name
  name         = "debian-9.9.0-amd64-netinst.iso"
}
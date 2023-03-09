data "cloudavenue_catalogs" "test" {}

data "cloudavenue_catalog_medias" "example" {
  catalog_name = data.cloudavenue_catalogs.test.catalogs["catalog-example"].catalog_name
}
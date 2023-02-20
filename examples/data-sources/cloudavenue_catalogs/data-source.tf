data "cloudavenue_catalogs" "test" {}

output "catalogs_Name" {
  value = data.cloudavenue_catalogs.test.catalogs_name
}

output "catalog_Orange-Linux" {
  value = data.cloudavenue_catalogs.test.catalogs["Orange-Linux"]
}
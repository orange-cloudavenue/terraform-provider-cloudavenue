data "cloudavenue_vdcgs" "example" {
  filter_by_name = "example" // strict match
}

data "cloudavenue_vdcgs" "example_wildcard" {
  filter_by_name = "production-*" // Match with wildcard
}

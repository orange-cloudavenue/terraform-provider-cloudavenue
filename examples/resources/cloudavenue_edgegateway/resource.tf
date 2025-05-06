resource "cloudavenue_edgegateway" "example" {
  owner_name = cloudavenue_vdc.example.name
}

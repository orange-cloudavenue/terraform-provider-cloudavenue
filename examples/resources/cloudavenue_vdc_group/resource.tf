resource "cloudavenue_vdc_group" "example" {
  name = "example"
  vdc_ids = [
    cloudavenue_vdc.example-without-vdc-group.id,
  ]
}

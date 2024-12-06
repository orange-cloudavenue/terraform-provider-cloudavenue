resource "cloudavenue_vdcg" "example" {
  name = "example"
  vdc_ids = [
    cloudavenue_vdc.example.id,
  ]
}

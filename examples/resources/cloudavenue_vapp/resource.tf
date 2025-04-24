resource "cloudavenue_vapp" "example" {
  name        = "example"
  vdc         = cloudavenue_vdc.example.name
  description = "This is an example vApp"

  lease = {
    runtime_lease_in_sec = 3600
    storage_lease_in_sec = 3600
  }

  guest_properties = {
    "key" = "Value"
  }
}

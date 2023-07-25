resource "cloudavenue_vapp" "example" {
  name        = "MyVapp"
  description = "This is an example vApp"

  lease = {
    runtime_lease_in_sec = 3600
    storage_lease_in_sec = 3600
  }

  guest_properties = {
    "key" = "Value"
  }
}
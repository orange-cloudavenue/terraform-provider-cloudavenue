resource "cloudavenue_vapp" "example" {
  name        = "MyVapp"
  description = "This is an example vApp"
  power_on    = true

  lease = {
    runtime_lease_in_sec = 3600
    storage_lease_in_sec = 7200
  }

  guest_properties = {
    "key" = "Value"
  }
}
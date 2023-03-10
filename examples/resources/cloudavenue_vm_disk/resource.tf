resource "cloudavenue_vapp" "example" {
  name        = "vapp_example"
  description = "This is a example vapp"
}

resource "cloudavenue_vm_disk" "example-detachable" {
  vapp_id       = cloudavenue_vapp.example.id
  name          = "disk-example"
  bus_type      = "SATA"
  size_in_mb    = 2048
  is_detachable = true
}

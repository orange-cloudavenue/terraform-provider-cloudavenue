resource "cloudavenue_vm_disk" "example" {
  vapp_name       = "vapp_test3"
  vm_name         = "TestRomain"
  allow_vm_reboot = true
  internal_disk = {
    bus_type    = "sata"
    size_in_mb  = "500"
    bus_number  = 0
    unit_number = 1
  }
}
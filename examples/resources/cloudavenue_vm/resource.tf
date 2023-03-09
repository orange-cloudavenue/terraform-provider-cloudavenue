data "cloudavenue_catalog_vapp_template" "example" {
  catalog_name = "Orange-Linux"
  vapp_name    = "debian_10_X64"
}

resource "cloudavenue_vapp" "example" {
  vapp_name   = "vapp_example"
  description = "This is a example vapp"
}

resource "cloudavenue_vm" "example" {
  vm_name          = "example-vm"
  description      = "This is a example vm"
  accept_all_eulas = true
  vapp_name        = cloudavenue_vapp.example.vapp_name
  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
}
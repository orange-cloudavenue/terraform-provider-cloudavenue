data "cloudavenue_catalog_vapp_template" "example" {
  catalog_name  = "Orange-Linux"
  template_name = "debian_10_X64"
}

resource "cloudavenue_vapp" "example" {
  name        = "vapp_example"
  description = "This is a example vapp"
}

resource "cloudavenue_vm" "example" {
  name             = "example-vm"
  description      = "This is a example vm"
  accept_all_eulas = true
  vapp_name        = cloudavenue_vapp.example.name
  vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
}
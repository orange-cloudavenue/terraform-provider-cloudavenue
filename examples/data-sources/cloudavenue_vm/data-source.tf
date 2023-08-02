data "cloudavenue_vapp" "example" {
  name = "example-vapp"
}

data "cloudavenue_vm" "example" {
  name      = "example-vm"
  vapp_name = data.cloudavenue_vapp.example.name
}

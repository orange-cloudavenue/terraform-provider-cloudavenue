data "cloudavenue_vm_affinity_rule" "example" {
  id = cloudavenue_vm_affinity_rule.example.id
}

resource "cloudavenue_vm_affinity_rule" "example" {
  name     = "example-affinity-rule"
  polarity = "Affinity"

  vm_ids = [
    cloudavenue_vm.example.id,
    cloudavenue_vm.example2.id,
  ]
}

resource "cloudavenue_vm" "example" {
  name      = "example-vm"
  vapp_name = cloudavenue_vapp.example.name
  deploy_os = {
    vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
  }
  settings = {
    customization = {
      auto_generate_password = true
    }
  }
  resource = {
  }

  state = {
  }
}

resource "cloudavenue_vm" "example2" {
  name      = "example-vm2"
  vapp_name = cloudavenue_vapp.example.name
  deploy_os = {
    vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
  }
  settings = {
    customization = {
      auto_generate_password = true
    }
  }
  resource = {
  }

  state = {
  }
}

data "cloudavenue_catalog_vapp_template" "example" {
  catalog_name  = "Orange-Linux"
  template_name = "debian_10_X64"
}

resource "cloudavenue_vapp" "example" {
  name        = "vapp_example"
  description = "This is a example vapp"
}
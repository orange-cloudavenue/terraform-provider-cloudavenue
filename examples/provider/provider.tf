provider "cloudavenue" {
  org      = var.org
  user     = var.user
  password = var.password
  # VMware/VCD endpoint.
  url = var.url
  # Supported Cerberus endpoint override, including internal URLs when required.
  core_api = var.core_api
}

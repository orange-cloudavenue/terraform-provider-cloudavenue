provider "cloudavenue" {
  org      = var.org
  user     = var.user
  password = var.password
  # VMware/VCD endpoint.
  url = var.url
  # Override the Cloud Avenue API endpoint URL (authentication and backend). Useful when accessing through a private or internal network.
  core_api = var.core_api
}

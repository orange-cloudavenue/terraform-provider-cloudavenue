data "cloudavenue_edgegateway" "example" {
  name = "tn01e02ocb0006205spt101"
}

resource "cloudavenue_alb_pool" "example" {
  edge_gateway_id = data.cloudavenue_edgegateway.example.id
  name            = "albpool-name"
  persistence_profile = {
    type = "CLIENT_IP"
  }
  members = [
    {
      ip_address = "192.168.99.11"
      port       = "80"
    },
    {
      ip_address = "192.168.10.2"
      port       = "80"
    },
    {
      ip_address = "192.168.1.3"
      port       = "80"
    }
  ]
  health_monitors = ["TCP"]
}
resource "cloudavenue_alb_virtual_service" "example" {
  name            = "albvs-name"
  description     = "description"
  edge_gateway_id = data.cloudavenue_edgegateway.example.id
  enabled         = true
  pool_id         = cloudavenue_alb_pool.example.id
  virtual_ip      = "192.168.10.10"
  certificate_id  = "urn:vcloud:certificateLibraryItem:f9caac3a-2555-477e-ae58-0740687d4daf"
  service_type    = "HTTPS"
  service_ports = [
    {
      port_start = 443
      port_type  = "TCP_PROXY"
      port_ssl   = true
    },
    {
      port_start = 8080
      port_type  = "TCP_PROXY"
      port_ssl   = true
    },
    {
      port_start = 8088
      port_type  = "TCP_PROXY"
      port_ssl   = true
    }
  ]
}
data "cloudavenue_edgegateway" "example" {
  name = "tn01e02ocb0006205spt101"
}

data "cloudavenue_alb_pool" "example" {
  name = "albpool-name"
}
resource "cloudavenue_alb_virtual_service" "example" {
  name            = "albvs-name"
  description     = "description"
  edge_gateway_id = data.cloudavenue_edgegateway.example.id
  enabled         = true
  pool_id         = cloudavenue_alb_pool.example.id
  virtual_ip      = "192.168.10.10"
  service_type    = "HTTP"
  service_ports = [
    {
      port_start = 80
    }
  ]
}
resource "cloudavenue_elb_virtual_service" "example" {
  name    = "example"
  enabled = true

  virtual_ip = "192.168.0.1"

  pool_id         = cloudavenue_elb_pool.example.id
  edge_gateway_id = cloudavenue_edgegateway.example.id

  service_type = "HTTP"
  service_ports = [
    {
      start = 80
    }
  ]
}

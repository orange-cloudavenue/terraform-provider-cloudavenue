resource "cloudavenue_edgegateway_static_route" "example" {
  name        = "example"
  description = "example description"

  edge_gateway_id = cloudavenue_edgegateway.example.id

  network_cidr = "192.168.2.0/24"
  next_hops = [
    {
      ip_address = "192.168.2.254"
    },
    {
      ip_address     = "192.168.2.253"
      admin_distance = 2
    }
  ]
}

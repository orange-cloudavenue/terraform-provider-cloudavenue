resource "cloudavenue_edgegateway_static_route" "example" {
  edge_gateway_name = "myEdgeName"
  name              = "example"
  description       = "example description"
  network_cidr      = "192.168.2.0/24"
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

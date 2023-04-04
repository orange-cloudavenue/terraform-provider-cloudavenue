resource "cloudavenue_network_routed" "example" {
  name            = "ExampleNetworkRouted"
  gateway         = "192.168.10.254"
  prefix_length   = 24
  edge_gateway_id = "urn:vcloud:gateway:dde5d31a-2f32-43ef-b3b3-127245958298"
}

data "cloudavenue_network_routed" "example" {
  name            = "ExampleNetworkRouted"
  edge_gateway_id = "urn:vcloud:gateway:dde5d31a-2f32-43ef-b3b3-127245958298"
}
data "cloudavenue_alb_virtual_service" "example" {
  edge_gateway_name = data.cloudavenue_edgegateway.example.name
  name              = "MyALBVirtualServiceName"
}

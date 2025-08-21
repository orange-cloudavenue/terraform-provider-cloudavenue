# Retrieving a CloudAvenue Tier-0 by Edge Gateway ID
data "cloudavenue_tier0_vrf" "example" {
  edge_gateway_id = cloudavenue_edgegateway.example.id
}

# Retrieving a CloudAvenue Tier-0 by Edge Gateway Name
data "cloudavenue_tier0_vrf" "example" {
  edge_gateway_name = cloudavenue_edgegateway.example.name
}

# Retrieving a CloudAvenue Tier-0 by Name
data "cloudavenue_tier0_vrf" "example" {
  name = "prvrf01eocb0001234allsp01"
}

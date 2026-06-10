# Simple profile with multiple App IDs (no sub-attributes)
resource "cloudavenue_edgegateway_network_context_profile" "example" {
  edge_gateway_name = cloudavenue_edgegateway.example.name
  name              = "my-custom-profile"
  description       = "Custom Layer 7 profile matching SSH and DNS traffic"

  attribute = [
    { app_id = "SSH", sub_attribute = [] },
    { app_id = "DNS", sub_attribute = [] },
  ]
}

# Profile with a single App ID and sub-attributes (TLS constraints)
resource "cloudavenue_edgegateway_network_context_profile" "ssl_strict" {
  edge_gateway_name = cloudavenue_edgegateway.example.name
  name              = "ssl-tls12-only"
  description       = "SSL restricted to TLS 1.2 and 1.3"

  attribute = [
    {
      app_id = "SSL"
      sub_attribute = [
        {
          type   = "TLS_VERSION"
          values = ["TLS_V12", "TLS_V13"]
        }
      ]
    }
  ]
}

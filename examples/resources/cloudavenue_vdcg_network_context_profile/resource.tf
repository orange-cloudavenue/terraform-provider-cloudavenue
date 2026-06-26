# Profile matching SSL traffic restricted to TLS 1.2 and 1.3
resource "cloudavenue_vdcg_network_context_profile" "example" {
  vdc_group_name = cloudavenue_vdcg.example.name
  name           = "ssl-tls12-only"
  description    = "Allow only TLS 1.2 and 1.3"

  app_id = {
    values = ["SSL"]
    sub_attributes = [
      {
        type   = "TLS_VERSION"
        values = ["TLS_V12", "TLS_V13"]
      }
    ]
  }
}


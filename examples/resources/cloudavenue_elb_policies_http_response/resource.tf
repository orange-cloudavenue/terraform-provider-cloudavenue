resource "cloudavenue_elb_policies_http_response" "example" {
  virtual_service_id = cloudavenue_elb_virtual_service.example.id
  policies = [
    // Policy 1
    {
      name = "example"

      // Define the criteria for the policy
      criteria = {
        client_ip = {
          criteria = "IS_IN"
          ip_addresses = [
            "12.13.14.15",
            "12.13.14.0/24",
            "12.13.14.1-12.13.14.15"
          ]
        }
        cookie = {
          criteria = "BEGINS_WITH"
          name     = "example"
          value    = "example"
        }
        http_methods = {
          criteria = "IS_IN"
          methods  = ["GET", "POST"]
        }
        path = {
          criteria = "CONTAINS"
          paths    = ["/example"]
        }
        protocol = "HTTPS"
        query = [
          "example=example"
        ]
        request_headers = [
          {
            criteria = "CONTAINS"
            name     = "X-EXAMPLE"
            values   = ["example"]
          },
          {
            criteria = "BEGINS_WITH"
            name     = "X-CUSTOM"
            values   = ["example"]
          }
        ]
        response_headers = [
          {
            criteria = "CONTAINS"
            name     = "X-RESPONSE"
            values   = ["example"]
          },
          {
            criteria = "BEGINS_WITH"
            name     = "X-RESPONSE-CUSTOM"
            values   = ["example"]
          }
        ]
        service_ports = {
          criteria = "IS_IN"
          ports    = ["443"]
        }
        location = {
          criteria = "BEGINS_WITH"
          values = [
            "example.com"
          ]
        }
        status_code = {
          criteria = "IS_IN"
          codes    = ["200", "302"]
        }
      }

      // Define the action to take when the criteria is met
      actions = {

        location_rewrite = {
          host       = "example.org"
          protocol   = "HTTPS"
          keep_query = true
          port       = 443
        }

        modify_headers = [
          {
            action = "ADD"
            name   = "X-FROM-OLD-DOMAIN"
            value  = "example.com"
          }
        ]
      }
    } // End policy 1
  ]   // End policies
}

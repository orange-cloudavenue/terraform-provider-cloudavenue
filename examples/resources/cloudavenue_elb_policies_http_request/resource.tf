resource "cloudavenue_elb_policies_http_request" "example" {
  virtual_service_id = cloudavenue_elb_virtual_service.example.id
  policies = [
    # Policy 1
    {
      name = "example"

      # Define the criteria for the policy
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
        service_ports = {
          criteria = "IS_IN"
          ports    = ["80"]
        }
      }

      // Define the action to take when the criteria is met
      actions = {
        modify_headers = [
          {
            action = "ADD"
            name   = "X-SECURE"
            value  = "example"
          },
          {
            action = "REMOVE"
            name   = "X-EXAMPLE"
          }
        ]
        rewrite_url = {
          host = "example.com"
          path = "/example"
        }
      }
    } // End policy 1
  ]   // End policies
}

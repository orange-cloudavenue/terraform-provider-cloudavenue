resource "cloudavenue_elb_policies_http_security" "example" {
  virtual_service_id = cloudavenue_elb_virtual_service.example.id
  policies = [
    # Policy 1
    {
      name = "example"

      # Define the criteria for the policy
      criteria = {
        // Client IP criteria
        // This criteria checks if the client IP address is in the specified list or range
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS"
        // The ip_addresses can be a single IP, a CIDR range, or a range of IPs
        client_ip = {
          criteria = "IS_IN"
          ip_addresses = [
            "12.13.14.15",
            "12.13.14.0/24",
            "12.13.14.1-12.13.14.15"
          ]
        }
        // Cookie criteria
        // This criteria checks if the specified cookie name and value match the criteria
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS"
        // The name and value are the cookie name and value to check
        cookie = {
          criteria = "BEGINS_WITH"
          name     = "example"
          value    = "example"
        }
        // HTTP methods criteria
        // This criteria checks if the HTTP method is in the specified list
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS"
        // The methods are the HTTP methods to check
        http_methods = {
          criteria = "IS_IN"
          methods  = ["GET", "POST"]
        }
        // Path criteria
        // This criteria checks if the request path matches the specified criteria
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS"
        // The paths are the request paths to check
        // The paths can be a single path or a list of paths
        path = {
          criteria = "CONTAINS"
          paths    = ["/example"]
        }
        // Protocol criteria
        // This criteria checks if the request protocol matches the specified criteria
        protocol = "HTTPS"
        // Query criteria
        // This criteria checks if the query string matches the specified criteria
        query = [
          "example=example"
        ]
        // Request headers criteria
        // This criteria checks if the request headers match the specified criteria
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS" etc...
        // The name and values are the request header name and values to check
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
        // Response headers criteria
        // This criteria checks if the response headers match the specified criteria
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS" etc...
        // The name and values are the response header name and values to check
        // The response headers can be a single header or a list of headers
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
        // Service ports criteria
        // This criteria checks if the service port is in the specified list
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS" etc...
        // The ports are the service ports to check
        service_ports = {
          criteria = "IS_IN"
          ports    = ["80", "443"]
        }
        // Location criteria
        // This criteria checks if the request location matches the specified criteria
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS" etc...
        // The values are the request location to check
        location = {
          criteria = "IS_IN"
          values = [
            "example.com",
            "example.org"
          ]
        }
        // Status code criteria
        // This criteria checks if the response status code is in the specified list
        // The criteria can be "IS_IN", "IS_NOT_IN", "BEGINS_WITH", "ENDS_WITH", or "CONTAINS" etc...
        // The codes are the response status codes to check
        status_code = {
          criteria = "IS_IN"
          codes    = ["200", "301", "302"]
        }
      }

      // Define the action to take when the criteria is met
      actions = {
        redirect_to_https = 443
      }
    } // End policy 1
  ]   // End policies
}

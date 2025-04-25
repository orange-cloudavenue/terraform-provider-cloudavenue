resource "cloudavenue_elb_policies_http_security" "example" {
  virtual_service_id = cloudavenue_elb_virtual_service.example.id
  policies = [
    # Policy 1
    {
      name = "example"

      # Define the criteria for the policy
      criteria = {
        // This example checks if the request is using HTTP
        // and redirects to HTTPS on port 8443
        protocol = "HTTP"
        service_ports = {
          criteria = "IS_IN"
          ports    = ["80"]
        }
      }

      // The action rate_limit can be set to limit the number of requests
      // Only one action can be set at a time
      // The action close_connection can be set to close the connection when the rate limit is reached.
      actions = {
        rate_limit = {
          count            = 1000
          period           = 60
          close_connection = true
        }
      }
    } // End policy 1
  ]   // End policies
}

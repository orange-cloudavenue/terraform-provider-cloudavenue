resource "cloudavenue_elb_policies_http_security" "example" {
  virtual_service_id = cloudavenue_elb_virtual_service.example.id
  policies = [
    # Policy 1
    {
      name = "example"

      # Define the criteria for the policy
      criteria = {
        // This example checks if the request is using HTTP on port 80.
        protocol = "HTTP"
        service_ports = {
          criteria = "IS_IN"
          ports    = ["80"]
        }
      }

      // The action rate_limit can be set to limit the number of requests.
      // The action redirect can be set to redirect the request on port 8443 when the rate limit is reached.
      // Only one action can be set at a time
      actions = {
        rate_limit = {
          count  = 1000
          period = 60
          redirect = {
            port     = 8443
            protocol = "HTTPS"
          }
        }
      }
    } // End policy 1
  ]   // End policies
}

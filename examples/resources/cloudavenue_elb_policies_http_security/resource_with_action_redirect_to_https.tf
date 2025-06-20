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

      // The action redirect_to_https must be a port number
      // This action redirects the request to HTTPS on the specified port.
      // Only one action can be set at a time
      actions = {
        redirect_to_https = 443
      }
    } // End policy 1
  ]   // End policies
}

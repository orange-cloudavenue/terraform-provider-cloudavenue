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

      // The action connection can be "ALLOW" or "CLOSE"
      // Only one ation can be set at a time
      actions = {
        connection = "ALLOW"
      }
    } // End policy 1
  ]   // End policies
}

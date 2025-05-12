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

      // The action send_response
      // Only one action can be set at a time
      // The send_response action can be used to send a custom response
      actions = {
        send_response = {
          status_code  = "403"
          content      = base64("Access Denied")
          content-type = "text/plain"
        }
      }
    } // End policy 1
  ]   // End policies
}

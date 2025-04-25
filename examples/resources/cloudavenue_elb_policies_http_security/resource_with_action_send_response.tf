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

      // The action send_response
      // The send_response action can be used to send a custom response
      // Only one action can be set at a time
      actions = {
        send_response = {
          status_code  = 403
          content      = base64("Access Denied")
          content_type = "text/plain"
        }
      }
    } // End policy 1
  ]   // End policies
}

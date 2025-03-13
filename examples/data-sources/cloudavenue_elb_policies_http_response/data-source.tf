data "cloudavenue_elb_policies_http_response" "example" {
  virtual_service_id = cloudavenue_elb_virtual_service.example.id
}

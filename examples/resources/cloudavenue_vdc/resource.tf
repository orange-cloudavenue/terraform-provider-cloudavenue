resource "cloudavenue_vdc" "example" {
  name                  = "MyVDC"
  description           = "Example VDC created by Terraform"
  cpu_allocated         = 22000
  memory_allocated      = 30
  cpu_speed_in_mhz      = 2200
  billing_model         = "PAYG"
  disponibility_class   = "ONE-ROOM"
  service_class         = "STD"
  storage_billing_model = "PAYG"

  storage_profiles = [
    {
      class   = "gold"
      default = true
      limit   = 500
    },
    {
      class   = "silver"
      default = false
      limit   = 500
    },
  ]
}

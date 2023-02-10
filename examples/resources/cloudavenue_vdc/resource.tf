resource "cloudavenue_vdc" "example" {
  name                  = "MyVDC"
  vdc_group             = "MyVDCGroup"
  description           = "Example VDC created by Terraform"
  cpu_allocated         = 6000
  memory_allocated      = 10
  cpu_speed_in_mhz      = 1200
  billing_model         = "PAYG"
  disponibility_class   = "ONE-ROOM"
  service_class         = "STD"
  storage_billing_model = "PAYG"

  storage_profile {
    class   = "gold"
    default = true
    limit   = 500
  }

  storage_profile {
    class   = "silver"
    default = false
    limit   = 500
  }

}

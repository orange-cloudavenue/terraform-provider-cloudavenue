resource "cloudavenue_vdc" "example" {
  name                      = "MyVDC"
  vdc_group                 = "MyVDCGroup"
  description               = "Example VDC created by Terraform"
  cpu_allocated             = 5
  memory_allocated          = 128
  vcpu_in_mhz2              = 2200
  vdc_billing_model         = "PAYG"
  vdc_disponibility_class   = "ONE-ROOM"
  vdc_service_class         = "STD"
  vdc_storage_billing_model = "PAYG"

  vdc_storage_profiles = [
    {
      class   = "gold"
      default = true
      limit   = 500
    }
  ]
}

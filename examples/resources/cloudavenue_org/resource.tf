import {
  to = cloudavenue_org.example
  # The id field is informational and is only mandatory when using the import function.
  # For clarity, it is recommended to specify the name of your organization as the id value. 
  # This helps to easily identify the imported resource in your Terraform state.
  # ex : cav01xx00000 
  id = "yourOrganizationName"
}

resource "cloudavenue_org" "example" {
  name                  = "Your Organization Name"
  description           = "This is an example organization"
  email                 = "example@mycompagny.com"
  internet_billing_mode = "PAYG"
}

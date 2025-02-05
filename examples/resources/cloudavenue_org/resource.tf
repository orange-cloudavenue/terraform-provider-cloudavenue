import {
  to = cloudavenue_org.example
  # ex : cav01xx00000   
  id = "yourOrganizationName"
}

resource "cloudavenue_org" "example" {
  name        = "Your Organization Name"
  description = "This is an example organization"
  email       = "example@mycompagny.com"
}

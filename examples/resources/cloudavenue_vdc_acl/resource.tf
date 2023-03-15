resource "cloudavenue_vdc_acl" "example" {
  vdc                   = "VDC_Test" # Optional
  everyone_access_level = "ReadOnly"
}

resource "cloudavenue_vdc_acl" "example" {
  vdc = "VDC_Test" # Optional
  shared_with = [
    {
      access_level = "ReadOnly"
      user_id      = "urn:vcloud:user:53665519-7036-43ea-ba97-63fc5a2aabe7"
    }
  ]
}
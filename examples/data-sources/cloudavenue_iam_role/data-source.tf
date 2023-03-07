resource "cloudavenue_iam_role" "example" {
  name        = "roletest"
  description = "A test role"
  rights = [
    "Catalog: Add vApp from My Cloud",
    "Catalog: Edit Properties",
    "Catalog: View Private and Shared Catalogs",
    "Organization vDC Compute Policy: View",
    "vApp Template / Media: Edit",
    "vApp Template / Media: View",
  ]
}

data "cloudavenue_iam_role" "example" {
  name = cloudavenue_iam_role.example.name
}

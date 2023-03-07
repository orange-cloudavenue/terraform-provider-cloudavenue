resource "cloudavenue_iam_group" "example" {
  name        = "OrgTest"
  role        = "Organization Administrator"
  description = "org test from go test"
}

data "cloudavenue_iam_group" "example" {
  name = "your_value"
}
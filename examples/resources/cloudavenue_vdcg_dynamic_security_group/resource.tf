resource "cloudavenue_vdcg_dynamic_security_group" "example" {
  name         = "example"
  description  = "example of dynamic security group with criteria"
  vdc_group_id = cloudavenue_vdcg.example.id
  criteria = [
    { # OR
      rules = [
        { # AND
          type     = "VM_NAME"
          value    = "test"
          operator = "STARTS_WITH"
        },
        { # AND
          type     = "VM_NAME"
          value    = "front"
          operator = "CONTAINS"
        },
        { # AND
          type     = "VM_TAG"
          value    = "prod"
          operator = "ENDS_WITH"
        },
      ]
    },
    { # OR
      rules = [
        { # AND
          type     = "VM_TAG"
          value    = "test"
          operator = "STARTS_WITH"
        },
        { # AND
          type     = "VM_TAG"
          value    = "web-front"
          operator = "CONTAINS"
        }
      ]
    },
    { # OR
      rules = [
        { # AND
          type     = "VM_TAG"
          value    = "prod"
          operator = "STARTS_WITH"
        },
        { # AND
          type     = "VM_TAG"
          value    = "test-xx"
          operator = "EQUALS"
        }
      ]
    }
  ]
}

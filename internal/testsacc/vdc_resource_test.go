package testsacc

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const testAccVDCResourceConfig = `
resource "cloudavenue_vdc" "example" {
	name                  = "MyVDC1"
	vdc_group             = "MyGroup"
	description           = "Example vDC created by Terraform"
	cpu_allocated         = 22000
	memory_allocated      = 30
	cpu_speed_in_mhz      = 2200
	billing_model         = "PAYG"
	disponibility_class   = "ONE-ROOM"
	service_class         = "STD"
	storage_billing_model = "PAYG"
  
	storage_profiles = [{
	  class   = "gold"
	  default = true
	  limit   = 500
	}]
  
  }
`

// Temp test during deprecation notice of vdc_group attribute
// This test will be removed when vdc_group attribute will be removed
// Used in vdc_group_resource_test.go.
const TestAccVDCResourceConfigWithoutVDCGroup = `
resource "cloudavenue_vdc" "example" {
	name                  = "MyVDCExample"
	description           = "Example vDC created by Terraform"
	cpu_allocated         = 22000
	memory_allocated      = 30
	cpu_speed_in_mhz      = 2200
	billing_model         = "PAYG"
	disponibility_class   = "ONE-ROOM"
	service_class         = "STD"
	storage_billing_model = "PAYG"
  
	storage_profiles = [{
	  class   = "gold"
	  default = true
	  limit   = 500
	}]
  
}

resource "cloudavenue_vdc" "example2" {
	name                  = "MyVDCExample2"
	description           = "Example vDC created by Terraform"
	cpu_allocated         = 22000
	memory_allocated      = 30
	cpu_speed_in_mhz      = 2200
	billing_model         = "PAYG"
	disponibility_class   = "ONE-ROOM"
	service_class         = "STD"
	storage_billing_model = "PAYG"
  
	storage_profiles = [{
	  class   = "gold"
	  default = true
	  limit   = 500
	}]
  
}
`

func TestAccVDCResource(t *testing.T) {
	const resourceName = "cloudavenue_vdc.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Apply test
				Config: testAccVDCResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.VDC.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "name", "MyVDC1"),
					resource.TestCheckResourceAttr(resourceName, "vdc_group", "MyGroup"),
					resource.TestCheckResourceAttr(resourceName, "description", "Example vDC created by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
					resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
					resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
					resource.TestCheckResourceAttr(resourceName, "billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "disponibility_class", "ONE-ROOM"),
					resource.TestCheckResourceAttr(resourceName, "service_class", "STD"),
					resource.TestCheckResourceAttr(resourceName, "storage_billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.class", "gold"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.default", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.limit", "500"),
				),
			},
			// Uncomment if you want to test update or delete this block
			{
				// Update test
				Config: strings.Replace(testAccVDCResourceConfig, "30", "40", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(uuid.VDC.String()+`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)),
					resource.TestCheckResourceAttr(resourceName, "name", "MyVDC1"),
					resource.TestCheckResourceAttr(resourceName, "vdc_group", "MyGroup"),
					resource.TestCheckResourceAttr(resourceName, "description", "Example vDC created by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
					resource.TestCheckResourceAttr(resourceName, "memory_allocated", "40"),
					resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
					resource.TestCheckResourceAttr(resourceName, "billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "disponibility_class", "ONE-ROOM"),
					resource.TestCheckResourceAttr(resourceName, "service_class", "STD"),
					resource.TestCheckResourceAttr(resourceName, "storage_billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.class", "gold"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.default", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.limit", "500"),
				),
			},
			// ImportruetState testing
			{
				// Import test with vdc
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           "MyVDC1",
				ImportStateVerifyIgnore: []string{"vdc_group"},
			},
		},
	})
}

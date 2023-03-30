package vdc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tests "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/tests/common"
)

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -test
const testAccVDCGroupDataSourceConfig = `
data "cloudavenue_vdc_group" "example" {
	name = "MyVDCGroup"
}
`

func TestAccVDCGroupDataSource(t *testing.T) {
	const dataSourceName = "data.cloudavenue_vdc_group.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{ // Read testing
			{
				Config: testAccVDCGroupDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", "MyVDCGroup"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "dfw_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "local_egress"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_pool_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_provider_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "universal_networking_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdcs.0.fault_domain_tag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdcs.0.network_provider_scope"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdcs.0.is_remote_org"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdcs.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdcs.0.site_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdcs.0.site_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdcs.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vdcs.0.id"),
				),
			},
		},
	})
}

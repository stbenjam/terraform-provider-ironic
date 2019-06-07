// +build acceptance

package ironic

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceAPI(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ironic_api" "wait" {
					timeout = 30
				}`,
				Check: resource.TestCheckResourceAttr("data.ironic_api.wait", "timeout", "30"),
			},
		},
	})
}

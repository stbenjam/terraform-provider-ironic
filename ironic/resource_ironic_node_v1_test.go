package ironic

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// Creates a node in Ironic that uses the fake-hardware driver
func TestAccIronicNode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		//CheckDestroy: testAccNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeResource,
				Check: resource.ComposeTestCheckFunc(
		//			testAccCheckNodeExists(),
				),
			},
		},
	})
}

const (
	testAccNodeResource = `
		resource "ironic_node_v1" "node-0" {
			name = "node-0"
			driver = "fake-hardware"
			boot_interface = "pxe"
			driver_info {
				ipmi_port      = "6230",
				ipmi_username  = "admin",
				deploy_kernel  = "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
				ipmi_address   = "192.168.122.1",
				deploy_ramdisk = "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
				ipmi_password  = "admin",
			}
		}`
)

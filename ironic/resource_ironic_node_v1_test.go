package ironic

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIronicNode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNodeDestroy,
		Steps: []resource.TestStep{
			// Create a node and check that it exists
			{
				Config: testAccNodeResource(""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists("ironic_node_v1.node-0"),
					resource.TestCheckResourceAttr("ironic_node_v1.node-0",
						"provision_state", "enroll",
					),
				),
			},

			// Change the node's provision state to 'available', and inspect the node
			// in the process.
			{
				Config: testAccNodeResource(`
					target_provision_state = "provide"
					inspect = true
				`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists("ironic_node_v1.node-0"),
					resource.TestCheckResourceAttr("ironic_node_v1.node-0",
						"provision_state", "available"),
					resource.TestCheckResourceAttrSet("ironic_node_v1.node-0",
						"inspection_finished_at"),
				),
			},

			// Change the node's power state to 'power on', with a timeout
			{
				Config: testAccNodeResource(`
					target_power_state = "power on"
					power_state_timeout = 10
				`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists("ironic_node_v1.node-0"),
					resource.TestCheckResourceAttr("ironic_node_v1.node-0",
						"power_state", "power on"),
				),
			},

			// Change the node's power state to 'power off'.
			{
				Config: testAccNodeResource("target_power_state = \"power off\""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists("ironic_node_v1.node-0"),
					resource.TestCheckResourceAttr("ironic_node_v1.node-0",
						"power_state", "power off"),
				),
			},

			// Change the node's power state to 'rebooting', it probably
			// doesn't make a whole lot of sense for a terraform user to
			// declare a node's state as forever rebooting, as it'd reboot
			// every time, but we should check anyway that if they do say
			// rebooting, power_state goes to power_on and terraform exits
			// successfully.
			{
				Config: testAccNodeResource("target_power_state = \"rebooting\""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists("ironic_node_v1.node-0"),
					resource.TestCheckResourceAttr("ironic_node_v1.node-0",
						"power_state", "power on"),
				),
			},

			// Change the node's provision state to 'active'
			{
				Config: testAccNodeResource(`
					target_provision_state = "active"
				`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeExists("ironic_node_v1.node-0"),
					resource.TestCheckResourceAttr("ironic_node_v1.node-0",
						"provision_state", "active"),
				),
			},
		},
	})
}

func testAccCheckNodeExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*gophercloud.ServiceClient)

		rs, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no node ID is set")
		}

		_, err := nodes.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("node (%s) not found: %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccNodeDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*gophercloud.ServiceClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "ironic_node_v1" {
			continue
		}

		_, err := nodes.Get(client, rs.Primary.ID).Extract()
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return fmt.Errorf("unexpected error: %s, expected 404", err)
		}
	}

	return nil
}

func testAccNodeResource(extraValue string) string {
	return fmt.Sprintf(`resource "ironic_node_v1" "node-0" {
			name = "node-0"
			driver = "fake-hardware"

			boot_interface = "pxe"
			deploy_interface = "fake"
			management_interface = "fake"
			power_interface = "fake"
            inspect_interface = "fake"
			resource_class = "baremetal"
			vendor_interface = "no-vendor"

			driver_info = {
				ipmi_port      = "6230"
				ipmi_username  = "admin"
				ipmi_password  = "admin"
			}

			%s
		}`, extraValue)
}

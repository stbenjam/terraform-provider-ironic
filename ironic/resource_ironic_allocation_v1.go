package ironic

import (
	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform/helper/schema"
)

// Schema resource definition for an Ironic allocation.
func resourceAllocationV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceAllocationV1Create,
		Read:   resourceAllocationV1Read,
		Delete: resourceAllocationV1Delete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"candidate_nodes": {
				Type:     schema.TypeList,
				Optional: true,
			},
			"traits": {
				Type:     schema.TypeList,
				Optional: true,
			},
			"extra": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"node_uuid": {
				Type: schema.TypeString,
				Computed: true,
			},
			"state": {
				Type: schema.TypeString,
				Computed: true,
			},
			"last_error": {
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Create a allocation, including driving Ironic's state machine
func resourceAllocationV1Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)
}

// Read the allocation's data from Ironic
func resourceAllocationV1Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)
}

// Delete a allocation from Ironic
func resourceAllocationV1Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)
}

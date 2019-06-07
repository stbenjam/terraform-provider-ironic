package ironic

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/http"
	"time"
)

func dataSourceIronicAPI() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIronicAPIRead,

		Schema: map[string]*schema.Schema{
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
		},
	}
}

func dataSourceIronicAPIRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)

	wait := time.Duration(d.Get("timeout").(int))

	// Repeatedly try to connect to the API, up to the defined number of seconds.  If we cannot connect, then we
	// return an error.
	apiIsUp := make(chan bool, 1)
	timeout := make(chan bool, 1)
	defer close(apiIsUp)
	defer close(timeout)

	go func() {
		for tries := 0; ; tries++ {
			// Return if we've hit the timeout
			select {
			case <-timeout:
				return
			default:
			}

			// Otherwise keep trying the IronicAPI
			resp, err := client.HTTPClient.Get(client.Endpoint)
			if err == nil && resp.StatusCode == http.StatusOK {
				apiIsUp <- true
				return
			} else if err != nil {
				log.Printf("[DEBUG] Ironic API get Error: %s", err.Error())
			} else {
				log.Printf("[DEBUG] Ironic HTTP status code: %d", resp.StatusCode)
			}
			log.Printf("[INFO] Ironic API is not yet up, attempt %d.", tries)
			time.Sleep(5 * time.Second)
		}

	}()

	select {
	case _, ok := <-apiIsUp:
		if !ok {
			return fmt.Errorf("api channel unexpectedly closed")
		}
	case <-time.After(time.Duration(wait) * time.Second):
		timeout <- true
		errMsg := fmt.Sprintf("Ironic api failed to come up after %d seconds", wait)
		log.Printf("[ERROR] %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	return nil
}

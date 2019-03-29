package ironic

import (
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	th "github.com/openshift-metalkube/terraform-provider-ironic/testhelper"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"ironic": testAccProvider,
	}
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	v := os.Getenv("IRONIC_ENDPOINT")
	if v == "" {
		t.Fatal("IRONIC_ENDPOINT must be set for acceptance tests")
	}
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
}

func TestAccProvider(t *testing.T) {
	testAccPreCheck(t)

	p := Provider()
	raw, err := config.NewRawConfig(map[string]interface{}{
		"url":          "http://localhost:6385/v1",
		"microversion": "1.52",
	})
	th.AssertNoError(t, err)

	err = p.Configure(terraform.NewResourceConfig(raw))
	th.AssertNoError(t, err)
}

func TestAccProvider_urlRequired(t *testing.T) {
	testAccPreCheck(t)

	p := Provider()
	raw, err := config.NewRawConfig(map[string]interface{}{})
	th.AssertNoError(t, err)

	err = p.Configure(terraform.NewResourceConfig(raw))
	th.AssertError(t, err, "url is required")
}

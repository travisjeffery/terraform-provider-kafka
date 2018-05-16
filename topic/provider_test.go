package topic

import (
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"kafka": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatal(err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	c, err := config.NewRawConfig(map[string]interface{}{
		"hosts": []string{"localhost:9092"},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = testAccProvider.Configure(terraform.NewResourceConfig(c))
	if err != nil {
		t.Fatal(err)
	}
}

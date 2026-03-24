package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTokenResource_basic(t *testing.T) {
	resourceName := "datafy_token.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTokenResourceConfig("regression-test-token", "60m"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "token_id"),
					resource.TestCheckResourceAttrSet(resourceName, "secret"),
					resource.TestCheckResourceAttrSet(resourceName, "expires"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "description", "regression-test-token"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "60m"),
				),
			},
		},
	})
}

func TestAccTokenResource_noTTL(t *testing.T) {
	resourceName := "datafy_token.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTokenResourceConfigNoTTL("regression-test-token-no-ttl"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token_id"),
					resource.TestCheckResourceAttrSet(resourceName, "secret"),
					resource.TestCheckResourceAttr(resourceName, "description", "regression-test-token-no-ttl"),
				),
			},
		},
	})
}

func testAccTokenResourceConfig(description, ttl string) string {
	return fmt.Sprintf(`
resource "datafy_account" "test" {
  name = "regression-test-token"
}

resource "datafy_token" "test" {
  account_id  = datafy_account.test.id
  description = %q
  ttl         = %q
  role_ids    = []
}
`, description, ttl)
}

func testAccTokenResourceConfigNoTTL(description string) string {
	return fmt.Sprintf(`
resource "datafy_account" "test" {
  name = "regression-test-token-no-ttl"
}

resource "datafy_token" "test" {
  account_id  = datafy_account.test.id
  description = %q
  role_ids    = []
}
`, description)
}

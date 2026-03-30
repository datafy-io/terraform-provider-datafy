package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTokenDataSource_basic(t *testing.T) {
	resourceName := "data.datafy_token.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTokenDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "token_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "expires"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
		},
	})
}

func testAccTokenDataSourceConfig() string {
	return `
resource "datafy_account" "test" {
  name = "regression-test-token-ds"
}

resource "datafy_token" "test" {
  account_id  = datafy_account.test.id
  description = "regression-test-token-ds"
  ttl         = "60m"
  role_ids    = []
}

data "datafy_token" "test" {
  account_id = datafy_account.test.id
  token_id   = datafy_token.test.token_id
}
`
}

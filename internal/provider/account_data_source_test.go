package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountDataSource_basic(t *testing.T) {
	resourceName := "data.datafy_account.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "parent_account_id"),
				),
			},
		},
	})
}

func testAccAccountDataSourceConfig() string {
	return `
resource "datafy_account" "test" {
  name = "regression-test-account-ds"
}

data "datafy_account" "test" {
  id = datafy_account.test.id
}
`
}

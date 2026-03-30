package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleArnDataSource_basic(t *testing.T) {
	resourceName := "data.datafy_role_arn.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleArnDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
		},
	})
}

func testAccRoleArnDataSourceConfig() string {
	return `
resource "datafy_account" "test" {
  name = "regression-test-rolearn-ds"
}

resource "datafy_role_arn" "test" {
  account_id      = datafy_account.test.id
  arn             = "arn:aws:iam::123456789012:role/regression-test-role-ds"
  skip_validation = true
}

data "datafy_role_arn" "test" {
  account_id = datafy_account.test.id

  depends_on = [datafy_role_arn.test]
}
`
}

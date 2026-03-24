package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleArnResource_basic(t *testing.T) {
	resourceName := "datafy_role_arn.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleArnResourceConfig(
					"arn:aws:iam::123456789012:role/regression-test-role",
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttr(resourceName, "arn", "arn:aws:iam::123456789012:role/regression-test-role"),
					resource.TestCheckResourceAttr(resourceName, "skip_validation", "true"),
				),
			},
		},
	})
}

func TestAccRoleArnResource_update(t *testing.T) {
	resourceName := "datafy_role_arn.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleArnResourceConfig(
					"arn:aws:iam::123456789012:role/regression-test-role",
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "arn", "arn:aws:iam::123456789012:role/regression-test-role"),
				),
			},
			{
				Config: testAccRoleArnResourceConfig(
					"arn:aws:iam::123456789012:role/regression-test-role-updated",
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "arn", "arn:aws:iam::123456789012:role/regression-test-role-updated"),
				),
			},
		},
	})
}

func testAccRoleArnResourceConfig(arn string, skipValidation bool) string {
	return fmt.Sprintf(`
resource "datafy_account" "test" {
  name = "regression-test-rolearn"
}

resource "datafy_role_arn" "test" {
  account_id      = datafy_account.test.id
  arn             = %q
  skip_validation = %t
}
`, arn, skipValidation)
}

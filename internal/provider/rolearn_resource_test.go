package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRoleArnResource_basic(t *testing.T) {
	resourceName := "datafy_role_arn.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleArnDestroy,
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
		CheckDestroy:             testAccCheckRoleArnDestroy,
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

func testAccCheckRoleArnDestroy(s *terraform.State) error {
	client := newTestClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "datafy_role_arn" {
			continue
		}

		_, err := client.GetAccountRoleArn(context.Background(), &datafy.GetAccountRoleArnRequest{
			AccountId: rs.Primary.Attributes["account_id"],
		})
		if err == nil {
			return fmt.Errorf("role_arn for account %s still exists after destroy", rs.Primary.Attributes["account_id"])
		}
	}

	return nil
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

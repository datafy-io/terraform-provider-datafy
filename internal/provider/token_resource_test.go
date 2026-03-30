package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTokenResource_basic(t *testing.T) {
	resourceName := "datafy_token.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTokenDestroy,
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
		CheckDestroy:             testAccCheckTokenDestroy,
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

func testAccCheckTokenDestroy(s *terraform.State) error {
	client := newTestClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "datafy_token" {
			continue
		}

		_, err := client.GetAccountToken(context.Background(), &datafy.GetAccountTokenRequest{
			AccountId: rs.Primary.Attributes["account_id"],
			TokenId:   rs.Primary.Attributes["token_id"],
		})
		if err == nil {
			return fmt.Errorf("token %s still exists after destroy", rs.Primary.Attributes["token_id"])
		}
	}

	return nil
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

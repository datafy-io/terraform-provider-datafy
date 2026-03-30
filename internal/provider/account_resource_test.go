package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAccountResource_basic(t *testing.T) {
	resourceName := "datafy_account.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountResourceConfig("regression-test-account"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "regression-test-account"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "parent_account_id"),
				),
			},
		},
	})
}

func TestAccAccountResource_update(t *testing.T) {
	resourceName := "datafy_account.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountResourceConfig("regression-test-account"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "regression-test-account"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testAccAccountResourceConfig("regression-test-account-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "regression-test-account-updated"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCheckAccountDestroy(s *terraform.State) error {
	client := newTestClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "datafy_account" {
			continue
		}

		_, err := client.GetAccount(context.Background(), &datafy.GetAccountRequest{
			AccountId: rs.Primary.Attributes["id"],
		})
		if err == nil {
			return fmt.Errorf("account %s still exists after destroy", rs.Primary.Attributes["id"])
		}
	}

	return nil
}

func testAccAccountResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "datafy_account" "test" {
  name = %q
}
`, name)
}

func newTestClient() *datafy.Client {
	token := os.Getenv("DATAFY_TOKEN")
	endpoint := os.Getenv("DATAFY_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.datafy.io"
	}
	return datafy.NewClient(token, endpoint)
}

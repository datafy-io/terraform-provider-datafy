package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAutoscalingRuleResource_basic(t *testing.T) {
	resourceName := "datafy_autoscaling_rule.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAutoscalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingRuleResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_id"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
				),
			},
		},
	})
}

func TestAccAutoscalingRuleResource_update(t *testing.T) {
	resourceName := "datafy_autoscaling_rule.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAutoscalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingRuleResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_id"),
				),
			},
			{
				Config: testAccAutoscalingRuleResourceConfigUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
				),
			},
		},
	})
}

func testAccCheckAutoscalingRuleDestroy(s *terraform.State) error {
	client := newTestClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "datafy_autoscaling_rule" {
			continue
		}

		_, err := client.GetAccountAutoscalingRule(context.Background(), &datafy.GetAccountAutoscalingRuleRequest{
			AccountId: rs.Primary.Attributes["account_id"],
			RuleId:    rs.Primary.Attributes["rule_id"],
		})
		if err == nil {
			return fmt.Errorf("autoscaling_rule %s still exists after destroy", rs.Primary.Attributes["rule_id"])
		}
	}

	return nil
}

func testAccAutoscalingRuleResourceConfig(active bool) string {
	activeStr := "true"
	if !active {
		activeStr = "false"
	}
	return `
resource "datafy_account" "test" {
  name = "regression-test-rule"
}

resource "datafy_autoscaling_rule" "test" {
  account_id = datafy_account.test.id
  active     = ` + activeStr + `
  rule       = jsonencode({
    "in" = [
      { "var" = "cluster_name" },
      ["regression-test-cluster"]
    ]
  })
}
`
}

func testAccAutoscalingRuleResourceConfigUpdated() string {
	return `
resource "datafy_account" "test" {
  name = "regression-test-rule"
}

resource "datafy_autoscaling_rule" "test" {
  account_id = datafy_account.test.id
  active     = false
  rule       = jsonencode({
    "and" = [
      {
        "in" = [
          { "var" = "cluster_name" },
          ["regression-test-cluster"]
        ]
      },
      {
        "in" = [
          { "var" = "node_group_name" },
          ["regression-test-nodegroup"]
        ]
      }
    ]
  })
}
`
}

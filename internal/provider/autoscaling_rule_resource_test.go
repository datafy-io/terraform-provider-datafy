package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoscalingRuleResource_basic(t *testing.T) {
	resourceName := "datafy_autoscaling_rule.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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

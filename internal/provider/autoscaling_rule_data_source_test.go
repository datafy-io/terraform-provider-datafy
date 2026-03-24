package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoscalingRuleDataSource_basic(t *testing.T) {
	resourceName := "data.datafy_autoscaling_rule.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingRuleDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_id"),
				),
			},
		},
	})
}

func testAccAutoscalingRuleDataSourceConfig() string {
	return `
resource "datafy_account" "test" {
  name = "regression-test-rule-ds"
}

resource "datafy_autoscaling_rule" "test" {
  account_id = datafy_account.test.id
  active     = true
  rule       = jsonencode({
    "in" = [
      { "var" = "cluster_name" },
      ["regression-test-cluster"]
    ]
  })
}

data "datafy_autoscaling_rule" "test" {
  account_id = datafy_account.test.id
  rule_id    = datafy_autoscaling_rule.test.rule_id
}
`
}

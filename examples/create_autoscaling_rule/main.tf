terraform {
  required_providers {
    datafy = {
      source = "datafy-io/datafy"
    }
  }
}

provider "datafy" {
  token = "eZa0qICnUV-COvO46NfDysUDN4bFKMeWssXVCIsIIn0.eyJzdW"
}

resource "datafy_account" "example" {
  name = "example-account"
}

resource "datafy_autoscaling_rule" "example1" {
  account_id = datafy_account.example.id
  active     = false
  rule = jsonencode({
    "in" : [
      { "var" : "instance_id" },
      [
        "i-1234567890",
        "i-1234567891"
      ]
    ]
  })
}

resource "datafy_autoscaling_rule" "example2" {
  account_id = datafy_account.example.id
  active     = true
  rule = jsonencode({
    "some" : [
      { "var" : "tags" },
      { "in" : [
        { "var" : "" },
        [
          "env:stg",
          "env:staging"
        ]
      ] }
    ]
  })
}

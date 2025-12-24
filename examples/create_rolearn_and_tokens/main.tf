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

resource "datafy_role_arn" "example" {
  account_id = datafy_account.example.id
  arn        = "arn:aws:iam::123456789012:role/example-role"
}

resource "datafy_token" "with_expiration" {
  account_id  = datafy_account.example.id
  description = "Token with expiration"
  ttl         = "60m"
  role_ids    = ["role1", "role2"]
}

resource "datafy_token" "no_expiration" {
  account_id  = datafy_account.example.id
  description = "Token with no expiration"
  role_ids    = ["role1", "role2"]
}

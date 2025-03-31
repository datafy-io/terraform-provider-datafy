terraform {
  required_providers {
    hashicups = {
      source = "datafy-io/datafy"
    }
  }
}

resource "datafy_account" "example" {
  name = "example-account"
}

resource "datafy_role_arn" "example" {
  account_id = datafy_account.example.id
  arn        = "arn:aws:iam::123456789012:role/example-role"
}

resource "datafy_token" "with_expiration" {
  account_id        = datafy_account.example.id
  description       = "Token with expiration"
  expire_in_minutes = 60
  role_ids          = ["role1", "role2"]
}

resource "datafy_token" "no_expiration" {
  account_id  = datafy_account.example.id
  description = "Token with no expiration"
  role_ids    = ["role1", "role2"]
}

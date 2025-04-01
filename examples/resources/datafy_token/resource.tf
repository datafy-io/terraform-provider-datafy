resource "datafy_token" "example" {
  account_id  = "79c406c5-7b64-43f2-ba76-9b01e74e3d90"
  description = "Example token"
  ttl         = "60m"
  role_ids    = ["role1", "role2"]
}

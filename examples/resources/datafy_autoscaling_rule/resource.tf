resource "datafy_autoscaling_rule" "example" {
  account_id = "79c406c5-7b64-43f2-ba76-9b01e74e3d90"
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

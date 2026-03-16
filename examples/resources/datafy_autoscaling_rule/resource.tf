resource "datafy_autoscaling_rule" "example" {
  account_id = datafy_account.example.id
  active     = true
  rule = jsonencode({
    "and" : [
      {
        "in" : [
          { "var" : "instance_id" },
          [
            "i-1234567890",
            "i-1234567891"
          ]
        ]
      }
    ]
  })
}

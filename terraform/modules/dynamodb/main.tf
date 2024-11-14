resource "aws_dynamodb_table" "item_table" {
  name         = var.table_name
  hash_key     = var.hash_key
  billing_mode = var.billing_mode

  attribute {
    name = var.hash_key
    type = "S"
  }

  tags = var.tags
}
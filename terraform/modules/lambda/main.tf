resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda_exec_role_${var.suffix}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb_access" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}

resource "aws_lambda_function" "function" {
  filename         = var.filename
  function_name    = var.function_name
  role             = aws_iam_role.lambda_exec_role.arn
  handler          = var.handler
  source_code_hash = filebase64sha256(var.filename)
  runtime          = var.runtime
  timeout          = var.timeout
}

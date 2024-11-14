
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

module "lambda" {
  source         = "./modules/lambda"
  filename       = "${path.module}/../function.zip"
  function_name  = "go_lambda_function"
  handler        = "bootstrap"
  runtime        = "provided.al2"  # Ensure this is a supported runtime
  timeout        = 30
  suffix         = random_string.suffix.result
}

resource "aws_api_gateway_rest_api" "api" {
  name = "enterprise-manager-api"
  body = templatefile("${path.module}/enterprise-manager-api.tpl", {
    lambda_arn = module.lambda.lambda_function_arn
  })
}

resource "aws_lambda_permission" "apigw_invoke" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api.execution_arn}/*/*"
}

resource "aws_api_gateway_deployment" "deployment" {
  depends_on  = [aws_api_gateway_rest_api.api]
  rest_api_id = aws_api_gateway_rest_api.api.id
}

resource "aws_api_gateway_stage" "stage" {
  deployment_id = aws_api_gateway_deployment.deployment.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
  stage_name    = "dev"
}


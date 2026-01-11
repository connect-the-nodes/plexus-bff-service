# 1. Lambda Log Group
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/zynchub-${var.environment}"
  retention_in_days = 14
}

# 2. API Gateway Access Log Group
resource "aws_cloudwatch_log_group" "api_gw" {
  name              = "/aws/api-gw/zynchub-${var.environment}"
  retention_in_days = 7
}

# 3. Correct v1 Stage with Logging
resource "aws_api_gateway_stage" "hub_stage" {
  deployment_id = aws_api_gateway_deployment.hub_deploy.id
  rest_api_id   = aws_api_gateway_rest_api.hub_api.id # Now using .id correctly
  stage_name    = var.environment

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw.arn
    format          = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      requestTime    = "$context.requestTime"
      httpMethod     = "$context.httpMethod"
      resourcePath   = "$context.resourcePath"
      status         = "$context.status"
      protocol       = "$context.protocol"
      responseLength = "$context.responseLength"
    })
  }
}
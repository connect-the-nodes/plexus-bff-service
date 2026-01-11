resource "aws_api_gateway_rest_api" "hub_api" {
  name        = "zynchub-api-${var.environment}"
  description = "API Gateway for Zynchub Digital Hub Service"

  body = var.openapi_spec

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "hub_deploy" {
  rest_api_id = aws_api_gateway_rest_api.hub_api.id

  triggers = {
    redeployment = sha1(aws_api_gateway_rest_api.hub_api.body)
  }

  lifecycle {
    create_before_destroy = true
  }
}

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
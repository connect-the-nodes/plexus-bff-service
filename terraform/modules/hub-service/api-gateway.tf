resource "aws_api_gateway_rest_api" "hub_api" {
  name        = "zynchub-api-${var.environment}"
  description = "API Gateway for Zynchub Digital Hub Service"

  body = templatefile("${path.module}/../../src/main/resources/static/openapi.yaml", {
    lambda_arn = aws_lambda_function.hub_service.arn
    region     = var.aws_region
  })

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
  rest_api_id   = aws_api_gateway_rest_api.hub_api.id
  stage_name    = var.environment
}
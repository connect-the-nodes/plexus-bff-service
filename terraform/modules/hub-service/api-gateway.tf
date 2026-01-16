# Render the OpenAPI spec dynamically
locals {
  openapi_spec = templatefile(var.openapi_spec_path, {
    nlb_dns     = aws_lb.nlb.dns_name         # internal NLB
    vpc_link_id = aws_api_gateway_vpc_link.main.id  # internal VPC link
  })
}

############################################
# API Gateway REST API (OpenAPI-driven)
############################################
resource "aws_api_gateway_rest_api" "hub_api" {
  name = "zynchub-api-${var.environment}"
  body = local.openapi_spec

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

############################################
# VPC Link (API Gateway -> NLB)
############################################
resource "aws_api_gateway_vpc_link" "main" {
  name        = "zynchub-vpc-link-${var.environment}"
  target_arns = [aws_lb.nlb.arn]
}

############################################
# API Deployment (Redeploys on OpenAPI change)
############################################
resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.hub_api.id

  triggers = {
    redeploy = sha1(local.openapi_spec)
  }

  lifecycle {
    create_before_destroy = true
  }
}

############################################
# API Stage
############################################
resource "aws_api_gateway_stage" "this" {
  stage_name    = var.environment
  rest_api_id   = aws_api_gateway_rest_api.hub_api.id
  deployment_id = aws_api_gateway_deployment.this.id

  access_log_settings {
    destination_arn = local.apigw_log_group_arn
    format = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      requestTime    = "$context.requestTime"
      httpMethod     = "$context.httpMethod"
      resourcePath   = "$context.resourcePath"
      status         = "$context.status"
      responseLength = "$context.responseLength"
    })
  }

  xray_tracing_enabled = false

  depends_on = [
    aws_api_gateway_account.this
  ]
}

############################################
# Enable Logs & Metrics for ALL Methods
############################################
resource "aws_api_gateway_method_settings" "all" {
  rest_api_id = aws_api_gateway_rest_api.hub_api.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  method_path = "*/*"

  settings {
    logging_level      = "INFO"
    data_trace_enabled = false
    metrics_enabled    = true
  }
}

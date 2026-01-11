resource "aws_lambda_function" "hub_service" {
  function_name = var.service_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = "com.zynchub.digital.hubservice.ZynchubApplication::handleRequest"
  runtime       = "java21"

  # Use S3 configuration
  s3_bucket = var.artifact_bucket
  s3_key    = var.artifact_key


  layers = []

  publish   = true
  snap_start {
    apply_on = "PublishedVersions"
  }

  environment {
    variables = {
      APP_VERSION = var.app_version
      SPRING_PROFILES_ACTIVE = var.environment
    }
  }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.hub_service.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.hub_api.execution_arn}/*/*"
}
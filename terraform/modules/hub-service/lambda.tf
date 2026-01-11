resource "aws_lambda_function" "hub_service" {
  function_name = var.service_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = "org.springframework.cloud.function.adapter.aws.FunctionInvoker::handleRequest"
  runtime       = "java21"

  s3_bucket = var.artifact_bucket
  s3_key    = var.artifact_key

  # This makes the Lambda "reactive" to code changes in S3
  environment {
    variables = {
      MAIN_CLASS             = "com.zynchub.digital.hubservice.ZynchubApplication"
      APP_VERSION = var.app_version
      DEPLOY_VERSION = var.app_version
      SPRING_PROFILES_ACTIVE = var.environment
    }
  }

  publish   = true
  snap_start {
    apply_on = "PublishedVersions"
  }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.hub_service.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.hub_api.execution_arn}/*/*"
}
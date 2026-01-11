resource "aws_lambda_function" "hub_service" {
  function_name = var.service_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = "com.zynchub.digital.hubservice.ZynchubApplication::handleRequest"
  runtime       = "java21"

  # Use the variables here
  s3_bucket = var.artifact_bucket
  s3_key    = var.artifact_key
  publish   = true
  tags = {
    Version     = var.app_version
    Environment = var.environment
  }

  source_code_hash = filebase64sha256(var.jar_path)
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
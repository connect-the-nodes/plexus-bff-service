output "api_url" {
  description = "The invocation URL for the API Gateway stage"
  value       = "${aws_api_gateway_stage.hub_stage.invoke_url}"
}

output "lambda_function_name" {
  description = "Name of the deployed Lambda function"
  value       = aws_lambda_function.hub_service.function_name
}
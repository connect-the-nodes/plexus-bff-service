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

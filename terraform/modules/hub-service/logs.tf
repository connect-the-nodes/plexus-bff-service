# 1. Lambda Log Group
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/zynchub-${var.environment}"
  retention_in_days = 14
}

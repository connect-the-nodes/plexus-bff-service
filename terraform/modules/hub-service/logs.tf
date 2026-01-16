######################################
# LOG GROUPS
######################################

locals {
  apigw_log_group_name = "/aws/apigateway/zynchub-${var.environment}"
  apigw_role_name      = "zynchub-apigw-logs-${var.environment}"
}

# ECS Log Group
resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/zynchub-${var.environment}"
  retention_in_days = 14
  tags              = var.tags
}

# API Gateway Log Group
resource "aws_cloudwatch_log_group" "apigw" {
  count             = var.create_apigw_log_group ? 1 : 0
  name              = local.apigw_log_group_name
  retention_in_days = 14
  tags              = var.tags
}

data "aws_cloudwatch_log_group" "apigw" {
  count = var.create_apigw_log_group ? 0 : 1
  name  = local.apigw_log_group_name
}

locals {
  apigw_log_group_arn = var.create_apigw_log_group ? aws_cloudwatch_log_group.apigw[0].arn : data.aws_cloudwatch_log_group.apigw[0].arn
}

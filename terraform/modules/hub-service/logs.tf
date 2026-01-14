######################################
# LOG GROUPS
######################################

# 1️⃣ ECS Log Group
resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/zynchub-${var.environment}"
  retention_in_days = 14
  tags              = var.tags
}

# 2️⃣ API Gateway Log Group
resource "aws_cloudwatch_log_group" "apigw" {
  name              = "/aws/apigateway/zynchub-${var.environment}"
  retention_in_days = 14
  tags              = var.tags
}

######################################
# API Gateway → CloudWatch Logs Role
######################################
resource "aws_iam_role" "apigw_cloudwatch_role" {
  name = "zynchub-apigw-logs-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "apigateway.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "apigw_cloudwatch_attach" {
  role       = aws_iam_role.apigw_cloudwatch_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}

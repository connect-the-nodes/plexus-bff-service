####################################
# ECS TASK EXECUTION ROLE
####################################
resource "aws_iam_role" "ecs_task_execution" {
  name = "zynchub-ecs-task-execution-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

####################################
# ECS TASK ROLE
####################################
resource "aws_iam_role" "ecs_task_role" {
  name = "zynchub-ecs-task-role-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

####################################
# ECS EC2 INSTANCE ROLE
####################################
resource "aws_iam_role" "ecs_instance_role" {
  name = "zynchub-ecs-ec2-role-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_instance_policy" {
  role       = aws_iam_role.ecs_instance_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_instance_profile" "ecs_instance_profile" {
  name = "zynchub-ecs-instance-profile-${var.environment}"
  role = aws_iam_role.ecs_instance_role.name
}

resource "aws_iam_policy" "appconfig_features_read" {
  name        = "zynchub-appconfig-features-${var.environment}"
  description = "Read feature flags from AWS AppConfig"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "appconfigdata:StartConfigurationSession",
          "appconfigdata:GetLatestConfiguration"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "appconfig_features_read" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.appconfig_features_read.arn
}

####################################
# API GATEWAY CLOUDWATCH LOGS
####################################
resource "aws_iam_role" "api_gw_cloudwatch" {
  count = var.create_apigw_cloudwatch_role ? 1 : 0
  name  = local.apigw_role_name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "apigateway.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

data "aws_iam_role" "api_gw_cloudwatch" {
  count = var.create_apigw_cloudwatch_role ? 0 : 1
  name  = local.apigw_role_name
}

resource "aws_iam_role_policy_attachment" "api_gw_logs_policy" {
  count      = var.create_apigw_cloudwatch_role ? 1 : 0
  role       = aws_iam_role.api_gw_cloudwatch[0].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}

locals {
  apigw_role_arn = var.create_apigw_cloudwatch_role ? aws_iam_role.api_gw_cloudwatch[0].arn : data.aws_iam_role.api_gw_cloudwatch[0].arn
}

resource "aws_api_gateway_account" "this" {
  cloudwatch_role_arn = local.apigw_role_arn
}

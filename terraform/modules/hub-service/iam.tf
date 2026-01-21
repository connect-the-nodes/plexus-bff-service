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

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_policy" "ecs_ssm_parameters" {
  name        = "zynchub-ecs-ssm-${var.environment}"
  description = "Allow ECS task execution to read SSM parameters"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameters",
          "ssm:GetParameter",
          "ssm:GetParametersByPath"
        ]
        Resource = [
          aws_ssm_parameter.auth_cognito_domain.arn,
          aws_ssm_parameter.auth_cognito_client_id.arn,
          aws_ssm_parameter.auth_cognito_redirect_uri.arn,
          aws_ssm_parameter.auth_cognito_post_login_redirect_uri.arn,
          aws_ssm_parameter.security_jwt_issuer_uri.arn,
          aws_ssm_parameter.security_jwt_jwk_set_uri.arn
        ]
      }
    ]
  })
  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "ecs_ssm_parameters" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = aws_iam_policy.ecs_ssm_parameters.arn
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

  tags = var.tags
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

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "ecs_instance_policy" {
  role       = aws_iam_role.ecs_instance_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_instance_profile" "ecs_instance_profile" {
  name = "zynchub-ecs-instance-profile-${var.environment}"
  role = aws_iam_role.ecs_instance_role.name
  tags = var.tags
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
          "appconfig:StartConfigurationSession",
          "appconfig:GetLatestConfiguration",
          "appconfigdata:StartConfigurationSession",
          "appconfigdata:GetLatestConfiguration"
        ]
        Resource = "*"
      }
    ]
  })
  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "appconfig_features_read" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.appconfig_features_read.arn
}

resource "aws_iam_policy" "redis_connect" {
  name        = "zynchub-redis-connect-${var.environment}"
  description = "Allow ECS task to connect to Redis"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "elasticache:Connect"
        ]
        Resource = [
          aws_elasticache_replication_group.redis.arn,
          aws_elasticache_user.redis_iam.arn
        ]
      }
    ]
  })
  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "redis_connect" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.redis_connect.arn
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

  tags = var.tags
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

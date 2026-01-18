############################################
# ECS Optimized AMI (Amazon Linux 2 ECS)
############################################
data "aws_ami" "ecs_optimized" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-hvm-*-x86_64-ebs"]
  }
}

############################################
# ECS Cluster
############################################
resource "aws_ecs_cluster" "main" {
  name = "zynchub-cluster-${var.environment}"
  tags = var.tags
}

############################################
# ECS Task Definition (EC2 launch type)
############################################
resource "aws_ecs_task_definition" "app" {
  family                   = "zynchub-task-${var.environment}"
  requires_compatibilities = ["EC2"]
  network_mode             = "bridge"

  cpu    = "256"
  memory = "512"

  task_role_arn      = aws_iam_role.ecs_task_role.arn
  execution_role_arn = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([
    {
      name      = "zynchub-app"
      image     = "${var.ecr_repository_url}:${var.app_version}"
      essential = true

      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "SPRING_PROFILES_ACTIVE", value = var.environment }
        , { name = "AWS_APP_CONFIG_FEATURES_APPLICATION_ID", value = aws_appconfig_application.features.id }
        , { name = "AWS_APP_CONFIG_FEATURES_ENVIRONMENT_ID", value = aws_appconfig_environment.features.environment_id }
        , { name = "AWS_APP_CONFIG_FEATURES_CONFIGURATION_ID", value = aws_appconfig_configuration_profile.features.configuration_profile_id }
        , { name = "SPRING_DATA_REDIS_HOST", value = aws_elasticache_replication_group.redis.primary_endpoint_address }
        , { name = "SPRING_DATA_REDIS_PORT", value = tostring(aws_elasticache_replication_group.redis.port) }
        , { name = "SPRING_DATA_REDIS_SSL_ENABLED", value = "true" }
        , { name = "SPRING_DATA_REDIS_IAM_ENABLED", value = "true" }
        , { name = "SPRING_DATA_REDIS_USERID", value = aws_elasticache_user.redis_iam.user_id }
        , { name = "SPRING_DATA_REDIS_REPLICATIONGROUPID", value = aws_elasticache_replication_group.redis.replication_group_id }
        , { name = "SPRING_DATA_REDIS_REGION", value = var.aws_region }
        , { name = "SECURITY_ENABLED", value = var.environment == "local" ? "false" : "true" }
        , { name = "SECURITY_JWT_ISSUER_URI", value = local.cognito_issuer }
        , { name = "SECURITY_JWT_JWK_SET_URI", value = "${local.cognito_issuer}/.well-known/jwks.json" }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.ecs.name
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  tags = var.tags
}

############################################
# Launch Template (EC2 capacity for ECS)
############################################
resource "aws_launch_template" "ecs" {
  name_prefix   = "zynchub-ecs-${var.environment}-"
  image_id      = data.aws_ami.ecs_optimized.id
  instance_type = "t3.small"

  iam_instance_profile {
    name = aws_iam_instance_profile.ecs_instance_profile.name
  }

  vpc_security_group_ids = [aws_security_group.ecs_tasks.id]

  user_data = base64encode(<<EOF
#!/bin/bash
echo ECS_CLUSTER=${aws_ecs_cluster.main.name} >> /etc/ecs/ecs.config
EOF
  )

  lifecycle {
    create_before_destroy = true
  }

  tag_specifications {
    resource_type = "instance"
    tags          = var.tags
  }

  tag_specifications {
    resource_type = "volume"
    tags          = var.tags
  }
}

############################################
# Auto Scaling Group (REQUIRED for EC2 ECS)
############################################
resource "aws_autoscaling_group" "ecs" {
  name = "zynchub-asg-${var.environment}"

  min_size         = 2
  max_size         = 2
  desired_capacity = 2

  launch_template {
    id      = aws_launch_template.ecs.id
    version = "$Latest"
  }

  # Free-tier friendly: put instances in PUBLIC subnets so no NAT needed
  vpc_zone_identifier = aws_subnet.public[*].id

  health_check_type = "EC2"

  tag {
    key                 = "Name"
    value               = "zynchub-ecs-instance-${var.environment}"
    propagate_at_launch = true
  }

  tag {
    key                 = "Environment"
    value               = var.environment
    propagate_at_launch = true
  }

  tag {
    key                 = "Service"
    value               = var.service_name
    propagate_at_launch = true
  }
}

############################################
# ECS Service
############################################
resource "aws_ecs_service" "main" {
  name            = "zynchub-service-${var.environment}"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = 1
  launch_type     = "EC2"

  load_balancer {
    target_group_arn = aws_lb_target_group.app.arn
    container_name   = "zynchub-app"
    container_port   = 8080
  }

  depends_on = [
    aws_autoscaling_group.ecs,
    aws_lb_listener.app
  ]

  tags = var.tags
}

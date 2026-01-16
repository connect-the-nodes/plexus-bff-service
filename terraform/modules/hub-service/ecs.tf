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
}

############################################
# ECS Task Definition (EC2 launch type)
############################################
resource "aws_ecs_task_definition" "app" {
  family                   = "zynchub-task-${var.environment}"
  requires_compatibilities = ["EC2"]
  network_mode             = "awsvpc"

  cpu    = "512"
  memory = "800"

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
}

############################################
# Launch Template (EC2 capacity for ECS)
############################################
resource "aws_launch_template" "ecs" {
  name_prefix   = "zynchub-ecs-${var.environment}-"
  image_id      = data.aws_ami.ecs_optimized.id
  instance_type = "t3.micro"

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
}

############################################
# Auto Scaling Group (REQUIRED for EC2 ECS)
############################################
resource "aws_autoscaling_group" "ecs" {
  name = "zynchub-asg-${var.environment}"

  min_size         = 1
  max_size         = 1
  desired_capacity = 1

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

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.app.arn
    container_name   = "zynchub-app"
    container_port   = 8080
  }

  depends_on = [
    aws_autoscaling_group.ecs,
    aws_lb_listener.app
  ]
}

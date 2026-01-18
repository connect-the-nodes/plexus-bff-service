######################################
# REDIS (ELASTICACHE)
######################################

resource "aws_security_group" "redis" {
  name        = "zynchub-redis-${var.environment}"
  description = "Redis access from ECS"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_tasks.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.tags
}

resource "aws_elasticache_subnet_group" "redis" {
  name       = "zynchub-redis-${var.environment}"
  subnet_ids = aws_subnet.public[*].id

  tags = var.tags
}

resource "aws_elasticache_user" "redis_iam" {
  user_id       = "zynchub-${var.environment}-redis"
  user_name     = "zynchub-${var.environment}-redis"
  engine        = "redis"
  access_string = "on ~* +@all"

  authentication_mode {
    type = "iam"
  }

  tags = var.tags
}

resource "aws_elasticache_user_group" "redis" {
  user_group_id = "zynchub-${var.environment}-redis"
  engine        = "redis"
  user_ids      = ["default", aws_elasticache_user.redis_iam.user_id]

  tags = var.tags
}

resource "aws_elasticache_replication_group" "redis" {
  replication_group_id          = "zynchub-${var.environment}-redis"
  description                   = "Zynchub Redis session store"
  engine                        = "redis"
  engine_version                = "7.1"
  node_type                     = "cache.t4g.micro"
  num_cache_clusters            = 1
  port                          = 6379
  automatic_failover_enabled    = false
  multi_az_enabled              = false
  at_rest_encryption_enabled    = true
  transit_encryption_enabled    = true
  apply_immediately             = true
  parameter_group_name          = "default.redis7"
  security_group_ids            = [aws_security_group.redis.id]
  subnet_group_name             = aws_elasticache_subnet_group.redis.name
  user_group_ids                = [aws_elasticache_user_group.redis.id]

  tags = var.tags
}

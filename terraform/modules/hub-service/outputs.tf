
output "ecs_task_execution_role_arn" {
  value = aws_iam_role.ecs_task_execution.arn
}

output "ecs_task_role_arn" {
  value = aws_iam_role.ecs_task_role.arn
}

output "ecs_instance_profile" {
  value = aws_iam_instance_profile.ecs_instance_profile.name
}

output "api_gw_cloudwatch_role_arn" {
  value = local.apigw_role_arn
}


output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnets" {
  value = aws_subnet.public[*].id
}

output "ecs_sg_id" {
  value = aws_security_group.ecs_tasks.id
}

output "nlb_arn" {
  value = aws_lb.nlb.arn
}

output "nlb_dns" {
  value = aws_lb.nlb.dns_name
}

output "nlb_target_group_arn" {
  value = aws_lb_target_group.app.arn
}

output "nlb_dns_name" {
  description = "DNS name of the Network Load Balancer"
  value       = aws_lb.nlb.dns_name
}

output "vpc_link_id" {
  description = "API Gateway VPC Link ID"
  value       = aws_api_gateway_vpc_link.main.id
}

output "redis_primary_endpoint" {
  description = "Redis primary endpoint address"
  value       = aws_elasticache_replication_group.redis.primary_endpoint_address
}

output "redis_user_id" {
  description = "Redis IAM user ID"
  value       = aws_elasticache_user.redis_iam.user_id
}

output "cognito_user_pool_id" {
  description = "Cognito user pool ID"
  value       = aws_cognito_user_pool.main.id
}

output "cognito_user_pool_client_id" {
  description = "Cognito user pool client ID"
  value       = aws_cognito_user_pool_client.app.id
}

output "cognito_domain" {
  description = "Cognito hosted UI domain"
  value       = aws_cognito_user_pool_domain.main.domain
}

output "cognito_issuer" {
  description = "Cognito issuer URL"
  value       = local.cognito_issuer
}

output "ssm_auth_cognito_domain" {
  description = "SSM parameter for Cognito domain"
  value       = aws_ssm_parameter.auth_cognito_domain.name
}

output "ssm_auth_cognito_client_id" {
  description = "SSM parameter for Cognito client ID"
  value       = aws_ssm_parameter.auth_cognito_client_id.name
}

output "ssm_security_jwt_issuer_uri" {
  description = "SSM parameter for JWT issuer URI"
  value       = aws_ssm_parameter.security_jwt_issuer_uri.name
}

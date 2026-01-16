
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

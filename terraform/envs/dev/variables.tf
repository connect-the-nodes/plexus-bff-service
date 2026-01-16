variable "environment" {
  description = "Deployment environment (dev, prod)"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "app_version" {
  description = "Docker image tag pushed to ECR"
  type        = string
}

variable "ecr_repository_url" {
  description = "ECR repository URL"
  type        = string
}

variable "tags" {
  description = "Common resource tags"
  type        = map(string)
}

variable "create_apigw_log_group" {
  description = "Whether to create the API Gateway log group"
  type        = bool
}

variable "create_apigw_cloudwatch_role" {
  description = "Whether to create the API Gateway CloudWatch role"
  type        = bool
}



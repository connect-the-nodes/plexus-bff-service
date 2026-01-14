variable "environment" {
  type = string
}

variable "aws_region" {
  type = string
}

variable "service_name" {
  type = string
}

variable "app_version" {
  type = string
}

variable "ecr_repository_url" {
  type = string
}

variable "openapi_spec_path" {
  description = "Path to the OpenAPI YAML spec"
  type        = string
}

variable "nlb_dns" {
  description = "DNS name of the Network Load Balancer"
  type        = string
}

variable "vpc_link_id" {
  description = "API Gateway VPC Link ID"
  type        = string
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Resource tags for cost tracking and management"
}

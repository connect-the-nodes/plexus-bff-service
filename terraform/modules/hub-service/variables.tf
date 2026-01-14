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

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Resource tags for cost tracking and management"
}

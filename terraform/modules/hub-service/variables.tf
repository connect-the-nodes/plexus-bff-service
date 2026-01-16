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

variable "features_config_path" {
  description = "Path to the features YAML for AWS AppConfig"
  type        = string
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Resource tags for cost tracking and management"
}

variable "create_apigw_log_group" {
  type        = bool
  default     = true
  description = "Whether to create the API Gateway log group in this module"
}

variable "create_apigw_cloudwatch_role" {
  type        = bool
  default     = true
  description = "Whether to create the API Gateway CloudWatch role in this module"
}

variable "environment" {
  type        = string
  description = "Deployment environment (dev, stage, prod)"
}

variable "aws_region" {
  type        = string
  description = "AWS region for deployment"
}

variable "jar_path" {
  type        = string
  description = "Path to the compiled Spring Boot JAR file"
}

variable "app_name" {
  type        = string
  default     = "zynchub"
  description = "Application name prefix for resources"
}

variable "lambda_memory" {
  type        = number
  default     = 2048
  description = "Memory allocated to the Lambda function"
}

variable "lambda_timeout" {
  type        = number
  default     = 30
  description = "Lambda execution timeout in seconds"
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Resource tags for cost tracking and management"
}


variable "app_version" {
  type        = string
  description = "The version of the application extracted from pom.xml"
}

variable "service_name" {
  type        = string
  default     = "zynchub-service"
  description = "The name used for resource grouping and log groups"
}
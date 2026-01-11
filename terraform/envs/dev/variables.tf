variable "environment" {
  type        = string
  description = "Target environment name"
}

variable "aws_region" {
  type        = string
  description = "AWS Region"
}

# Note: jar_path is usually handled via local path in main.tf
# but we declare it here if we want to override it from tfvars
variable "jar_path" {
  type        = string
  default     = "./zynchub-digital-hub-service.jar"
  description = "The JAR file name copied by the GitHub Action"
}

variable "app_version" {
  type        = string
  description = "The version of the application extracted from pom.xml"
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Resource tags"
}
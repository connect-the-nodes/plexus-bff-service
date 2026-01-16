terraform {
  required_version = ">= 1.6.0"
}

data "aws_caller_identity" "current" {}

module "hub_service" {
  source = "../../modules/hub-service"

  # Environment & Provider
  environment = var.environment
  aws_region  = var.aws_region

  # Application
  service_name       = "zynchub"
  app_version        = var.app_version
  ecr_repository_url = var.ecr_repository_url

  # OpenAPI spec path (module renders it)
  openapi_spec_path = "${path.root}/../../../src/main/resources/static/openapi.yaml"

  # Tags
  tags = var.tags

  create_apigw_log_group       = var.create_apigw_log_group
  create_apigw_cloudwatch_role = var.create_apigw_cloudwatch_role
}

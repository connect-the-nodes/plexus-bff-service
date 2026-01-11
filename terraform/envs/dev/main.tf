data "aws_caller_identity" "current" {}

locals {
  # Manually construct ARN to prevent circular dependency
  lambda_name = "zynchub-service-${var.environment}"
  lambda_arn  = "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${local.lambda_name}"
}

module "hub_service" {
  source       = "../../modules/hub-service"
  environment  = var.environment
  aws_region   = var.aws_region
  app_version  = var.app_version
  service_name = local.lambda_name

  # Pass rendered spec into the module
  openapi_spec = templatefile("${path.root}/../../../src/main/resources/static/openapi.yaml", {
    lambda_arn = local.lambda_arn
    region     = var.aws_region
  })
}
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
  jar_path     = var.jar_path
  artifact_bucket = aws_s3_bucket.lambda_artifacts.id
  artifact_key    = "hub-service/app.jar"
  # Pass rendered spec into the module
  openapi_spec = templatefile("${path.root}/../../../src/main/resources/static/openapi.yaml", {
    lambda_arn = local.lambda_arn
    region     = var.aws_region
  })
}

# Create the bucket for your Lambda JARs
resource "aws_s3_bucket" "lambda_artifacts" {
  bucket = "zynchub-artifacts-${data.aws_caller_identity.current.account_id}-${var.environment}"

  # Prevent accidental deletion of the bucket
  lifecycle {
    prevent_destroy = false
  }
}

# Block public access (Enterprise Security Best Practice)
resource "aws_s3_bucket_public_access_block" "artifacts_protection" {
  bucket                  = aws_s3_bucket.lambda_artifacts.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
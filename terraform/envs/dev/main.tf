module "hub_service" {
  source      = "../../modules/hub-service"
  environment = var.environment
  aws_region  = var.aws_region
  jar_path    = "${path.module}/zynchub-digital-hub-service.jar"
}
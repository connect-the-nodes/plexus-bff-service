############################################
# AWS AppConfig for Feature Flags
############################################
resource "aws_appconfig_deployment_strategy" "all_at_once" {
  name                           = "AllAtOnce"
  deployment_duration_in_minutes = 0
  growth_factor                  = 100
  growth_type                    = "LINEAR"
  final_bake_time_in_minutes     = 0
  replicate_to                   = "NONE"
  description                    = "All at once rollout"
}

resource "aws_appconfig_application" "features" {
  name        = "zynchub-features-${var.environment}"
  description = "Zynchub feature flags configuration"
  tags        = var.tags
}

resource "aws_appconfig_environment" "features" {
  application_id = aws_appconfig_application.features.id
  name           = var.environment
  description    = "Zynchub ${var.environment} feature flags environment"
  tags           = var.tags
}

resource "aws_appconfig_configuration_profile" "features" {
  application_id = aws_appconfig_application.features.id
  name           = "features"
  location_uri   = "hosted"
  type           = "AWS.Freeform"
  description    = "Feature flags configuration"
  tags           = var.tags
}

resource "aws_appconfig_hosted_configuration_version" "features" {
  application_id           = aws_appconfig_application.features.id
  configuration_profile_id = aws_appconfig_configuration_profile.features.configuration_profile_id
  content_type             = "text/yaml"
  content                  = file(var.features_config_path)
}

resource "aws_appconfig_deployment" "features" {
  application_id           = aws_appconfig_application.features.id
  environment_id           = aws_appconfig_environment.features.environment_id
  configuration_profile_id = aws_appconfig_configuration_profile.features.configuration_profile_id
  configuration_version    = aws_appconfig_hosted_configuration_version.features.version_number
  deployment_strategy_id   = aws_appconfig_deployment_strategy.all_at_once.id
  description              = "Deploy feature flags"
  tags                     = var.tags
}

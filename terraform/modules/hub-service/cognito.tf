######################################
# COGNITO (USER POOL + HOSTED UI)
######################################

locals {
  cognito_domain_prefix = coalesce(var.cognito_domain_prefix, "zynchub-${var.environment}")
  cognito_issuer        = "https://cognito-idp.${var.aws_region}.amazonaws.com/${aws_cognito_user_pool.main.id}"
}

resource "aws_cognito_user_pool" "main" {
  name = "zynchub-${var.environment}"

  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length    = 8
    require_uppercase = true
    require_lowercase = true
    require_numbers   = true
    require_symbols   = false
  }

  tags = var.tags
}

resource "aws_cognito_user_pool_client" "app" {
  name         = "zynchub-ui-${var.environment}"
  user_pool_id = aws_cognito_user_pool.main.id

  generate_secret                      = false
  supported_identity_providers         = ["COGNITO"]
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = var.cognito_allowed_oauth_scopes
  callback_urls                        = var.cognito_callback_urls
  logout_urls                          = var.cognito_logout_urls

  explicit_auth_flows = [
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_SRP_AUTH"
  ]

  prevent_user_existence_errors = "ENABLED"
}

resource "aws_cognito_user_pool_domain" "main" {
  domain       = local.cognito_domain_prefix
  user_pool_id = aws_cognito_user_pool.main.id
}

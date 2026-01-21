######################################
# SSM PARAMETERS (AUTH/JWT)
######################################

locals {
  auth_cognito_domain             = "https://${aws_cognito_user_pool_domain.main.domain}.auth.${var.aws_region}.amazoncognito.com"
  auth_cognito_redirect_uri       = var.cognito_callback_urls[0]
  auth_post_login_redirect_uri    = coalesce(var.auth_post_login_redirect_uri, var.cognito_callback_urls[0])
  security_jwt_issuer_uri_value   = local.cognito_issuer
  security_jwt_jwk_set_uri_value  = "${local.cognito_issuer}/.well-known/jwks.json"
}

resource "aws_ssm_parameter" "auth_cognito_domain" {
  name  = "/zynchub/${var.environment}/auth/cognito/domain"
  type  = "String"
  value = local.auth_cognito_domain
  tags  = var.tags
}

resource "aws_ssm_parameter" "auth_cognito_client_id" {
  name  = "/zynchub/${var.environment}/auth/cognito/client-id"
  type  = "String"
  value = aws_cognito_user_pool_client.app.id
  tags  = var.tags
}

resource "aws_ssm_parameter" "auth_cognito_redirect_uri" {
  name  = "/zynchub/${var.environment}/auth/cognito/redirect-uri"
  type  = "String"
  value = local.auth_cognito_redirect_uri
  tags  = var.tags
}

resource "aws_ssm_parameter" "auth_cognito_post_login_redirect_uri" {
  name  = "/zynchub/${var.environment}/auth/cognito/post-login-redirect-uri"
  type  = "String"
  value = local.auth_post_login_redirect_uri
  tags  = var.tags
}

resource "aws_ssm_parameter" "security_jwt_issuer_uri" {
  name  = "/zynchub/${var.environment}/security/jwt/issuer-uri"
  type  = "String"
  value = local.security_jwt_issuer_uri_value
  tags  = var.tags
}

resource "aws_ssm_parameter" "security_jwt_jwk_set_uri" {
  name  = "/zynchub/${var.environment}/security/jwt/jwk-set-uri"
  type  = "String"
  value = local.security_jwt_jwk_set_uri_value
  tags  = var.tags
}

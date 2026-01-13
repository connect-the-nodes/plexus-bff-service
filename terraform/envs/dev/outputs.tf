output "deployment_api_endpoint" {
  description = "Access the service at this URL"
  value       = module.hub_service.api_url
}

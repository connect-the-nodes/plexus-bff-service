environment = "dev"
aws_region  = "eu-west-1"

# This is injected dynamically by GitHub Actions.
# For manual applies, "latest" matches the moving tag pushed by the workflow.
app_version = "latest"

# Example:
# 410521973628.dkr.ecr.eu-west-1.amazonaws.com/zynchub-repo-dev
ecr_repository_url = "410521973628.dkr.ecr.eu-west-1.amazonaws.com/zynchub-repo-dev"

tags = {
  Project     = "Zynchub"
  Owner       = "Fiaz Zeeshan Zackariya"
  Environment = "dev"
}

create_apigw_log_group = true
create_apigw_cloudwatch_role = true

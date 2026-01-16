environment = "dev"
aws_region  = "eu-west-1"

# This is injected dynamically by GitHub Actions
# Example: 9a3f7c2b8a1e
app_version = "bff7bddfed8939c16e107eaaf8e958702aa37c81"

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

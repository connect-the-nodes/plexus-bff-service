terraform {
  backend "s3" {
    bucket         = "zynchub-hub-terraform-state"
    key            = "hub-service/dev/terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "zynchub-terraform-locks"
    encrypt        = true
  }
}
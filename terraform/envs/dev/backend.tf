terraform {
  backend "s3" {
    bucket         = "zynchub-hub-terraform-state"
    key            = "hub-service/dev/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "zynchub-terraform-locks"
    encrypt        = true
  }
}
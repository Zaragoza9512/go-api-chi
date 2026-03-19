# Backend: Configuracion de Remote State en S3

terraform {
  backend "s3" {
    bucket = "terraform-state-zaragoza"
    key = "infrastructure/terraform.tfstate"
    region = "us-east-1"
    dynamodb_table = "terraform-state-locks"
    encrypt = true
  }
}
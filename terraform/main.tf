# Provider: Configurar AWS
provider "aws" {
  region = var.aws_region
}

# Data Source: Buscar la AMI más reciente de Amazon Linux 2023
data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# Resource: Instancia EC2
resource "aws_instance" "api_server" {
  ami           = data.aws_ami.amazon_linux_2023.id
  instance_type = var.instance_type
  key_name      = var.key_pair_name
  subnet_id     = var.subnet_id
  
  vpc_security_group_ids = [var.security_group_id]

  tags = {
    Name        = "go-api-server"
    Environment = "Production"
    ManagedBy   = "Terraform"
    Project     = "go-api-chi"
  }
}

# Resource: Instancia RDS PostgreSQL
resource "aws_db_instance" "postgres" {
  identifier          = var.db_identifier
  engine              = var.db_engine
  engine_version      = var.db_engine_version
  instance_class      = var.db_instance_class
  allocated_storage   = var.db_allocated_storage
  db_name             = var.db_name
  username            = var.db_username
  password            = var.db_password
  
  publicly_accessible = false
  skip_final_snapshot = true
  storage_type        = "gp3"

  tags = {
    Name        = "go-api-postgres"
    Environment = "Production"
    ManagedBy   = "Terraform"
    Project     = "go-api-chi"
  }
}
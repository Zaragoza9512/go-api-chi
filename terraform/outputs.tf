output "ec2_public_ip" {
  description = "Ip publica de la instancia EC2"
  value = aws_instance.api_server.public_ip
}

output "ec2_instance_id" {
  description = "ID de la instancia EC2"
  value = aws_instance.api_server.id
}

output "rds_endpoint" {
  description = "Endpoint de concexion a la base de datos"
  value = aws_db_instance.postgres.endpoint 
}

output "rds_port" {
  description = "Puerto de la base de datos"
  value = aws_db_instance.postgres.port
}
#-General
variable "aws_region" {
  description = "Region de AWS donde se crearan los recursos"
  type = string
  default = "us-east-1"
}

#-EC2
variable "instance_type" {
  description = "Tipo de instancia EC2"
  type = string
  default = "t3.micro"
}

variable "key_pair_name" {
  description = "Nombre de la Key Pair para el acceso SSH"
  type = string
  default = "go-api-key"
}

variable "subnet_id" {
  description = "ID de la subnet para la instancia EC2"
  type        = string
  default     = "subnet-0e0fe92f8f2613dd1"
}

variable "security_group_id" {
  description = "ID del Security Group para la instancia EC2"
  type        = string
  default     = "sg-0b8ac9ea84059a59a"
}

#-RDS
variable "db_identifier" {
  description = "Identificador único de la instancia RDS"
  type        = string
  default     = "go-api-postgres-db"
}

variable "db_engine" {
  description = "Motor de base de datos"
  type        = string
  default     = "postgres"
}

variable "db_engine_version" {
  description = "Versión del motor de base de datos"
  type        = string
  default     = "17.2"
}

variable "db_instance_class" {
  description = "Clase de instancia RDS (tamaño de CPU/RAM)"
  type        = string
  default     = "db.t4g.micro"
}

variable "db_allocated_storage" {
  description = "Almacenamiento en GB para la base de datos"
  type        = number
  default     = 20
}

variable "db_name" {
  description = "Nombre de la base de datos inicial"
  type        = string
  default     = "appdb"
}

variable "db_username" {
  description = "Usuario administrador de la base de datos"
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Contraseña del usuario administrador"
  type        = string
  sensitive   = true
  # Sin default - debe proporcionarse
}

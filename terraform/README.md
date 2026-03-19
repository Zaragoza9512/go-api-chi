# Infraestructura go-api-chi

Infraestructura como código para el proyecto go-api-chi usando Terraform.

## Recursos Creados

- **EC2**: Instancia t3.micro para la API
- **RDS**: PostgreSQL 17.2 (db.t4g.micro)
- **Remote State**: S3 + DynamoDB para gestión de estado

## Requisitos

- Terraform >= 1.14
- AWS CLI configurado
- Bucket S3: `terraform-state-zaragoza`
- Tabla DynamoDB: `terraform-state-locks`

## Uso
```bash
# Inicializar
terraform init

# Ver plan
terraform plan

# Aplicar cambios
terraform apply

# Destruir todo
terraform destroy
```

## Variables

Ver `variables.tf` para la lista completa de variables configurables.

La contraseña de la base de datos debe proporcionarse en `terraform.tfvars`.
# 🔧 Troubleshooting Guide

Guía de solución de problemas comunes.

---

## Tabla de Contenidos

- [Docker](#docker)
- [API](#api)
- [Base de Datos](#base-de-datos)
- [Terraform](#terraform)
- [CI/CD](#cicd)

---

## Docker

### Error: "Cannot connect to Docker daemon"

**Síntoma:**

#Cannot connect to the Docker daemon at unix:///var/run/docker.sock
**Causa:** Docker no está corriendo o no tienes permisos.

**Solución:**
```bash
# Ubuntu/WSL2
sudo service docker start

# Verificar estado
sudo service docker status

# Agregar usuario a grupo docker (evita usar sudo)
sudo usermod -aG docker $USER
newgrp docker
```

---

### Error: "Port already in use"

**Síntoma:**

#ERROR: for api  Cannot start service api: driver failed programming external connectivity:
#Bind for 0.0.0.0:8080 failed: port is already allocated
**Causa:** Otro proceso usa el puerto 8080.

**Solución:**

**Opción 1:** Detener el proceso que usa el puerto
```bash
# Ver qué usa el puerto 8080
sudo lsof -i :8080

# Matar el proceso (reemplaza PID)
kill -9 PID
```

**Opción 2:** Cambiar puerto en docker-compose.yml
```yaml
api:
  ports:
    - "8081:8080"  # Cambiar de 8080 a 8081
```

---

### Error: "no space left on device"

**Síntoma:**

#ERROR: no space left on device
**Causa:** Docker acumuló muchas imágenes/contenedores/volúmenes.

**Solución:**
```bash
# Ver uso de espacio
docker system df

# Limpiar todo lo no usado
docker system prune -a --volumes

# Limpiar solo imágenes
docker image prune -a

# Limpiar solo volúmenes
docker volume prune
```

---

### Contenedor se reinicia constantemente

**Síntoma:**
```bash
docker ps
# Status: Restarting (1) Less than a second ago
```

**Diagnóstico:**
```bash
# Ver logs del contenedor
docker logs go_api_container

# Ver logs en tiempo real
docker logs -f go_api_container
```

**Causas comunes:**
1. Error en la aplicación (revisar logs)
2. Base de datos no disponible (esperar a que RDS esté lista)
3. Variables de entorno incorrectas

---

## API

### Error: "Token inválido o expirado"

**Síntoma:**
```json
{
  "error": "Token inválido o expirado"
}
```

**Causa:** Token JWT expirado (24 horas) o mal formado.

**Solución:**
```bash
# Obtener nuevo token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Verificar que el token sea válido
# Debe empezar con: eyJ...
```

---

### Error: "Credenciales inválidas"

**Síntoma:**
```json
{
  "error": "Credenciales inválidas"
}
```

**Causa:** Usuario o contraseña incorrectos.

**Solución:**

Usuario y contraseña actuales (hardcodeados en `security.go`):

#Username: admin
#Password: password
**Nota:** En versión futura se implementará tabla de usuarios.

---

### API no responde (timeout)

**Síntoma:**
```bash
curl: (28) Failed to connect to localhost port 8080: Connection timed out
```

**Diagnóstico:**
```bash
# ¿El contenedor está corriendo?
docker ps | grep go_api

# ¿Escucha en el puerto correcto?
docker exec go_api_container netstat -tuln | grep 8080

# Ver logs
docker logs go_api_container
```

**Soluciones:**
1. Verificar que docker-compose esté corriendo
2. Revisar firewall
3. Verificar que el puerto esté mapeado correctamente

---

### Error: "database connection failed"

**Síntoma:**

#Error connecting to database: connection refused
**Causa:** PostgreSQL no está disponible.

**Solución:**
```bash
# Verificar que PostgreSQL está corriendo
docker ps | grep postgres

# Ver logs de PostgreSQL
docker logs go_db_container

# Verificar conectividad desde API
docker exec -it go_api_container ping postgres

# Probar conexión manual
docker exec -it go_db_container psql -U postgres -d ecom_db
```

---

## Base de Datos

### Error: "database does not exist"

**Síntoma:**

#ERROR: database "ecom_db" does not exist
**Causa:** La base de datos no se inicializó correctamente.

**Solución:**
```bash
# Recrear contenedores con volúmenes limpios
docker-compose down -v
docker-compose up --build
```

---

### Olvidé la contraseña de PostgreSQL

**Solución:**
```bash
# La contraseña está en .env
cat .env | grep POSTGRES_PASSWORD

# O usa el default del docker-compose.yml
# Password: postgres
```

---

### Datos de prueba no se cargaron

**Causa:** El script `init.sql` no se ejecutó.

**Solución:**
```bash
# Verificar que init.sql existe
ls -la init.sql

# Recrear con volúmenes limpios
docker-compose down -v
docker-compose up --build

# Verificar datos manualmente
docker exec -it go_db_container psql -U postgres -d ecom_db -c "SELECT * FROM productos;"
```

---

## Terraform

### Error: "timeout while waiting for plugin"

**Síntoma:**

#Error: timeout while waiting for plugin to start
**Causa:** Plugin de AWS tarda en iniciar (común en WSL2).

**Solución:**
```bash
# Aumentar timeout
TF_PLUGIN_TIMEOUT=60 terraform plan
TF_PLUGIN_TIMEOUT=60 terraform apply

# Limpiar caché si persiste
rm -rf .terraform .terraform.lock.hcl
terraform init
```

---

### Error: "Backend initialization required"

**Síntoma:**

#Error: Backend initialization required
**Causa:** Cambio en configuración de backend sin reinicializar.

**Solución:**
```bash
terraform init -reconfigure
```

---

### Error: "NoSuchBucket"

**Síntoma:**

#Error: Failed to get existing workspaces: S3 bucket does not exist
**Causa:** El bucket S3 para remote state no existe.

**Solución:**
```bash
# Crear la infraestructura de remote state primero
cd ~/terraform-remote-state-setup
terraform apply

# Luego volver al proyecto principal
cd ~/terraform-infrastructure
terraform init
```

---

### Error: "invalid AWS Region"

**Síntoma:**

#Error: invalid AWS Region: us_east_1
**Causa:** Región con guión bajo en lugar de guión normal.

**Solución:**
```hcl
# Incorrecto
default = "us_east_1"

# Correcto
default = "us-east-1"
```

---

### Error: "Error acquiring state lock"

**Síntoma:**

#Error: Error acquiring the state lock
**Causa:** Otra persona/proceso tiene el lock, o lock quedó huérfano.

**Diagnóstico:**
```bash
# Ver el lock en DynamoDB
aws dynamodb scan --table-name terraform-state-locks
```

**Solución:**

**Si es lock legítimo:** Esperar a que termine.

**Si es lock huérfano:** Forzar desbloqueo
```bash
terraform force-unlock LOCK_ID
```

⚠️ **CUIDADO:** Solo usar si estás seguro que nadie más está ejecutando Terraform.

---

## CI/CD

### GitHub Actions falla en build

**Síntoma:** Workflow falla con error de compilación.

**Diagnóstico:**
```bash
# Ver logs en GitHub
# Actions → Workflow → Failed job → Ver detalles
```

**Causas comunes:**
1. Tests fallan localmente (correr `go test ./...`)
2. Dependencias faltantes (actualizar `go.mod`)
3. Errores de sintaxis

**Solución:**
```bash
# Probar localmente primero
go build .
go test ./...

# Si pasa local, revisar configuración de GitHub Actions
```

---

### Docker Hub push falla

**Síntoma:**

#denied: requested access to the resource is denied
**Causa:** Credenciales de Docker Hub incorrectas o expiradas.

**Solución:**
1. Ir a GitHub → Settings → Secrets
2. Verificar `DOCKER_USERNAME` y `DOCKER_PASSWORD`
3. Regenerar token en Docker Hub si es necesario

---

## Comandos Útiles de Diagnóstico

### Docker

```bash
# Ver todos los contenedores
docker ps -a

# Ver uso de recursos
docker stats

# Inspeccionar contenedor
docker inspect go_api_container

# Ver redes
docker network ls
docker network inspect go-api-chi_default

# Entrar al contenedor
docker exec -it go_api_container sh
```

### Logs

```bash
# API
docker logs -f go_api_container

# PostgreSQL
docker logs -f go_db_container

# Últimas 100 líneas
docker logs --tail 100 go_api_container

# Con timestamps
docker logs -t go_api_container
```

### Base de Datos

```bash
# Conectar a PostgreSQL
docker exec -it go_db_container psql -U postgres -d ecom_db

# Ver tablas
\dt

# Ver datos
SELECT * FROM productos;

# Salir
\q
```

### Terraform

```bash
# Ver state actual
terraform show

# Listar recursos
terraform state list

# Ver detalles de recurso
terraform state show aws_instance.api_server

# Ver outputs
terraform output
```

---

## Checklist de Diagnóstico

Cuando algo no funciona, revisar en orden:

### 1. Docker
- [ ] Docker está corriendo
- [ ] Contenedores están UP
- [ ] Puertos mapeados correctamente
- [ ] Sin errores en logs

### 2. Variables de Entorno
- [ ] Archivo `.env` existe
- [ ] Variables correctas
- [ ] Sin espacios extra

### 3. Red
- [ ] Puertos no están en uso
- [ ] Firewall permite conexiones
- [ ] Contenedores en misma red

### 4. Base de Datos
- [ ] PostgreSQL corriendo
- [ ] Base de datos existe
- [ ] Credenciales correctas
- [ ] Datos inicializados

### 5. API
- [ ] Compilación exitosa
- [ ] Sin errores en logs
- [ ] Responde a health check

---

## Obtener Ayuda

Si el problema persiste:

1. **Revisar logs detalladamente**
```bash
   docker-compose logs
```

2. **Buscar en Issues de GitHub**
   - https://github.com/Zaragoza9512/go-api-chi/issues

3. **Crear Issue nuevo**
   - Incluir logs completos
   - Pasos para reproducir
   - Versiones de software

4. **Stack Overflow**
   - Tags: `go`, `docker`, `postgresql`, `terraform`

---

**Última actualización:** Marzo 2026

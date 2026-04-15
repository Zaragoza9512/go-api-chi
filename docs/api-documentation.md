# 📡 API Documentation

Documentación completa de todos los endpoints de la API.

**Base URL:** `http://localhost:8080`

**Autenticación:** JWT Bearer Token (excepto `/login`)

---

## Tabla de Contenidos

- [Autenticación](#autenticación)
- [Productos](#productos)
- [Códigos de Estado](#códigos-de-estado)
- [Errores](#errores)

---

## Autenticación

La API usa JWT (JSON Web Tokens) para autenticación.

### Flujo de Autenticación

---

### POST /login

Obtiene un token JWT para autenticación.

**Endpoint:** `POST /login`

**Headers:**
```http
Content-Type: application/json
```

**Body:**
```json
{
  "username": "admin",
  "password": "password"
}
```

**Respuesta Exitosa (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk4NTc2MDAsInVzZXJuYW1lIjoiYWRtaW4ifQ.xxx"
}
```

**Respuesta Error (401 Unauthorized):**
```json
{
  "error": "Credenciales inválidas"
}
```

**Ejemplo con cURL:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password"
  }'
```

**Notas:**
- El token expira en 24 horas
- Usuario y password actuales están hardcodeados
- En producción: implementar tabla de usuarios

---

## Productos

Todos los endpoints de productos requieren autenticación.

**Header requerido:**
```http
Authorization: Bearer {token}
```

---

### GET /productos

Obtiene la lista completa de productos.

**Endpoint:** `GET /productos`

**Headers:**
```http
Authorization: Bearer {token}
```

**Query Parameters:** Ninguno (por ahora)

**Respuesta Exitosa (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Laptop Dell XPS 15",
    "description": "Laptop de alto rendimiento con procesador Intel i7",
    "price": 1499.99,
    "stock": 10
  },
  {
    "id": 2,
    "name": "Mouse Logitech MX Master 3",
    "description": "Mouse inalámbrico ergonómico",
    "price": 99.99,
    "stock": 50
  }
]
```

**Respuesta Error (401 Unauthorized):**
```json
{
  "error": "Token inválido o expirado"
}
```

**Ejemplo con cURL:**
```bash
TOKEN="tu_token_aqui"

curl -X GET http://localhost:8080/productos \
  -H "Authorization: Bearer $TOKEN"
```

**Notas:**
- Retorna array vacío `[]` si no hay productos
- Futuro: implementar paginación

---

### GET /productos/{id}

Obtiene un producto específico por ID.

**Endpoint:** `GET /productos/{id}`

**Path Parameters:**
- `id` (integer, requerido) - ID del producto

**Headers:**
```http
Authorization: Bearer {token}
```

**Respuesta Exitosa (200 OK):**
```json
{
  "id": 1,
  "name": "Laptop Dell XPS 15",
  "description": "Laptop de alto rendimiento con procesador Intel i7",
  "price": 1499.99,
  "stock": 10
}
```

**Respuesta Error (404 Not Found):**
```json
{
  "error": "Producto no encontrado"
}
```

**Ejemplo con cURL:**
```bash
TOKEN="tu_token_aqui"

curl -X GET http://localhost:8080/productos/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Notas:**
- El ID debe ser un número entero positivo
- Retorna 404 si el producto no existe

---

### POST /productos

Crea un nuevo producto.

**Endpoint:** `POST /productos`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer {token}
```

**Body:**
```json
{
  "name": "Teclado Mecánico Keychron K8",
  "description": "Teclado mecánico inalámbrico con switches Gateron",
  "price": 89.99,
  "stock": 25
}
```

**Validaciones:**
- `name`: Requerido, string, max 255 caracteres
- `description`: Opcional, string
- `price`: Requerido, número positivo
- `stock`: Requerido, entero no negativo

**Respuesta Exitosa (201 Created):**
```json
{
  "id": 3,
  "name": "Teclado Mecánico Keychron K8",
  "description": "Teclado mecánico inalámbrico con switches Gateron",
  "price": 89.99,
  "stock": 25
}
```

**Respuesta Error (400 Bad Request):**
```json
{
  "error": "Datos inválidos: price debe ser mayor a 0"
}
```

**Ejemplo con cURL:**
```bash
TOKEN="tu_token_aqui"

curl -X POST http://localhost:8080/productos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Teclado Mecánico Keychron K8",
    "description": "Teclado mecánico inalámbrico con switches Gateron",
    "price": 89.99,
    "stock": 25
  }'
```

**Notas:**
- El `id` es auto-generado por la base de datos
- Los campos vacíos retornan error de validación

---

### PUT /productos/{id}

Actualiza un producto existente.

**Endpoint:** `PUT /productos/{id}`

**Path Parameters:**
- `id` (integer, requerido) - ID del producto a actualizar

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer {token}
```

**Body:**
```json
{
  "name": "Laptop Dell XPS 15 (2024)",
  "description": "Laptop actualizada con mejores specs",
  "price": 1599.99,
  "stock": 8
}
```

**Validaciones:**
- Mismas validaciones que POST
- Todos los campos son requeridos (no es PATCH)

**Respuesta Exitosa (200 OK):**
```json
{
  "id": 1,
  "name": "Laptop Dell XPS 15 (2024)",
  "description": "Laptop actualizada con mejores specs",
  "price": 1599.99,
  "stock": 8
}
```

**Respuesta Error (404 Not Found):**
```json
{
  "error": "Producto no encontrado"
}
```

**Ejemplo con cURL:**
```bash
TOKEN="tu_token_aqui"

curl -X PUT http://localhost:8080/productos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Laptop Dell XPS 15 (2024)",
    "description": "Laptop actualizada con mejores specs",
    "price": 1599.99,
    "stock": 8
  }'
```

**Notas:**
- Reemplaza completamente el producto
- Para actualización parcial, se necesitaría un endpoint PATCH

---

### DELETE /productos/{id}

Elimina un producto.

**Endpoint:** `DELETE /productos/{id}`

**Path Parameters:**
- `id` (integer, requerido) - ID del producto a eliminar

**Headers:**
```http
Authorization: Bearer {token}
```

**Respuesta Exitosa (204 No Content):**

**Respuesta Error (404 Not Found):**
```json
{
  "error": "Producto no encontrado"
}
```

**Ejemplo con cURL:**
```bash
TOKEN="tu_token_aqui"

curl -X DELETE http://localhost:8080/productos/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Notas:**
- La eliminación es permanente
- No hay confirmación adicional
- Futuro: implementar soft delete

---

## Códigos de Estado

| Código | Significado | Cuándo se usa |
|--------|-------------|---------------|
| 200 | OK | Request exitoso (GET, PUT) |
| 201 | Created | Recurso creado exitosamente (POST) |
| 204 | No Content | Eliminación exitosa (DELETE) |
| 400 | Bad Request | Datos inválidos en el body |
| 401 | Unauthorized | Token inválido, expirado o faltante |
| 404 | Not Found | Recurso no encontrado |
| 500 | Internal Server Error | Error del servidor |

---

## Errores

### Formato de Errores

Todos los errores siguen este formato:

```json
{
  "error": "Descripción del error"
}
```

### Errores Comunes

**Token Faltante:**
```json
{
  "error": "Token de autorización requerido"
}
```

**Token Inválido:**
```json
{
  "error": "Token inválido o expirado"
}
```

**Credenciales Incorrectas:**
```json
{
  "error": "Credenciales inválidas"
}
```

**Producto No Encontrado:**
```json
{
  "error": "Producto no encontrado"
}
```

**Datos Inválidos:**
```json
{
  "error": "Datos inválidos: [detalle específico]"
}
```

---

## Rate Limiting

**Actualmente:** No implementado

**Futuro:** 
- 100 requests por minuto por IP
- Header `X-RateLimit-Remaining` en respuestas

---

## Versionamiento

**Versión Actual:** v1 (implícita)

**Futuro:** Versionamiento en URL (`/api/v1/productos`)

---

## Ejemplos Completos

### Flujo Completo: Crear y Listar Productos

```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' \
  | jq -r '.token')

echo "Token obtenido: $TOKEN"

# 2. Crear producto
curl -X POST http://localhost:8080/productos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Monitor LG 27 pulgadas",
    "description": "Monitor 4K con HDR",
    "price": 399.99,
    "stock": 15
  }'

# 3. Listar todos los productos
curl -X GET http://localhost:8080/productos \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.'

# 4. Obtener producto específico (ID 1)
curl -X GET http://localhost:8080/productos/1 \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.'

# 5. Actualizar producto
curl -X PUT http://localhost:8080/productos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Monitor LG 27 pulgadas 4K",
    "description": "Monitor 4K actualizado con mejor refresh rate",
    "price": 449.99,
    "stock": 12
  }'

# 6. Eliminar producto
curl -X DELETE http://localhost:8080/productos/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## Testing con Scripts

### Script de prueba completo

Crea un archivo `test-api.sh`:

```bash
#!/bin/bash

API_URL="http://localhost:8080"

echo "=== Testing API ==="
echo ""

# Login
echo "1. Login..."
TOKEN=$(curl -s -X POST $API_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' \
  | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ Login failed"
  exit 1
fi
echo "✅ Login successful"
echo ""

# Get all products
echo "2. Getting all products..."
PRODUCTS=$(curl -s -X GET $API_URL/productos \
  -H "Authorization: Bearer $TOKEN")
echo "✅ Products retrieved"
echo "$PRODUCTS" | jq '.'
echo ""

# Create product
echo "3. Creating new product..."
NEW_PRODUCT=$(curl -s -X POST $API_URL/productos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Product",
    "description": "Testing API",
    "price": 99.99,
    "stock": 10
  }')
PRODUCT_ID=$(echo "$NEW_PRODUCT" | jq -r '.id')
echo "✅ Product created with ID: $PRODUCT_ID"
echo ""

# Get specific product
echo "4. Getting product $PRODUCT_ID..."
curl -s -X GET "$API_URL/productos/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.'
echo "✅ Product retrieved"
echo ""

# Update product
echo "5. Updating product $PRODUCT_ID..."
curl -s -X PUT "$API_URL/productos/$PRODUCT_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Product Updated",
    "description": "Updated description",
    "price": 149.99,
    "stock": 5
  }' | jq '.'
echo "✅ Product updated"
echo ""

# Delete product
echo "6. Deleting product $PRODUCT_ID..."
curl -s -X DELETE "$API_URL/productos/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN"
echo "✅ Product deleted"
echo ""

echo "=== All tests passed! ==="
```

**Uso:**
```bash
chmod +x test-api.sh
./test-api.sh
```

---

## Notas para Desarrolladores

### Headers Importantes

```http
Content-Type: application/json        # Para POST/PUT
Authorization: Bearer {token}         # Para autenticación
```

### Manejo de Tokens

- Guardar token en variable de entorno o storage
- Renovar token antes de expiración
- Manejar errores 401 (token expirado)

### Best Practices

1. **Siempre validar respuestas**
```javascript
   if (response.status === 401) {
     // Token expirado, renovar
   }
```

2. **Usar HTTPS en producción**
   - Nunca enviar tokens por HTTP

3. **Manejar errores gracefully**
```javascript
   try {
     const data = await fetch(url);
   } catch (error) {
     console.error('API Error:', error);
   }
```

---

## Próximas Funcionalidades

- [ ] Paginación en `/productos`
- [ ] Búsqueda y filtros
- [ ] Categorías de productos
- [ ] Tabla de usuarios
- [ ] Refresh tokens
- [ ] Rate limiting
- [ ] Webhooks

---

**Última actualización:** Marzo 2026  
**Versión API:** 1.0

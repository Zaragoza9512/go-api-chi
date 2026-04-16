# 🤝 Contributing to go-api-chi

¡Gracias por tu interés en contribuir! Este documento te guiará en el proceso.

---

## 📋 Tabla de Contenidos

- [Código de Conducta](#código-de-conducta)
- [Cómo Contribuir](#cómo-contribuir)
- [Configuración del Entorno](#configuración-del-entorno)
- [Workflow de Desarrollo](#workflow-de-desarrollo)
- [Estándares de Código](#estándares-de-código)
- [Commits](#commits)
- [Pull Requests](#pull-requests)
- [Reportar Bugs](#reportar-bugs)
- [Proponer Features](#proponer-features)

---

## 📜 Código de Conducta

Este proyecto sigue principios de respeto y colaboración:

- ✅ Ser respetuoso con otros contributors
- ✅ Aceptar críticas constructivas
- ✅ Enfocarse en lo mejor para el proyecto
- ❌ No usar lenguaje ofensivo
- ❌ No hacer ataques personales

---

## 🚀 Cómo Contribuir

Hay muchas formas de contribuir:

### 1. Reportar Bugs
- Usa el template de Issues
- Incluye pasos para reproducir
- Adjunta logs si es posible

### 2. Sugerir Features
- Abre un Issue con el tag `enhancement`
- Explica el caso de uso
- Describe la solución propuesta

### 3. Escribir Código
- Arreglar bugs
- Implementar features
- Mejorar tests
- Optimizar performance

### 4. Mejorar Documentación
- Corregir typos
- Agregar ejemplos
- Traducir docs
- Escribir tutoriales

---

## 🛠️ Configuración del Entorno

### Prerequisitos

```bash
# Versiones mínimas
Go 1.23+
Docker 20.10+
Docker Compose v2.0+
Git
```

### Setup

1. **Fork el repositorio**
   
   Click en "Fork" en GitHub

2. **Clonar tu fork**
```bash
   git clone git@github.com:TU_USUARIO/go-api-chi.git
   cd go-api-chi
```

3. **Agregar remote upstream**
```bash
   git remote add upstream git@github.com:Zaragoza9512/go-api-chi.git
```

4. **Instalar dependencias**
```bash
   go mod download
```

5. **Configurar variables de entorno**
```bash
   cp .env.example .env
   nano .env
```

6. **Levantar el entorno**
```bash
   docker-compose up --build
```

7. **Verificar que funciona**
```bash
   ./test-api.sh
```

---

## 🔄 Workflow de Desarrollo

### 1. Sincronizar con upstream

Antes de empezar, sincroniza tu fork:

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

### 2. Crear una rama

```bash
# Formato: tipo/descripcion-corta
git checkout -b feature/add-user-table
git checkout -b fix/jwt-expiration-bug
git checkout -b docs/api-examples
```

**Tipos de ramas:**
- `feature/` - Nueva funcionalidad
- `fix/` - Corrección de bug
- `docs/` - Documentación
- `refactor/` - Refactorización
- `test/` - Tests

### 3. Hacer cambios

```bash
# Editar archivos
nano handlers.go

# Probar localmente
go run .
go test ./...

# Verificar con linter (opcional)
golangci-lint run
```

### 4. Commit

```bash
git add .
git commit -m "feat: add user authentication table"
```

Ver [sección de Commits](#commits) para el formato.

### 5. Push a tu fork

```bash
git push origin feature/add-user-table
```

### 6. Abrir Pull Request

1. Ve a tu fork en GitHub
2. Click en "Compare & pull request"
3. Llena el template
4. Espera review

---

## 📏 Estándares de Código

### Go Style

Seguimos [Effective Go](https://go.dev/doc/effective_go):

**✅ Buenas prácticas:**

```go
// Nombres descriptivos
func GetProductByID(id int) (*Product, error) {
    // ...
}

// Manejar errores
if err != nil {
    return nil, fmt.Errorf("failed to get product: %w", err)
}

// Comentarios en exported functions
// GetProductByID retrieves a product by its ID from the database.
// Returns an error if the product is not found.
func GetProductByID(id int) (*Product, error) {
    // ...
}
```

**❌ Evitar:**

```go
// Nombres vagos
func Get(i int) (*Product, error) {

// Ignorar errores
result, _ := GetProduct(id)

// Sin comentarios en exported functions
func GetProductByID(id int) (*Product, error) {
```

### Formato

```bash
# Formatear código
go fmt ./...

# Verificar imports
goimports -w .

# Linter (si está instalado)
golangci-lint run
```

### Tests

- Escribir tests para nueva funcionalidad
- Mantener coverage > 70%
- Tests deben ser rápidos

```go
// Example test
func TestGetProductByID(t *testing.T) {
    // Arrange
    product := &Product{ID: 1, Name: "Test"}
    
    // Act
    result, err := GetProductByID(1)
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.Name != product.Name {
        t.Errorf("expected %s, got %s", product.Name, result.Name)
    }
}
```

---

## 💬 Commits

Usamos [Conventional Commits](https://www.conventionalcommits.org/).

### Formato

### Tipos

| Tipo | Descripción | Ejemplo |
|------|-------------|---------|
| `feat` | Nueva funcionalidad | `feat: add user login endpoint` |
| `fix` | Corrección de bug | `fix: resolve JWT token expiration bug` |
| `docs` | Documentación | `docs: update API endpoints in README` |
| `style` | Formato (sin cambio lógico) | `style: format code with gofmt` |
| `refactor` | Refactorización | `refactor: extract validation logic` |
| `test` | Tests | `test: add unit tests for handlers` |
| `chore` | Tareas de mantenimiento | `chore: update dependencies` |

### Ejemplos

**Feature:**

### Reglas

- ✅ Presente imperativo ("add" no "added")
- ✅ Primera línea < 72 caracteres
- ✅ Cuerpo opcional con detalles
- ✅ Footer con referencias a issues
- ❌ No terminar con punto

---

## 🔀 Pull Requests

### Template

Al abrir un PR, llena este template:

```markdown
## Descripción
Breve descripción de los cambios.

## Tipo de cambio
- [ ] Bug fix
- [ ] Nueva funcionalidad
- [ ] Breaking change
- [ ] Documentación

## ¿Cómo probarlo?
Pasos para verificar los cambios.

## Checklist
- [ ] Tests pasan localmente
- [ ] Código formateado (go fmt)
- [ ] Documentación actualizada
- [ ] Sin breaking changes (o documentados)
```

### Proceso de Review

1. **CI debe pasar** ✅
   - Tests
   - Build
   - Linting

2. **Code review**
   - Al menos 1 aprobación requerida
   - Resolver comentarios

3. **Merge**
   - Squash and merge (preferido)
   - O merge commit

### Durante el Review

**Si te piden cambios:**

```bash
# Hacer cambios
nano handlers.go

# Commit
git add .
git commit -m "fix: address review comments"

# Push
git push origin feature/add-user-table
```

El PR se actualiza automáticamente.

---

## 🐛 Reportar Bugs

### Antes de reportar

1. Buscar en Issues existentes
2. Verificar con última versión
3. Reproducir el bug

### Template de Bug Report

```markdown
## Descripción del Bug
Descripción clara del problema.

## Pasos para Reproducir
1. Ejecutar '...'
2. Hacer click en '...'
3. Ver error

## Comportamiento Esperado
Qué debería pasar.

## Comportamiento Actual
Qué pasa realmente.

## Logs
```
Pegar logs relevantes aquí

---

## 💡 Proponer Features

### Template de Feature Request

```markdown
## Problema
¿Qué problema resuelve esta feature?

## Solución Propuesta
Cómo funcionaría.

## Alternativas Consideradas
Otras opciones evaluadas.

## Contexto Adicional
Screenshots, ejemplos, etc.
```

### Discusión

- Las features se discuten en Issues primero
- Se busca consenso antes de implementar
- Prioridad según roadmap del proyecto

---

## ✅ Checklist Pre-Commit

Antes de hacer commit, verifica:

- [ ] Código compila sin errores
- [ ] Tests pasan
- [ ] Código formateado (`go fmt`)
- [ ] Sin warnings del linter
- [ ] Documentación actualizada
- [ ] Commit message sigue convención
- [ ] Sin credenciales en código
- [ ] `.env` no se incluye

---

## 🎯 Áreas que Necesitan Ayuda

Áreas donde contribuciones son especialmente bienvenidas:

### Backend
- [ ] Implementar tabla de usuarios real
- [ ] Agregar paginación
- [ ] Implementar categorías
- [ ] Búsqueda y filtros

### Testing
- [ ] Aumentar coverage
- [ ] Tests de integración
- [ ] Tests E2E

### Documentación
- [ ] Más ejemplos de uso
- [ ] Tutoriales
- [ ] Traducción a inglés

### DevOps
- [ ] Kubernetes manifests
- [ ] Helm charts
- [ ] Scripts de deployment

---

## 📚 Recursos

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Chi Router Docs](https://go-chi.io/)

---

## 🙏 Agradecimientos

Gracias a todos los contributors:

<!-- Lista generada automáticamente por GitHub -->

---

## ❓ Preguntas

Si tienes dudas:

1. Buscar en [Issues](https://github.com/Zaragoza9512/go-api-chi/issues)
2. Abrir nueva [Discussion](https://github.com/Zaragoza9512/go-api-chi/discussions)
3. Email: zaragoza95.luis@gmail.com

---

¡Happy coding! 🚀

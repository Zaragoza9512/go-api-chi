# ====================================================================
# Multi-stage build para optimizar tamaño de imagen
# ====================================================================

# Etapa 1: Builder (compilar el binario Go)
FROM golang:1.23-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias primero (mejor uso de caché)
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
# CGO_ENABLED=0: binario estático (no necesita librerías externas)
# -ldflags="-w -s": reduce tamaño del binario
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api .

# ====================================================================
# Etapa 2: Runtime (imagen final mínima)
# ====================================================================
FROM alpine:latest

# Instalar certificados SSL (necesarios para HTTPS)
RUN apk --no-cache add ca-certificates

# Crear usuario no-root por seguridad
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copiar el binario desde la etapa builder
COPY --from=builder /app/api .

# Cambiar ownership al usuario no-root
RUN chown -R appuser:appuser /app

# Cambiar a usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./api"]
# Fase de compilación
FROM golang:1.23 as builder

# Instalar dependencias necesarias
RUN apt-get update && apt-get install -y gcc libc6-dev make sqlite3 libsqlite3-dev

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos esenciales primero para aprovechar caché de Docker
COPY go.mod go.sum Makefile ./

# Descargar dependencias
RUN go mod tidy

# Copiar el código fuente completo
COPY . .

# Compilar el binario
RUN make build

# Fase de ejecución
FROM debian:stable-slim

# Instalar SQLite
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0 && rm -rf /var/lib/apt/lists/*

# Establecer directorio de trabajo
WORKDIR /app

# Copiar el binario compilado desde la fase de construcción
COPY --from=builder /app/main /usr/local/bin/main

# Dar permisos de ejecución al binario
RUN chmod +x /usr/local/bin/main

# Comando de inicio
ENTRYPOINT ["/usr/local/bin/main"]

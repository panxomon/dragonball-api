BUILDPATH=$(CURDIR)
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
API_NAME=dragonball-test
LOWER_API_NAME=$(shell echo $(API_NAME) | tr A-Z a-z)

dir:
	@echo "full path: $(BUILDPATH)"

# Compilar el binario
build: tidy
	@echo "Compiling binary..."
	@CGO_ENABLED=1 go build -o /app/main ./cmd/main.go
	@echo "Binary generated at /app/main"

# Descargar y ordenar dependencias
.PHONY: tidy
tidy:
	@go mod tidy

# Instalar herramientas y dependencias adicionales
.PHONY: dep
dep: tidy
	@echo "Installing dependencies..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/go-critic/go-critic/cmd/gocritic@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go mod download
	@echo "OK"

# Ejecutar pruebas unitarias
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -short $(PKG_LIST)
	@echo "OK"

# Ejecutar la aplicación localmente
.PHONY: run
run:
	@go run cmd/main.go

# Formatear el código
.PHONY: fmt
fmt: tidy
	@echo "Formatting Go code..."
	@go fmt ./... $(PKG_LIST)
	@echo "OK"

# Linter y verificación de código
.PHONY: lint
lint: tidy
	@echo "Checking code style..."
	@staticcheck ./... $(PKG_LIST)
	@go vet ./... $(PKG_LIST)
	@gocritic check ./... $(PKG_LIST)
	@echo "OK"

# Pruebas de carrera de datos
.PHONY: race
race: tidy
	@go test ./... -race -short $(PKG_LIST)

# Generar mocks
.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	@rm -rf ./tests/mocks/
	@mockery --all --dir internal --output ./tests/mocks --case underscore
	@echo "OK"

# Generar documentación con Swag
.PHONY: docs
docs:
	@echo "Generating docs..."
	@swag fmt
	@swag init --g ./cmd/main.go --markdownFiles ./docs/api --codeExampleFiles ./docs/examples --parseInternal
	@echo "OK"

# Restaurar contexto de Docker
.PHONY: docker-context
docker-context:
	@echo "Setting docker context to default again"
	@docker context use default
	@echo "OK"

# Construir y ejecutar Docker Compose
.PHONY: build_and_up
build_and_up:
	@echo "Building and starting the application..."
	@docker-compose build && docker-compose up

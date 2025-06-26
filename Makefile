.PHONY: build-wasm clean up up-fast down logs help rebuild-ext-proc

# Default target
help:
	@echo "Available targets:"
	@echo "  build-wasm      - Build the WASM filter"
	@echo "  up              - Build WASM and start environment (forces rebuild)"
	@echo "  up-fast         - Start environment without forcing rebuild"
	@echo "  rebuild-ext-proc - Force rebuild just the ext-proc service"
	@echo "  down            - Stop the environment"
	@echo "  logs            - Show logs from all services"
	@echo "  clean           - Clean build artifacts"
	@echo "  help            - Show this help message"

# Build the WASM filter
build-wasm:
	@echo "Building WASM filter..."
	cargo build --target wasm32-wasip1 --release
	@echo "WASM filter built successfully!"

# Build WASM and start the environment (with forced rebuild)
up: build-wasm
	docker-compose up --build

# Start environment without rebuild
up-fast: build-wasm
	docker-compose up

# Stop the environment
down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Force rebuild just the ext-proc service
rebuild-ext-proc:
	docker-compose build --no-cache ext-proc
	docker-compose up -d ext-proc

# Clean build artifacts
clean:
	cargo clean
	docker-compose down --volumes --remove-orphans 
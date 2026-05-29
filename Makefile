.PHONY: dev up down migrate seed test build clean


build:
	docker-compose -f docker-compose.yml build

up:
	docker-compose -f docker-compose.yml up -d

# Stop all services
down:
	docker-compose -f docker-compose.yml down


docker-restart:
	docker-compose -f docker-compose.yml up -d --build

docker-logs:
	docker-compose -f docker-compose.yml logs -f

# Run database migrations
migrate:
	cd backend && go run ./cmd/migrate

# Generate image variants for all existing media
generate-variants:
	cd backend && go run ./cmd/server generate-variants

# Run backend dev server
dev-backend:
	cd backend && go run ./cmd/server

# Run frontend dev server
dev-frontend:
	cd frontend && npm run dev

# Run both dev servers
dev:
	@echo "Start MySQL first: make up"
	@echo "Then run in separate terminals:"
	@echo "  make dev-backend"
	@echo "  make dev-frontend"

# Run tests
test:
	cd backend && go test ./... -v

# Run tests with coverage
test-cover:
	cd backend && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out


# Install dependencies
install:
	cd backend && go mod tidy
	cd frontend && npm install

# Clean
clean:
	docker-compose down -v
	rm -rf frontend/.nuxt frontend/.output frontend/node_modules

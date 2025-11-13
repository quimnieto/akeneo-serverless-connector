.PHONY: help install start stop backend frontend clean build

help:
	@echo "Available commands:"
	@echo "  make install  - Install all dependencies"
	@echo "  make start    - Start both backend and frontend"
	@echo "  make backend  - Start backend server"
	@echo "  make frontend - Start frontend dev server"
	@echo "  make build    - Build both backend and frontend"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make stop     - Stop all running processes"

install:
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "✓ All dependencies installed"

backend:
	@echo "Starting backend server..."
	@echo "Using config: config/settings.$${ENV:-dev}.json"
	cd backend && go run cmd/server/main.go

frontend:
	@echo "Starting frontend dev server..."
	cd frontend && npm start

start:
	@echo "Starting backend and frontend..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:3000"
	@echo ""
	@echo "Run 'make backend' and 'make frontend' in separate terminals"

build:
	@echo "Building backend..."
	cd backend && go build -o ../bin/server cmd/server/main.go
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "✓ Build complete"
	@echo "  Backend binary: bin/server"
	@echo "  Frontend build: frontend/build/"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf frontend/build/
	rm -rf frontend/node_modules/
	cd backend && go clean
	@echo "✓ Clean complete"

stop:
	@echo "Stopping processes..."
	@pkill -f "go run cmd/server/main.go" || true
	@pkill -f "react-scripts start" || true
	@echo "✓ Processes stopped"

# Akeneo Event Platform Configuration

Simple app to manage Akeneo Event Platform configurations (subscriber and subscriptions).

## Prerequisites

- Go 1.21+
- Node.js 18+
- Akeneo PIM instance with Event Platform enabled

## Quick Start

1. Configure your Akeneo instance in `backend/config/settings.dev.json`:
```json
{
  "akeneo": {
    "base_url": "https://your-akeneo-instance.com",
    "client_id": "your_client_id",
    "client_secret": "your_client_secret",
    "username": "your_username",
    "password": "your_password"
  },
  "server": {
    "port": "8080",
    "gin_mode": "debug"
  },
  "cors": {
    "allowed_origins": "http://localhost:3000"
  }
}
```

2. Install dependencies:
```bash
make install
```

3. Start the application (in separate terminals):
```bash
# Terminal 1 - Backend (defaults to dev environment)
make backend

# Or specify environment
ENV=integration make backend
ENV=prod make backend

# Terminal 2 - Frontend
make frontend
```

Backend runs on `http://localhost:8080`, Frontend on `http://localhost:3000`

## Environment Configuration

The app uses environment-specific JSON config files located in `backend/config/`:
- `backend/config/settings.dev.json` - Development (default)
- `backend/config/settings.integration.json` - Integration/Staging
- `backend/config/settings.prod.json` - Production

Copy the example and create your environment configs:
```bash
cp backend/config/settings.example.json backend/config/settings.dev.json
# Edit backend/config/settings.dev.json with your credentials
```

Set the `ENV` or `ENVIRONMENT` environment variable to switch between configurations:
```bash
export ENV=integration  # or dev, prod
# or
export ENVIRONMENT=integration
```

## Makefile Commands

- `make install` - Install all dependencies
- `make backend` - Start backend server
- `make frontend` - Start frontend dev server
- `make build` - Build production binaries
- `make clean` - Clean build artifacts
- `make stop` - Stop all running processes
- `make help` - Show available commands

## Usage

### Subscriber
- Configure the webhook URL where Akeneo will send events
- Toggle active/inactive status

### Subscriptions
- Create subscriptions for specific connection codes
- Select which event types to subscribe to (fetched from Akeneo)
- Activate/deactivate or delete subscriptions

## API Endpoints

- `GET /api/subscriber` - Get current subscriber
- `POST /api/subscriber` - Create subscriber
- `PATCH /api/subscriber` - Update subscriber
- `GET /api/subscriptions` - List all subscriptions
- `POST /api/subscriptions` - Create subscription
- `PATCH /api/subscriptions/:code` - Update subscription
- `DELETE /api/subscriptions/:code` - Delete subscription
- `GET /api/event-types` - Get available event types

## Production Build

Frontend:
```bash
cd akeneo/frontend
npm run build
```

Backend:
```bash
cd akeneo/backend
go build -o server cmd/server/main.go
./server
```

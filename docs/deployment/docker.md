# Docker Deployment

This guide covers deploying Laju Go using Docker with multi-stage builds for optimized container images.

## Prerequisites

- Docker installed on your system
- Docker Compose (optional, for multi-container setups)

## Quick Start

### Build and Run

```bash
# Build the image
docker build -t laju-go .

# Run the container
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  -v $(pwd)/storage:/root/storage \
  --name laju-go \
  laju-go
```

### Access Application

Visit `http://localhost:8080`

## Dockerfile

### Multi-Stage Build

```dockerfile
# ===========================================
# Stage 1: Build Frontend
# ===========================================
FROM node:20-alpine AS frontend

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy frontend source
COPY . .

# Build frontend
RUN npm run build

# ===========================================
# Stage 2: Build Go Binary
# ===========================================
FROM golang:1.22-alpine AS backend

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend assets
COPY --from=frontend /app/dist ./dist

# Build Go binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o laju-go .

# ===========================================
# Stage 3: Production Image
# ===========================================
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

# Copy binary from backend stage
COPY --from=backend /app/laju-go .

# Copy migrations
COPY --from=backend /app/migrations ./migrations

# Copy templates
COPY --from=backend /app/templates ./templates

# Create data and storage directories
RUN mkdir -p /app/data /app/storage/avatars && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["./laju-go"]
```

## Docker Compose

### Basic Setup

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    container_name: laju-go
    ports:
      - "8080:8080"
    volumes:
      - ./data:/root/data
      - ./storage:/root/storage
    environment:
      - APP_ENV=production
      - APP_PORT=8080
      - DB_PATH=/root/data/app.db
      - SESSION_SECRET=${SESSION_SECRET}
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s
```

### With Nginx Reverse Proxy

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    container_name: laju-go
    expose:
      - "8080"
    volumes:
      - ./data:/root/data
      - ./storage:/root/storage
    environment:
      - APP_ENV=production
      - APP_PORT=8080
      - DB_PATH=/root/data/app.db
      - SESSION_SECRET=${SESSION_SECRET}
    restart: unless-stopped
    networks:
      - laju-network

  nginx:
    image: nginx:alpine
    container_name: laju-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - laju-network

networks:
  laju-network:
    driver: bridge
```

### Nginx Configuration

```nginx
# nginx/nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream laju_app {
        server app:8080;
    }

    server {
        listen 80;
        server_name _;

        location / {
            proxy_pass http://laju_app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }
    }
}
```

## Environment Variables

### Using .env File

```bash
# .env
SESSION_SECRET=your-32-character-secret-key
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
```

### Docker Compose with .env

```yaml
services:
  app:
    environment:
      - SESSION_SECRET=${SESSION_SECRET}
      - GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASS=${SMTP_PASS}
```

## Data Persistence

### Volume Mounts

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  -v $(pwd)/storage:/root/storage \
  laju-go
```

### Named Volumes

```yaml
version: '3.8'

services:
  app:
    build: .
    volumes:
      - laju-data:/root/data
      - laju-storage:/root/storage

volumes:
  laju-data:
  laju-storage:
```

## Building

### Build Image

```bash
docker build -t laju-go:latest .
```

### Build with Tag

```bash
docker build -t laju-go:1.0.0 -t laju-go:latest .
```

### Build Specific Stage

```bash
# Build frontend only
docker build --target frontend -t laju-go-frontend .

# Build backend only
docker build --target backend -t laju-go-backend .
```

## Running

### Run Container

```bash
docker run -d \
  --name laju-go \
  -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  -v $(pwd)/storage:/root/storage \
  -e SESSION_SECRET=your-secret-key \
  laju-go
```

### Run with Docker Compose

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Rebuild and restart
docker-compose up -d --build
```

## Container Management

### View Logs

```bash
# Real-time logs
docker logs -f laju-go

# Last 100 lines
docker logs --tail 100 laju-go

# With timestamps
docker logs -f --timestamps laju-go
```

### Execute Commands

```bash
# Open shell in container
docker exec -it laju-go /bin/sh

# Run database migrations
docker exec laju-go goose -dir migrations sqlite3 /root/data/app.db up

# Check database
docker exec laju-go sqlite3 /root/data/app.db ".tables"
```

### Container Health

```bash
# Check health status
docker inspect --format='{{.State.Health.Status}}' laju-go

# View health check logs
docker inspect --format='{{json .State.Health}}' laju-go | jq
```

### Stop and Remove

```bash
# Stop container
docker stop laju-go

# Remove container
docker rm laju-go

# Remove image
docker rmi laju-go:latest
```

## Production Deployment

### Push to Registry

```bash
# Tag image
docker tag laju-go:latest registry.example.com/laju-go:1.0.0

# Login to registry
docker login registry.example.com

# Push image
docker push registry.example.com/laju-go:1.0.0
```

### Deploy to Server

```bash
# On production server
docker pull registry.example.com/laju-go:1.0.0

docker run -d \
  --name laju-go \
  -p 8080:8080 \
  -v /var/lib/laju/data:/root/data \
  -v /var/lib/laju/storage:/root/storage \
  --env-file /etc/laju-go/.env \
  --restart unless-stopped \
  registry.example.com/laju-go:1.0.0
```

### Docker Swarm

```yaml
# docker-stack.yml
version: '3.8'

services:
  app:
    image: registry.example.com/laju-go:1.0.0
    ports:
      - "8080:8080"
    volumes:
      - laju-data:/root/data
      - laju-storage:/root/storage
    environment:
      - SESSION_SECRET=${SESSION_SECRET}
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure

volumes:
  laju-data:
  laju-storage:
```

```bash
# Deploy stack
docker stack deploy -c docker-stack.yml laju-go

# View services
docker stack services laju-go

# Scale service
docker service scale laju-go_app=5
```

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/docker.yml
name: Docker Build and Push

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to Registry
        uses: docker/login-action@v2
        with:
          registry: registry.example.com
          username: ${{ secrets.REGISTRY_USERNAME }}
          password: ${{ secrets.REGISTRY_PASSWORD }}

      - name: Build and Push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: |
            registry.example.com/laju-go:${{ github.ref_name }}
            registry.example.com/laju-go:latest
          cache-from: type=registry,ref=registry.example.com/laju-go:buildcache
          cache-to: type=registry,ref=registry.example.com/laju-go:buildcache,mode=max
```

## Optimization Tips

### Reduce Image Size

```dockerfile
# Use multi-stage builds (already implemented)
# Use Alpine base images
# Copy only necessary files
# Remove unnecessary dependencies

# Before: ~1.2GB
FROM golang:1.22
# ...

# After: ~50MB
FROM alpine:latest
# ...
```

### Layer Caching

```dockerfile
# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code last
COPY . .
```

### Use .dockerignore

```
# .dockerignore
node_modules
npm-debug.log
.git
.gitignore
*.md
.env
data/
storage/
tmp/
```

## Troubleshooting

### Container Won't Start

**Check logs**:

```bash
docker logs laju-go
```

**Common issues**:
- Missing environment variables
- Database path not writable
- Port already in use

### Database Not Persisting

**Solution**: Ensure volume is mounted correctly

```bash
# Check volume mount
docker inspect laju-go | grep -A 10 Mounts

# Verify data exists
ls -la ./data/
```

### Permission Denied

**Solution**: Run as root or fix permissions

```dockerfile
# Option 1: Run as root (not recommended for production)
USER root

# Option 2: Fix permissions in Dockerfile
RUN chown -R appuser:appgroup /app
```

### Health Check Failing

**Solution**: Increase start period or fix endpoint

```yaml
healthcheck:
  test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
  interval: 30s
  timeout: 3s
  retries: 3
  start_period: 30s  # Increase from 5s
```

## Best Practices

### 1. Use Multi-Stage Builds

```dockerfile
# ✅ Good: Multi-stage build
FROM node:20-alpine AS frontend
# ...
FROM golang:1.22-alpine AS backend
# ...
FROM alpine:latest
COPY --from=backend /app/laju-go .

# ❌ Bad: Single stage
FROM golang:1.22
# Everything in one layer
```

### 2. Don't Run as Root

```dockerfile
# ✅ Good: Non-root user
RUN adduser -D appuser
USER appuser

# ❌ Bad: Root user
# No USER instruction
```

### 3. Use Specific Tags

```dockerfile
# ✅ Good: Specific versions
FROM node:20-alpine
FROM golang:1.22-alpine

# ❌ Bad: Latest tag
FROM node:latest
FROM golang:latest
```

### 4. Minimize Layers

```dockerfile
# ✅ Good: Combined commands
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    && rm -rf /var/cache/apk/*

# ❌ Bad: Multiple RUN commands
RUN apk add ca-certificates
RUN apk add sqlite
```

### 5. Use Health Checks

```dockerfile
# ✅ Good: Health check defined
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --spider http://localhost:8080/health || exit 1

# ❌ Bad: No health check
```

## Next Steps

- [Production Deployment](production.md) - Bare-metal deployment
- [Optimization Guide](optimization.md) - Performance tuning
- [CI/CD Guide](../guide/cicd.md) - Continuous integration and deployment

# Docker Compose 生产部署 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为项目添加 Docker 生产部署能力：3 个 Dockerfile + nginx 配置 + docker-compose.prod.yml + 环境变量模板 + Makefile 扩展 + SSR 配置修改。

**Architecture:** 4 容器编排 — nginx 反向代理（80/443）→ frontend (Node SSR :3000) + backend (Go :8080) → mysql (内部 :3306)。多阶段构建减小镜像体积。

**Tech Stack:** Docker, docker-compose, Go 1.25 alpine, Node 22 alpine, Nginx alpine, MySQL 8.0

---

## File Map

| 文件 | 操作 | 职责 |
|------|------|------|
| `backend/Dockerfile` | Create | Go 多阶段构建（编译→运行） |
| `backend/.dockerignore` | Create | 排除不必要文件加速构建 |
| `frontend/Dockerfile` | Create | Nuxt 3 Node SSR 多阶段构建 |
| `frontend/.dockerignore` | Create | 排除 node_modules/.nuxt/.output |
| `nginx/Dockerfile` | Create | 基于 nginx:alpine 注入自定义配置 |
| `nginx/nginx.conf` | Create | 反向代理 /api/v1→backend, /*→frontend |
| `docker-compose.prod.yml` | Create | 4 service 编排 + 自定义网络 + 数据卷 |
| `.env.production.example` | Create | 环境变量模板（提交到 git） |
| `frontend/nuxt.config.ts` | Modify | `ssr: false` → `ssr: true` |
| `Makefile` | Modify | 追加 docker-build/up/down/logs/restart |
| `.gitignore` | Modify | 添加 `.env.production` |

---

### Task 1: Backend Dockerfile + .dockerignore

**Files:**
- Create: `backend/Dockerfile`
- Create: `backend/.dockerignore`

- [ ] **Step 1: Create `backend/.dockerignore`**

```
.git
.env
*.out
coverage.out
__debug_bin
```

- [ ] **Step 2: Create `backend/Dockerfile`**

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/server .
COPY database/migrations ./database/migrations
EXPOSE 8080
CMD ["./server"]
```

- [ ] **Step 3: Verify build**

```bash
cd backend && docker build -t mygo-backend:test .
```

Expected: Image builds successfully, tagged as `mygo-backend:test`.

- [ ] **Step 4: Verify image size**

```bash
docker images mygo-backend:test --format "{{.Size}}"
```

Expected: Under 25MB.

- [ ] **Step 5: Commit**

```bash
git add backend/Dockerfile backend/.dockerignore
git commit -m "feat: add backend Dockerfile with multi-stage build"
```

---

### Task 2: Frontend Dockerfile + .dockerignore

**Files:**
- Create: `frontend/Dockerfile`
- Create: `frontend/.dockerignore`

- [ ] **Step 1: Create `frontend/.dockerignore`**

```
node_modules
.nuxt
.output
dist
.git
.gitignore
```

- [ ] **Step 2: Create `frontend/Dockerfile`**

```dockerfile
FROM node:22-alpine AS builder
WORKDIR /build
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npx nuxi build

FROM node:22-alpine
WORKDIR /app
COPY --from=builder /build/.output ./.output
COPY --from=builder /build/node_modules ./node_modules
EXPOSE 3000
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000
CMD ["node", ".output/server/index.mjs"]
```

- [ ] **Step 3: Verify build**

```bash
cd frontend && docker build -t mygo-frontend:test .
```

Expected: Image builds successfully (takes a few minutes for npm ci + build).

- [ ] **Step 4: Commit**

```bash
git add frontend/Dockerfile frontend/.dockerignore
git commit -m "feat: add frontend Dockerfile with multi-stage Node SSR build"
```

---

### Task 3: Nginx Dockerfile + nginx.conf

**Files:**
- Create: `nginx/Dockerfile`
- Create: `nginx/nginx.conf`

- [ ] **Step 1: Create `nginx/` directory**

```bash
mkdir -p nginx
```

- [ ] **Step 2: Create `nginx/nginx.conf`**

```nginx
worker_processes auto;

events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    keepalive_timeout  65;

    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml text/javascript image/svg+xml;
    gzip_min_length 256;

    server {
        listen 80;

        # Uploads served by backend
        location /uploads/ {
            proxy_pass http://backend:8080;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # API proxy to backend
        location /api/ {
            proxy_pass http://backend:8080;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # All other requests to frontend
        location / {
            proxy_pass http://frontend:3000;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

- [ ] **Step 3: Create `nginx/Dockerfile`**

```dockerfile
FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
```

- [ ] **Step 4: Verify build**

```bash
cd nginx && docker build -t mygo-nginx:test .
```

Expected: Image builds successfully.

- [ ] **Step 5: Commit**

```bash
git add nginx/Dockerfile nginx/nginx.conf
git commit -m "feat: add nginx reverse proxy config and Dockerfile"
```

---

### Task 4: docker-compose.prod.yml

**File:**
- Create: `docker-compose.prod.yml`

- [ ] **Step 1: Create `docker-compose.prod.yml`**

```yaml
services:
  mysql:
    image: mysql:8.0
    container_name: mygo-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: ${MYSQL_DATABASE}
      MYSQL_USER: ${MYSQL_USER}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
    volumes:
      - mysql_data:/var/lib/mysql
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
    networks:
      - mygo-net
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build: ./backend
    container_name: mygo-backend
    restart: unless-stopped
    environment:
      DB_HOST: mysql
      DB_PORT: "3306"
      DB_USER: ${MYSQL_USER}
      DB_PASSWORD: ${MYSQL_PASSWORD}
      DB_NAME: ${MYSQL_DATABASE}
      JWT_SECRET: ${JWT_SECRET}
      JWT_ACCESS_EXPIRY: ${JWT_ACCESS_EXPIRY:-24h}
      JWT_REFRESH_EXPIRY: ${JWT_REFRESH_EXPIRY:-168h}
      SERVER_PORT: "8080"
      CORS_ORIGIN: ${CORS_ORIGIN:-https://your-domain.com}
      ENV: production
    networks:
      - mygo-net
    depends_on:
      mysql:
        condition: service_healthy

  frontend:
    build: ./frontend
    container_name: mygo-frontend
    restart: unless-stopped
    environment:
      NITRO_HOST: "0.0.0.0"
      NITRO_PORT: "3000"
    networks:
      - mygo-net
    depends_on:
      - backend

  nginx:
    build: ./nginx
    container_name: mygo-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    networks:
      - mygo-net
    depends_on:
      - frontend
      - backend

networks:
  mygo-net:
    driver: bridge

volumes:
  mysql_data:
```

- [ ] **Step 2: Validate compose file syntax**

```bash
docker compose -f docker-compose.prod.yml config
```

Expected: No syntax errors, outputs the merged compose config.

- [ ] **Step 3: Commit**

```bash
git add docker-compose.prod.yml
git commit -m "feat: add production docker-compose with 4 services"
```

---

### Task 5: Environment variable template + .gitignore

**Files:**
- Create: `.env.production.example`
- Modify: `.gitignore`

- [ ] **Step 1: Check existing `.gitignore` for `.env` patterns**

Read `.gitignore` to see current state, then add `.env.production`.

- [ ] **Step 2: Create `.env.production.example`**

```
# MySQL
MYSQL_ROOT_PASSWORD=<generate-a-strong-password>
MYSQL_DATABASE=immigration
MYSQL_USER=immigration
MYSQL_PASSWORD=<generate-a-strong-password>

# JWT
JWT_SECRET=<at-least-32-random-characters>
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Server
CORS_ORIGIN=https://your-domain.com
```

- [ ] **Step 3: Add `.env.production` to `.gitignore`**

Read the current `.gitignore` first, then append:

```
# Production environment
.env.production
```

- [ ] **Step 4: Commit**

```bash
git add .env.production.example .gitignore
git commit -m "feat: add production env template and gitignore entry"
```

---

### Task 6: Enable SSR in nuxt.config.ts

**File:**
- Modify: `frontend/nuxt.config.ts:3`

- [ ] **Step 1: Change `ssr: false` to `ssr: true`**

In `frontend/nuxt.config.ts`, line 3, change:
```
  ssr: false,
```
to:
```
  ssr: true,
```

- [ ] **Step 2: Verify the change**

```bash
grep "ssr:" frontend/nuxt.config.ts
```

Expected output: `  ssr: true,`

- [ ] **Step 3: Commit**

```bash
git add frontend/nuxt.config.ts
git commit -m "feat: enable SSR for production Docker deployment"
```

---

### Task 7: Makefile Docker commands

**File:**
- Modify: `Makefile`

- [ ] **Step 1: Append docker commands to Makefile**

Append the following to the end of `Makefile`:

```makefile
# Docker 生产部署
docker-build:
	docker compose -f docker-compose.prod.yml build

docker-up:
	docker compose -f docker-compose.prod.yml up -d

docker-down:
	docker compose -f docker-compose.prod.yml down

docker-logs:
	docker compose -f docker-compose.prod.yml logs -f

docker-restart:
	docker compose -f docker-compose.prod.yml up -d --build
```

- [ ] **Step 2: Verify Makefile syntax**

```bash
make -n docker-build
```

Expected: Prints `docker compose -f docker-compose.prod.yml build` without errors.

- [ ] **Step 3: Commit**

```bash
git add Makefile
git commit -m "feat: add docker production commands to Makefile"
```

---

## Completion Checklist

- [ ] All 4 Docker images build successfully
- [ ] `docker compose -f docker-compose.prod.yml config` passes validation
- [ ] `make docker-build` runs without errors
- [ ] Backend health check returns 200: `curl http://localhost/api/v1/health` (after `docker-up`)
- [ ] Frontend pages accessible via nginx at `http://localhost`
- [ ] API calls proxied correctly through nginx
- [ ] MySQL data persists across container restarts

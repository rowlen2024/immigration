# Docker Compose 生产部署设计

## 1. 整体架构

```
┌────────────────────────────────────────────────┐
│                   Nginx (:80/:443)              │
│              SSL 终止 · 反向代理 · gzip         │
└─────┬──────────────────────────┬───────────────┘
      │ /api/v1/*                │ /*
      ▼                          ▼
┌──────────────┐         ┌──────────────────┐
│   backend    │         │    frontend       │
│   Go :8080   │         │  Node SSR :3000   │
│  (内部网络)   │         │   (内部网络)       │
└──────┬───────┘         └──────────────────┘
       │
       ▼
┌──────────────┐
│    mysql     │
│   :3306      │
│  (内部网络)   │
└──────────────┘
```

| 要点 | 说明 |
|------|------|
| 网络 | 自定义 bridge 网络 `mygo-net`，只有 nginx 对外暴露端口 |
| 后端 | backend 和 mysql 通过容器名互相发现，`DB_HOST=mysql` |
| 前端 | Node SSR 监听 3000，仅 nginx 可达，CORS 由 nginx 统一处理 |
| MySQL | 数据卷 `mysql_data` 持久化，只在 docker 内部网络暴露 3306 |
| 启动顺序 | mysql → backend（等 healthcheck）→ frontend → nginx |

关键决策：
- 后端内部不再需要 CORS 中间件（所有请求来自同域 nginx），但保留以兼容开发环境
- 前端当前 `nuxt.config.ts` 中 `ssr: false`（SPA 模式）。部署需改为 `ssr: true` 以启用服务端渲染
- 迁移脚本在 backend 容器启动时自动执行（main.go 已有 RunMigrations + AutoMigrate）

---

## 2. Dockerfile 设计

### 部署文件结构

```
immigration/
├── backend/Dockerfile
├── frontend/Dockerfile
├── nginx/
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml          # 保留，开发用（仅 MySQL）
├── docker-compose.prod.yml     # 新增，生产部署
├── .env.production             # 新增，生产环境变量（不提交到 git）
└── .env.production.example     # 新增，环境变量模板（提交到 git）
```

### Backend Dockerfile（多阶段构建）

```dockerfile
# 阶段1: 编译
FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# 阶段2: 运行
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/server .
COPY database/migrations ./database/migrations
EXPOSE 8080
CMD ["./server"]
```

- 编译启用 CGO_ENABLED=0 静态链接
- 运行镜像含二进制 + migrations 目录 + 时区
- 最终镜像约 18-22 MB
- 配置全部通过环境变量注入，不依赖 `.env` 文件

### Frontend Dockerfile（多阶段构建）

```dockerfile
# 阶段1: 构建
FROM node:22-alpine AS builder
WORKDIR /build
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npx nuxi build

# 阶段2: 运行
FROM node:22-alpine
WORKDIR /app
COPY --from=builder /build/.output ./.output
COPY --from=builder /build/node_modules ./node_modules
EXPOSE 3000
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000
CMD ["node", ".output/server/index.mjs"]
```

- 构建阶段需要完整 node_modules（含 devDependencies for sass/typescript）
- 运行阶段包含 .output/ 和 node_modules（仅生产依赖）
- 最终镜像约 200-300 MB

### Nginx（简单 Dockerfile + 静态配置）

```dockerfile
FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
```

nginx.conf 核心：
- 监听 80，`/api/v1/` 代理到 `backend:8080`，其余到 `frontend:3000`
- gzip 开启，静态资源缓存头
- SSL 证书通过 volume 挂载（443 部分可选）

---

## 3. docker-compose.prod.yml 设计

### MySQL

```yaml
mysql:
  image: mysql:8.0
  restart: unless-stopped
  environment:
    MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
    MYSQL_DATABASE: ${MYSQL_DATABASE}
    MYSQL_USER: ${MYSQL_USER}
    MYSQL_PASSWORD: ${MYSQL_PASSWORD}
  volumes:
    - mysql_data:/var/lib/mysql
  command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
  healthcheck:
    test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
    interval: 10s
    timeout: 5s
    retries: 5
```

### Backend

```yaml
backend:
  build: ./backend
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
  depends_on:
    mysql:
      condition: service_healthy
```

### Frontend

```yaml
frontend:
  build: ./frontend
  restart: unless-stopped
  environment:
    NITRO_HOST: "0.0.0.0"
    NITRO_PORT: "3000"
  depends_on:
    - backend
```

说明：useApi 已使用相对路径 `baseURL: '/api/v1'`，同域经 nginx 代理到后端，无需额外配置。

### Nginx

```yaml
nginx:
  build: ./nginx
  restart: unless-stopped
  ports:
    - "80:80"
    - "443:443"
  depends_on:
    - frontend
    - backend
```

---

## 4. 环境变量

`.env.production`（不提交到 git）：

```
MYSQL_ROOT_PASSWORD=<强密码>
MYSQL_DATABASE=immigration
MYSQL_USER=immigration
MYSQL_PASSWORD=<强密码>
JWT_SECRET=<至少32位随机字符串>
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h
CORS_ORIGIN=https://your-domain.com
```

---

## 5. 部署流程

### 首次部署

```bash
# 1. 配置环境变量
cp .env.production.example .env.production
# 编辑 .env.production，填入真实密码和域名

# 2. 构建镜像
docker compose -f docker-compose.prod.yml build

# 3. 启动全部服务
docker compose -f docker-compose.prod.yml up -d

# 4. 查看启动日志
docker compose -f docker-compose.prod.yml logs -f

# 5. 确认后端健康
curl http://localhost:80/api/v1/health
```

### 日常运维

```bash
# 代码更新后重新构建并部署
docker compose -f docker-compose.prod.yml up -d --build

# 仅重启某个服务
docker compose -f docker-compose.prod.yml restart backend

# 查看某个服务日志
docker compose -f docker-compose.prod.yml logs -f --tail=100 backend

# 停止所有服务
docker compose -f docker-compose.prod.yml down

# 完全清理含数据卷（慎用）
docker compose -f docker-compose.prod.yml down -v
```

### 数据库备份与恢复

```bash
# 备份
docker compose -f docker-compose.prod.yml exec mysql \
  mysqldump -u immigration -p${MYSQL_PASSWORD} immigration > backup.sql

# 恢复
docker compose -f docker-compose.prod.yml exec -T mysql \
  mysql -u immigration -p${MYSQL_PASSWORD} immigration < backup.sql
```

### SSL 证书

1. **自签名（测试）**：openssl 生成证书，nginx 配置指向 `/etc/nginx/ssl/` volume
2. **Let's Encrypt（生产）**：搭配 certbot 容器或 acme.sh，自动续期

---

## 6. Makefile 扩展

在现有 Makefile 追加：

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

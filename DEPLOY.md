# MyGo 移民 - Docker 部署指南

## 环境要求

- Docker 20.10+
- Docker Compose v2

## 首次部署

### 1. 配置环境变量

```bash
cp .env.production.example .env.production
```

编辑 `.env.production`，替换以下占位符：

| 变量 | 说明 |
|------|------|
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码 |
| `MYSQL_PASSWORD` | MySQL 应用用户密码 |
| `JWT_SECRET` | 至少 32 位随机字符串 |
| `CORS_ORIGIN` | 生产域名，如 `https://example.com` |

### 2. 构建并启动

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

首次启动会：构建前后端镜像 → 初始化 MySQL → 自动执行数据库迁移 → 启动服务。

### 3. 验证

```bash
curl http://localhost/api/v1/health
```

返回 `{"code":200,"message":"ok"}` 即部署成功。

## 架构

```
nginx (:80) → /api/*  → backend (Go :8080) → mysql (:3306)
            → /*      → frontend (Nuxt :3000)
```

所有服务在内部网络 `mygo-net` 中通信，仅 nginx 对外暴露端口。

## 常用命令

```bash
# 查看日志
docker compose -f docker-compose.prod.yml logs -f

# 仅重启后端
docker compose -f docker-compose.prod.yml restart backend

# 代码更新后重新构建
docker compose -f docker-compose.prod.yml up -d --build

# 停止所有服务
docker compose -f docker-compose.prod.yml down

# 完全清理（含数据，慎用）
docker compose -f docker-compose.prod.yml down -v
```

## 数据库备份

```bash
# 备份
docker compose -f docker-compose.prod.yml exec mysql \
  mysqldump -u immigration -p${MYSQL_PASSWORD} immigration > backup.sql

# 恢复
docker compose -f docker-compose.prod.yml exec -T mysql \
  mysql -u immigration -p${MYSQL_PASSWORD} immigration < backup.sql
```

## SSL 配置

### 测试环境（自签名）

```bash
mkdir -p ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ssl/privkey.pem -out ssl/fullchain.pem
```

在 `nginx/nginx.conf` 的 server 块中添加：

```nginx
listen 443 ssl;
ssl_certificate     /etc/nginx/ssl/fullchain.pem;
ssl_certificate_key /etc/nginx/ssl/privkey.pem;
```

### 生产环境（Let's Encrypt）

推荐使用 [acme.sh](https://github.com/acmesh-official/acme.sh) 或 certbot 获取免费证书，证书挂载到 nginx 容器的 `/etc/nginx/ssl/` 路径。

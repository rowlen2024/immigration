# MyGo 移民 - 部署指南

## 架构

```
用户 :80 → nginx → frontend:3000 (Nuxt 3 SSR)
                   → /api/* → backend:8080 (Go/Gin)
                   → /uploads/* → backend:8080
                              backend:8080 → mysql:3306
```

所有服务在内部网络 `mygo-net` 中通信，仅 nginx 对外暴露端口。

---

## 本地构建运行

```bash
# 配置环境变量
cp .env.production .env
# 编辑 .env，修改 JWT_SECRET

# 构建并启动
docker compose -f docker-compose.prod.yml up -d --build
```

验证：

```bash
curl http://localhost/api/v1/health
# → {"code":200,"message":"ok"}
```

---

## 云服务器部署（推荐：容器镜像仓库）

### 1. 创建镜像仓库

登录 [阿里云容器镜像服务](https://cr.console.aliyun.com/)，创建命名空间和 3 个仓库：

```
registry.cn-hangzhou.aliyuncs.com/<命名空间>/mygo-backend
registry.cn-hangzhou.aliyuncs.com/<命名空间>/mygo-frontend
registry.cn-hangzhou.aliyuncs.com/<命名空间>/mygo-nginx
```

> 也可用 Docker Hub（`docker.io/<username>/`）或腾讯云 TCR（`ccr.ccs.tencentyun.com/`），流程相同。

### 2. 修改部署文件中的仓库地址

编辑 `docker-compose.deploy.yml`，将 `your-namespace` 替换为你的命名空间名。

### 3. 构建并推送镜像

```bash
# 登录
docker login --username=<阿里云账号> registry.cn-hangzhou.aliyuncs.com

# 设置命名空间
NS=registry.cn-hangzhou.aliyuncs.com/<你的命名空间>

# 构建
docker build -t $NS/mygo-backend:latest  ./backend
docker build -t $NS/mygo-frontend:latest ./frontend
docker build -t $NS/mygo-nginx:latest    ./nginx

# 推送
docker push $NS/mygo-backend:latest
docker push $NS/mygo-frontend:latest
docker push $NS/mygo-nginx:latest
```

### 4. 上传部署文件到服务器

```bash
# 本地执行
scp docker-compose.deploy.yml root@<服务器IP>:/opt/mygo/
scp .env.production root@<服务器IP>:/opt/mygo/.env
```

### 5. 服务器上启动

```bash
# SSH 到服务器
ssh root@<服务器IP>
cd /opt/mygo

# 编辑 .env，修改 JWT_SECRET 等密钥
vim .env

# 登录镜像仓库
docker login --username=<阿里云账号> registry.cn-hangzhou.aliyuncs.com

# 启动
docker compose -f docker-compose.deploy.yml up -d

# 检查运行状态
docker ps -a --filter "name=mygo"

# 查看日志
docker compose -f docker-compose.deploy.yml logs -f
```

---

## 备选方案

### 方案B：服务器上直接构建

适合服务器性能好、不想折腾镜像仓库的场景。

```bash
# 服务器上
git clone <仓库地址> /opt/mygo
cd /opt/mygo
cp .env.production .env
# 编辑 .env，修改所有密钥和密码

docker compose -f docker-compose.prod.yml up -d --build
```

> Go 后端和 Nuxt 前端构建比较吃内存，建议服务器 4GB+ 内存。

### 方案C：离线导出传输

适合无法登录镜像仓库的离线环境。

```bash
# === 本地 ===
docker pull mysql:8.0

# 构建所有镜像
docker compose -f docker-compose.prod.yml build

# 导出为 tar
docker save -o mygo-backend.tar  mygo-backend:latest
docker save -o mygo-frontend.tar mygo-frontend:latest
docker save -o mygo-nginx.tar    mygo-nginx:latest
docker save -o mygo-mysql.tar    mysql:8.0

# 传输到服务器
scp *.tar docker-compose.deploy.yml .env.production root@<服务器IP>:/opt/mygo/

# === 服务器 ===
cd /opt/mygo
docker load -i mygo-mysql.tar
docker load -i mygo-backend.tar
docker load -i mygo-frontend.tar
docker load -i mygo-nginx.tar

cp .env.production .env
docker compose -f docker-compose.deploy.yml up -d
```

---

## 生产环境检查清单

- [ ] `.env` 中 `JWT_SECRET` 已改为 32 位以上随机字符串
- [ ] `.env` 中 `MYSQL_ROOT_PASSWORD` 和 `MYSQL_PASSWORD` 已改为强密码
- [ ] `CORS_ORIGIN` 已改为实际域名（如 `https://example.com`）
- [ ] 配置了 HTTPS 证书
- [ ] 设置了数据库定期备份
- [ ] 防火墙只开放 80/443，**不要暴露 3306/8080/3000**

---

## 配置 HTTPS

使用 acme.sh 申请免费 Let's Encrypt 证书：

```bash
# 安装
curl https://get.acme.sh | sh

# 申请证书
acme.sh --issue -d your-domain.com --nginx

# 安装证书到指定目录
mkdir -p /opt/mygo/nginx/ssl
acme.sh --install-cert -d your-domain.com \
  --key-file       /opt/mygo/nginx/ssl/key.pem \
  --fullchain-file /opt/mygo/nginx/ssl/cert.pem
```

然后在 nginx 配置中添加 443 监听并挂载证书目录。

---

## 数据库备份与恢复

```bash
# 备份（在服务器上执行）
docker exec mygo-mysql mysqldump \
  -u root -p${MYSQL_ROOT_PASSWORD} immigration \
  > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复
docker exec -i mygo-mysql mysql \
  -u root -p${MYSQL_ROOT_PASSWORD} immigration < backup.sql
```

建议设置 cron 定时备份：

```bash
# 每天凌晨 3 点备份
0 3 * * * cd /opt/mygo && docker exec mygo-mysql mysqldump -u root -p<密码> immigration > /opt/backups/immigration_$(date +\%Y\%m\%d).sql
```

---

## 常用运维命令

```bash
# 查看日志
docker compose -f docker-compose.deploy.yml logs -f

# 查看各容器状态
docker ps -a --filter "name=mygo"

# 重启单个服务
docker compose -f docker-compose.deploy.yml restart backend

# 更新镜像（推送新版本后）
docker compose -f docker-compose.deploy.yml pull
docker compose -f docker-compose.deploy.yml up -d

# 停止所有服务
docker compose -f docker-compose.deploy.yml down

# 完全清理（含数据，慎用）
docker compose -f docker-compose.deploy.yml down -v
```

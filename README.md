# Study-UPC

> 学院学习资料托管平台 - 面向学院的学习资料共享与管理平台

## 📘 项目简介

Study-UPC 是面向学院的学习资料管理平台，提供资料上传下载、审核与检索服务。项目采用前后端分离架构，支持多角色权限管理（学生、学委、管理员），确保内容质量与合规性。

### ✅ 核心功能

- **用户认证**：注册登录、基于角色的权限控制（RBAC）
- **资料管理**：资料上传、下载、分类、标签
- **审核流程**：学委审核、内容审核、违规举报处理
- **检索与推荐**：全文搜索、多条件筛选、热门资料推荐
- **通知系统**：审核结果通知、系统公告、站内消息
- **数据统计**：上传下载统计、用户活跃度分析

## 🛠 技术栈

| 层次 | 技术选型 |
|------|----------|
| **前端** | Vue 3 + TypeScript + Vite + Element Plus + Pinia |
| **后端** | Go 1.21 + Gin + GORM + JWT |
| **数据库** | PostgreSQL 15 + Redis 7 |
| **存储** | MinIO（开发）/ 阿里云 OSS（生产） |
| **部署** | Nginx + systemd（推荐） |

## 📁 项目结构

```
study-upc/
├── backend/                # Go 后端服务
│   ├── cmd/server/          # 程序入口
│   ├── configs/             # 配置文件
│   ├── internal/            # 私有代码
│   ├── logs/                # 日志目录
│   ├── migrations/          # 数据库迁移
│   ├── go.mod
│   └── go.sum
├── frontend/               # Vue 3 前端应用
│   ├── src/
│   ├── dist/                # 构建产物
│   ├── package.json
│   └── vite.config.ts
├── scripts/                # 运维/辅助脚本
│   ├── docker-up.bat
│   ├── docker-down.bat
│   ├── setup-postgres.sh
│   ├── setup-redis.sh
│   └── setup-nginx.sh
├── CLAUDE.md
└── README.md
```

## 🚀 本地开发

### 1) 准备依赖

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- 对象存储：MinIO（本地）/ OSS（生产环境）

### 2) 配置后端

编辑 `backend/configs/config.dev.yaml`，确保数据库、Redis、MinIO 配置可用。

### 3) 启动后端

```bash
cd backend
go run cmd/server/main.go
```

### 4) 启动前端

```bash
cd frontend
npm ci
npm run dev
```

默认 Vite 端口为 `5173`，后端默认端口为 `8080`。

## 🧩 生产部署（不使用 Docker）

### 1) 配置后端（阿里云 OSS）

编辑 `backend/configs/config.prod.yaml`：

- 配置 PostgreSQL 与 Redis 连接信息
- 修改 `jwt.secret`
- 配置 `oss` 为阿里云 OSS（access key、endpoint、region）

### 2) 构建并运行后端

```bash
cd backend
go build -o bin/server ./cmd/server
export APP_ENV=production
export CONFIG_PATH=configs/config.prod.yaml
./bin/server
```

健康检查：`http://127.0.0.1:8080/health`

### 3) 构建前端

```bash
cd frontend
npm ci
npm run build
```

### 4) 配置 Nginx 反向代理

示例（`/etc/nginx/conf.d/study-upc.conf`）：

```nginx
server {
    listen 80;
    server_name your.domain.com;

    root /path/to/study-upc/frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/v1/ {
        proxy_pass http://127.0.0.1:8080/api/v1/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

```bash
sudo nginx -t && sudo systemctl reload nginx
```

### 5) 使用部署脚本（可选）

```bash
sudo bash scripts/setup-postgres.sh study_upc your_db_password study_upc
sudo bash scripts/setup-redis.sh your_redis_password
sudo bash scripts/setup-nginx.sh example.com /opt/study-upc/frontend/dist 127.0.0.1 8080
```

## 🧰 脚本说明

- `scripts/setup-postgres.sh`：初始化 PostgreSQL 用户与数据库
- `scripts/setup-redis.sh`：设置 Redis 密码并限制本机访问
- `scripts/setup-nginx.sh`：生成 Nginx 站点配置并热加载
- `scripts/docker-up.bat`/`docker-down.bat`：Windows 下启动/停止基础服务（仅在有 Docker 时使用）

## 🤝 贡献指南

欢迎提交 Issue 或 Pull Request。

1. Fork 仓库
2. 创建分支（`feature/xxx`）
3. 提交修改（`feat: xxx`）
4. 推送分支并发起 PR

## 📄 许可证

本项目采用 [MIT](LICENSE) 许可证。

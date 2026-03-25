# 远楒文档

远楒文档是一个支持实时协作的文档系统，包含 Go 后端、Vue 前端和基于 Yjs 的协作栈。

## 项目特色
- **实时协作**：基于 Yjs + WebSocket + Redis/MQ
- **Connect RPC/gRPC**：强类型接口与统一协议
- **模块化服务**：后端 + 协作网关 + 异步 worker
- **配置驱动**：Consul KV 动态配置
- **Docker 优先**：Nginx 反向代理 + 一键启停

## 项目架构与链路
**核心服务**
- `frontend`（Vue）通过 Nginx 对外提供
- `backend`（Connect RPC/gRPC）
- `collab-core`（Node/Yjs）+ `collab-go` 网关（WebSocket）
- `snapshot-worker`（异步消费者/后台任务）
- 基础设施：Postgres / Redis / RabbitMQ / Consul

**典型请求链路**
1. 浏览器 -> Nginx（`/` 为 UI，`/grpc/` 为 gRPC‑web）
2. Nginx -> Backend（Connect gRPC‑web）
3. Backend -> Postgres / Redis / Consul / MQ
4. 协作链路：浏览器 -> Nginx `/ws/doc/` -> `collab-go` -> `collab-core`（Yjs）-> Redis/MQ/Backend
5. 异步消费：`snapshot-worker` 从 MQ 消费任务（如文档快照）

## 项目结构
- `backend`: Go 服务、数据迁移、Connect/gRPC 服务端、Worker。
- `frontend`: Vue 3 前端（Vite + Element Plus）。
- `proto`: Protobuf 定义与生成。
- `collab`: Node.js 的 Yjs 协作核心。
- `collab-go`: Go 协作网关。
- `deploy`: Docker Compose、Nginx、脚本与基础设施配置。
- `docs`: 文档。

## 启动方式（Docker Compose）
开发环境：
```bash
bash deploy/script/start.sh dev up
```

生产环境：
```bash
bash deploy/script/start.sh release up
```

## Prepare 脚本
用于从 example 生成配置文件：
```bash
bash deploy/script/prepare.sh
```
会生成/覆盖：
- `backend/config.toml`
- `frontend/nginx.conf`
- `deploy/nginx/nginx.conf`
- `deploy/nginx/conf.d/default.conf`
- `deploy/redis/redis.conf`
- `deploy/rabbitmq/rabbitmq.conf`

## 启动参数与 Token
这些值在 compose 和 `backend/config.toml` 中使用，**生产环境必须替换**。

### Docker Compose 环境变量
- `CONSUL_ROOT_TOKEN`
  Consul ACL 的 root token，默认值必须改。
- `BACKEND_INTERNAL_RPC_KEY`
  内部服务间通信的共享密钥。
- `VITE_GRPC_WEB_ENDPOINT`（前端）
  gRPC‑web 访问前缀（默认 `/grpc`）。

### 后端配置文件（`backend/config.toml`）
至少需要更新：
- **数据库**：`database.user` / `database.password` / `database.name`
- **Redis**：`redis.password`
- **Consul**：`consul.token` / `consul.address` / `consul.prefix`
- **消息队列**：`mq_config.type` / `mq_config.rabbitmq.url`
- **JWT**：`backend.jwt.secret`
- **安全配置**：`backend.security`

### 其他服务凭据
- **Postgres**：`POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB`
- **Redis**：`REDIS_PASSWORD`
- **RabbitMQ**：`RABBITMQ_DEFAULT_USER` / `RABBITMQ_DEFAULT_PASS`

## 生产环境注意事项
1. 替换所有默认 token/密码。
2. 将真实的 `backend/config.toml` 挂载到 backend 与 snapshot‑worker。
3. 若对外提供服务，配置 Nginx TLS（`deploy/nginx/conf.d/default.conf`）。
4. 为 Postgres / Redis / Consul / RabbitMQ 配置持久化卷。

## 服务与端口
默认开发环境（`deploy/docker-compose.dev.yml`）：
- **Nginx**：`http://localhost:8080`
- **Backend gRPC**：`:9090`（内部网络）
- **Backend gRPC‑web**：`http://localhost:8080/grpc/`
- **协作 WS**：`http://localhost:8080/ws/doc/`
- **Consul UI**：`http://localhost:8500`
- **Postgres**：`localhost:5432`
- **Redis**：`localhost:6379`
- **RabbitMQ**：`localhost:5672`（AMQP），管理页面可通过 Nginx `/rabbitmq/` 访问

其他服务（开发环境无端口暴露）：
- **snapshot-worker**：消费 MQ 任务（如文档快照）

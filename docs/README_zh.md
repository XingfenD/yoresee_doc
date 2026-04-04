# 远楒文档

远楒文档是一个支持实时协作的文档平台，技术栈为 Go 后端、Vue 3 前端和基于 Yjs 的协作服务。

## 运行拓扑

核心业务服务：
- `frontend`：Vue 3（Vite）
- `backend`：Go Connect RPC/gRPC 服务
- `collab-core`：Node.js Yjs 协作核心
- `collab`：`collab-go` WebSocket 网关
- `snapshot-worker`：文档快照同步 worker
- `notification-worker`：通知事件消费 worker
- `search-sync-worker`：搜索索引同步 worker

基础设施服务：
- `postgres`
- `redis`
- `rabbitmq`
- `consul`
- `minio`
- `elasticsearch`
- `nginx`（统一入口）

## 请求与事件链路

主请求链路：
1. 浏览器 -> Nginx
2. Nginx `/` -> Frontend
3. Nginx `/grpc/` -> Backend gRPC-web
4. Backend 按需访问 Postgres/Redis/MinIO/Consul/Elasticsearch/MQ

实时协作链路：
1. 浏览器 -> Nginx `/ws/doc/{docId}`
2. Nginx -> `collab-go`
3. `collab-go` 校验 JWT 后代理到 `collab-core`
4. `collab-core` 在内存和 Redis 中维护 Yjs 文档状态

异步事件链路：
1. `collab-core` 在 Redis 标记脏文档并发布脏文档事件
2. `snapshot-worker` 消费 RabbitMQ 的 `collab.dirty_docs`（组消费模式），并周期扫描脏集合兜底
3. `snapshot-worker` 调用 `collab-core` 的 `/internal/yjs/doc-snapshot/{docId}`，落库快照与正文；正文变化时发布搜索同步事件
4. `search-sync-worker` 消费 `search.sync.document` 并写入 Elasticsearch（ES 启用时）
5. `notification-worker` 消费 `notification.create` 并写入通知记录

## 健康检查与优雅停机

Backend 探针：
- `/health`
- `/readyz`
- `/livez`

协作服务探针：
- `collab-core`：`/health`、`/readyz`、`/livez`
- `collab-go`：`/health`、`/readyz`、`/livez`

停机行为：
- Backend 与协作服务支持 draining + 优雅停机。
- draining 阶段 readiness 会返回 `not_ready`。

## 部署模式

`deploy/docker-compose.dev.yml`：
- 源码挂载开发模式
- 暴露 backend 调试端口（`2345`）
- Postgres/Redis/RabbitMQ/Elasticsearch 端口对宿主机开放
- `start.sh dev ...` 会自动执行 `deploy/script/gen_proto.sh`

`deploy/docker-compose.yml`：
- 偏生产形态镜像
- 对外端口更少
- backend 与 worker 在镜像构建阶段完成编译

## 开发前置依赖

`start.sh dev ...` 会在 docker compose 启动前，先在宿主机执行 `deploy/script/gen_proto.sh`。

宿主机需要：
- `docker` + `docker compose`
- `go` + `protoc`（生成 Go 桩代码）
- `frontend/node_modules` 中的插件（`protoc-gen-es`、`protoc-gen-connect-es`）
- `PATH` 中可用的 `grpc_tools_node_protoc_plugin`（生成 `collab` 侧 Node gRPC 代码）

## 配置工作流

统一配置源：
- `deploy/.env`
- 模板：`deploy/.env.example`

交互式初始化/更新：
```bash
bash deploy/script/configure.sh
```

`configure.sh` 的作用：
- 交互式收集关键端口、密钥、环境变量
- 写入/更新 `deploy/.env`
- 由 `VITE_API_BASE_URL` 自动推导 `MINIO_BROWSER_REDIRECT_URL`
- 自动调用 `prepare.sh`

渲染配置文件：
```bash
bash deploy/script/prepare.sh
```

生成文件：
- `backend/config.toml`
- `frontend/nginx.conf`
- `deploy/nginx/nginx.conf`
- `deploy/nginx/conf.d/default.conf`
- `deploy/redis/redis.conf`
- `deploy/rabbitmq/rabbitmq.conf`

模板文件：
- `backend/config.toml.tmpl`
- `frontend/nginx.conf.tmpl`
- `deploy/nginx/nginx.conf.tmpl`
- `deploy/nginx/conf.d/default.conf.tmpl`
- `deploy/redis/redis.conf.tmpl`
- `deploy/rabbitmq/rabbitmq.conf.tmpl`

## 快速启动

推荐（交互式）：
```bash
bash deploy/script/configure.sh
bash deploy/script/start.sh dev up
```

非交互：
```bash
cp deploy/.env.example deploy/.env
bash deploy/script/prepare.sh
bash deploy/script/start.sh dev up
```

生产模式：
```bash
bash deploy/script/start.sh release up
```

backend 镜像启动行为：
- backend 容器会先执行 `migrate`、`db_init`、`es_init`，再启动 `cmd/main`。

其他常用动作：
```bash
bash deploy/script/start.sh dev rebuild
bash deploy/script/start.sh dev restart
bash deploy/script/start.sh dev clear
```

## 对外入口

默认开发环境：
- 主站：`http://localhost:8080`
- gRPC-web：`http://localhost:8080/grpc/`
- 协作 WS：`ws://localhost:8080/ws/doc/{docId}?token={jwt}`
- MinIO 控制台（经 Nginx）：`http://localhost:8080/minio/`
- 对象访问路径（经 Nginx）：`http://localhost:8080/storage/...`
- RabbitMQ 管理页（经 Nginx）：`http://localhost:8080/rabbitmq/`

开发环境基础设施直连端口：
- Postgres：`localhost:5432`
- Redis：`localhost:6379`
- RabbitMQ AMQP：`localhost:5672`
- Consul：`localhost:8500`
- Elasticsearch：`localhost:9200`

## MQ 说明

- 当前 worker 消费链路按 RabbitMQ 后端运行。
- MQ 接口层支持 Redis 与 RabbitMQ 两种实现。
- Redis Pub/Sub 不提供真正的组消费语义。
- 脏文档链路建议将 `DIRTY_DOC_MQ` 设为 `rabbitmq` 或 `both`，以保证 `snapshot-worker` 能从 RabbitMQ 收到事件。

## Protobuf 生成

手动生成：
```bash
bash deploy/script/gen_proto.sh
```

生成目标：
- `backend/pkg/gen`
- `collab-go/pkg/gen`
- `frontend/src/gen`
- `collab/src/gen`

## 仓库结构

- `backend`：API 服务、worker、仓储层、存储集成
- `frontend`：Vue 前端
- `collab`：Node.js 协作核心
- `collab-go`：Go WebSocket 协作网关
- `proto`：Protobuf 协议
- `deploy`：Compose、脚本、基础设施模板
- `docs`：项目文档

# Order Express（Docker + Postgres + Redis）

本项目使用 `docker-compose.yaml` 启动 **4 个常驻容器**：

- `frontend`：Nginx 托管前端静态资源，并反向代理 `/api/*` 到后端
- `backend`：Go API 服务（启动时可选自动 migrate/seed）
- `db`：PostgreSQL（数据持久化到 named volume）
- `redis`：Redis（用于缓存与限流；数据持久化到 named volume）

> 另外还提供一个按需运行的 `migrate` 工具服务（不常驻）：用于上线前执行数据库迁移。

## 1. 快速开始（本地）

### 1) 准备环境变量

复制一份根目录配置：

```bash
cp .env.example .env
```

默认本地会：

- `DB_AUTO_MIGRATE=true`：后端启动时自动执行数据库迁移
- `DB_AUTO_SEED=true`：后端启动时自动插入演示数据（幂等，不会重复插入）

### 2) 启动

```bash
docker compose up -d --build
```

访问：

- 前端：`http://localhost:8080`
- 后端（可选直连）：`http://localhost:3000/healthz`

如果端口冲突：

- 修改 `.env` 中的 `FRONTEND_PORT` / `BACKEND_PORT`

默认后台账号：

- 用户名：`admin`
- 密码：`admin123`

> 仅用于本地演示（`DB_AUTO_SEED=true`）。非本地环境请禁用自动 seed，并设置强密码与 `JWT_SECRET`。

## 2. 订单类型约定

`POST /api/orders` 的 `type` 支持三种：

- `dine_in`（堂食）：必须传 `tableId`
  - `tableId` 可传 **桌号**（如 `T05`）或 **桌台 id**（如 `table-05`），后端会自动解析并占用/释放桌台
- `takeaway`（外送）：不需要桌号，必须传 `deliveryAddress` + `contactPhone`
- `pickup`（自取）：不需要桌号/地址，必须传 `contactPhone`

## 3. 数据库迁移策略

### 本地（开发）

推荐保持 `.env` 中：

- `DB_AUTO_MIGRATE=true`
- `DB_AUTO_SEED=true`

这样 `docker compose up -d` 后即可直接使用。

### 上线（发布/迭代）

建议流程：**先迁移，再更新镜像**。

1) 先拉取/构建新版本镜像（示例）

```bash
docker compose build backend frontend
```

2) 运行一次性迁移任务（只跑迁移，不启动服务）

```bash
docker compose --profile tools run --rm migrate
```

3) 再更新前后端

```bash
docker compose up -d backend frontend
```

> 说明：数据库与 Redis 使用 named volume 持久化，更新前后端镜像不会影响数据。

## 4. “不使用外键”的实现方式

数据库层面：

- migrations 里**不创建外键约束**（DDL 不出现 `REFERENCES/FOREIGN KEY`）

逻辑层面：

- 在创建/更新时显式校验引用是否存在（例如创建菜品时先校验分类是否存在）
- 删除分类时若仍有菜品引用该分类会返回冲突错误（禁止删除）
- 关键写路径使用事务保证一致性（下单占用桌台、支付写入与订单更新、退款等）

## 5. 缓存与限流（Redis）

- 缓存：对 `categories/dishes/store` 等 GET 接口做短 TTL 缓存，并在写操作后失效
- 限流：对登录、下单、发起支付做按 IP 的分钟级限流（Redis 不可用时 fail-open，不阻断业务）

相关开关在 `.env`：

- `CACHE_TTL_SECONDS`
- `RL_ENABLED`、`RL_LOGIN_PER_MIN`、`RL_ORDER_CREATE_PER_MIN`、`RL_PAYMENT_PER_MIN`

## 6. 重置数据（谨慎）

会删除 Postgres/Redis 的持久化数据卷：

```bash
docker compose down -v
```

然后再启动：

```bash
docker compose up -d --build
```

## 7. 后端集成测试（可选）

后端有一组需要真实 Postgres 的集成测试，默认会跳过。

运行方式：

1) 确保本地有一个用于测试的 Postgres（建议单独库/实例）
2) 设置 `TEST_DATABASE_DSN` 后执行：

```bash
cd backend
env GOCACHE=/tmp/go-cache TEST_DATABASE_DSN='host=127.0.0.1 user=... password=... dbname=... port=5432 sslmode=disable' go test ./...
```

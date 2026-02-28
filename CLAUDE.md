# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概览

Order Express 是一个全栈餐厅点餐系统，支持堂食（dine_in）、外送（takeaway）、自取（pickup）三种订单类型。单仓库（monorepo）结构：Go 后端 + Vue 3 前端 + PostgreSQL + Redis，通过 Docker Compose 编排。

## 常用命令

### 全栈启动（Docker）

```bash
cp .env.example .env            # 首次需要
docker compose up -d --build    # 启动所有服务
docker compose down -v          # 重置所有数据（慎用）
```

- 前端：http://localhost:8080
- 后端健康检查：http://localhost:3000/healthz
- 本地默认管理员：`admin` / `admin123`（仅 `DB_AUTO_SEED=true` 时生效）

### 后端（Go 1.25，目录 `backend/`）

```bash
go build ./cmd/server             # 编译
go run ./cmd/server serve         # 启动服务（默认命令）
go run ./cmd/server migrate       # 仅执行数据库迁移
go run ./cmd/server seed          # 仅执行种子数据
go test ./...                     # 运行测试（集成测试需 TEST_DATABASE_DSN）
```

集成测试需要真实 Postgres：
```bash
cd backend
env TEST_DATABASE_DSN='host=127.0.0.1 user=... password=... dbname=... port=5432 sslmode=disable' go test ./...
```

### 前端（Vue 3 + Vite，目录 `frontend/`）

```bash
npm run dev           # 开发服务器（通过 VITE_DEV_PROXY_TARGET 代理 /api 到后端）
npm run build         # 生产构建（含 type-check）
npm run type-check    # 仅 TypeScript 类型检查
npm run lint          # oxlint + eslint（均带 --fix）
npm run format        # prettier 格式化
```

Node 版本要求：`^20.19.0 || >=22.12.0`

### 生产数据库迁移

```bash
docker compose --profile tools run --rm migrate
```

## 架构

### 后端 (`backend/`)

纯标准库 HTTP 服务器（无框架），使用 Go 1.22+ 增强的 `ServeMux` 路由模式。

```
cmd/server/main.go          – 入口；支持 serve|migrate|seed 子命令
internal/config/             – 基于环境变量的配置（config.Load()）
internal/router/router.go   – 所有路由定义集中在一个文件
internal/handlers/           – HTTP handler（按领域一文件一组：auth, orders, payments, dishes, categories, tables, store）
internal/middleware/          – auth（JWT Bearer）、CORS（白名单）、限流（Redis, fail-open）
internal/models/             – GORM 模型
internal/db/                 – 数据库连接、迁移（嵌入式 SQL）、种子数据
internal/db/migrations/      – 原始 SQL 迁移文件
internal/cache/              – Redis 缓存（前缀 key 管理，写操作后失效）
internal/idgen/              – ID 生成（crypto/rand 随机 hex + 前缀）
pkg/jwt/                     – JWT 生成与解析
```

### 前端 (`frontend/src/`)

```
main.ts                – 应用初始化（Pinia, Router, i18n, API client）
router/                – 路由定义，按 admin.ts / customer.ts 分组
stores/                – Pinia 状态管理（auth, app, menu, cart, order, table, settings）
api/                   – HTTP 客户端（api/client.ts）+ 按领域的接口模块；支持 mock 模式
views/                 – 页面组件，按 admin/ 和 customer/ 分组
components/            – 可复用 UI 组件
composables/           – Vue composables（useOrderPolling, useLocaleText, useResponsive）
types/                 – TypeScript 接口定义（按领域一文件）
i18n/                  – 国际化
styles/                – SCSS（变量、Element Plus 覆盖、布局）
utils/                 – 工具函数（storage, QR code, auth）
```

UI 库：Element Plus。状态：Pinia。构建：Vite。通过 unplugin 配置了 Vue 和 Element Plus 的自动导入。

## 关键设计约定

### 数据库：不使用外键

数据库 DDL 中**不创建外键约束**。引用完整性在 handler 逻辑层用事务保证：
- 创建/更新时显式校验引用是否存在（如创建菜品前校验分类存在）
- 删除被引用记录时返回冲突错误（如分类下有菜品则禁止删除）
- 关键写路径使用事务（下单占用桌台、支付更新订单、退款等）

### API 响应格式

统一 JSON 信封：`{"code": number, "data": T, "message": string}`

### 鉴权与授权

- 管理员写接口：`auth(requireAdmin(...))` 中间件链
- 公开读接口（菜单、门店信息、公共桌台列表）：无需鉴权
- 订单创建、支付发起：公开但有限流保护

### 订单类型与必填字段

- `dine_in`：必须传 `tableId`（桌号如 `T05` 或桌台 ID 如 `table-05`），后端自动管理桌台占用/释放
- `takeaway`：必须传 `deliveryAddress` + `contactPhone`
- `pickup`：必须传 `contactPhone`

### 缓存策略

对 categories/dishes/store 等 GET 接口做短 TTL Redis 缓存（`CACHE_TTL_SECONDS` 配置），写操作后主动失效。Redis key 使用 `oe:v1:` 前缀。

## 安全约定

以下安全模式已在代码中建立，修改代码时需维持：

1. **生产环境启动守卫**（`main.go`）：非 `local` 环境下，若 JWT_SECRET 或 DB_PASSWORD 为默认值则拒绝启动；`DB_AUTO_MIGRATE`/`DB_AUTO_SEED` 在非 local 环境不允许开启。
2. **CORS 白名单**（`middleware/cors.go`）：基于 `CORS_ALLOWED_ORIGINS` 白名单匹配，不允许 `credentials + wildcard` 组合（`main.go` 启动时校验）。
3. **HTTP Server 超时**（`main.go`）：已配置 ReadHeaderTimeout/ReadTimeout/WriteTimeout/IdleTimeout/MaxHeaderBytes，防 Slowloris 等慢速攻击。
4. **Nginx 安全头**（`nginx.conf`）：已配置 X-Content-Type-Options、X-Frame-Options、Referrer-Policy、Permissions-Policy。
5. **ID 生成使用 crypto/rand**（`idgen/idgen.go`）：资源 ID 使用加密随机 hex，不可预测。
6. **DevTools 仅开发模式**（`vite.config.ts`）：vueDevTools 仅在 `mode === 'development'` 启用。
7. **登录预填仅开发模式**（`LoginView.vue`）：默认账密仅在 `import.meta.env.DEV` 下预填。
8. **限流 fail-open**（`middleware/ratelimit.go`）：Redis 不可用时不阻断业务，降级放行。

### 已知待改进项

以下来自安全审计报告（`security_best_practices_report.md`），需在后续迭代中关注：

- **F-006**：限流的 `clientIP()` 信任 `X-Forwarded-For`/`X-Real-IP` 头，攻击者可伪造绕过。建议增加"是否信任代理"开关，或仅在确认在可信反代后面时才读取这些头。
- **F-007**：JWT 存储在 localStorage/sessionStorage，存在 XSS 盗取风险。长期可考虑 HttpOnly Cookie + CSRF 方案。
- **F-010**：部分 handler 将内部错误（含 `err.Error()`）直接返回客户端，存在信息泄露。应对外返回统一错误信息，详细错误仅写日志。

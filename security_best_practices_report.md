# Order Express 代码安全审计报告（Go 后端 + Vue 前端）

日期：2026-02-26  
范围：
- `backend/`：Go（`net/http` + GORM/Postgres + Redis）
- `frontend/`：Vue 3 + Vite，Nginx 静态托管并反代 `/api/*`
- 根目录：`docker-compose.yaml`、环境变量文件

> 说明：本报告基于仓库内代码与配置进行静态审计；若某些安全控制在基础设施层（WAF/网关/私网隔离/鉴权代理）实现，代码层面不可见，需在部署侧另行验证。

## 一、执行摘要（Executive Summary）

当前代码存在数个**可直接导致未授权访问/数据篡改**的高风险问题，最重要的是：

1. **后端 API 绝大多数写接口未做鉴权/鉴权中间件未启用（Critical）**：任何能访问后端端口的客户端都可创建/修改/删除菜品、分类、桌台、订单状态、发起支付/退款、修改门店信息等。
2. **CORS 配置“反射 Origin + 允许凭证”（Critical/High）**：在未来若引入 Cookie 会话或浏览器自动携带凭证的场景，可能直接引入跨站风险；同时也会让任意站点跨域读取 API 响应（取决于浏览器请求方式与前端实现）。
3. **默认管理员账号/密码与默认密钥/弱默认配置（High）**：包括前端登录页预填 `admin/admin123`、后端种子数据固定创建该账号，以及 JWT Secret 存在硬编码回退值等；若误用于生产将导致快速被接管。

建议优先级（从高到低）：
1) 给后端 API 建立清晰的 AuthN/AuthZ 边界并启用（鉴权与角色授权）  
2) 修正 CORS（尤其是 `Allow-Credentials` 与 Origin 策略）  
3) 上生产前强制要求显式配置密钥/密码，禁用自动 seed/migrate 与默认账号  
4) 补齐 HTTP Server 超时与请求体大小限制，降低 DoS 风险  

---

## 二、发现项（Findings）

### Critical

#### F-001：后端绝大多数 API 未启用鉴权/授权（AuthN/AuthZ 缺失）

- 严重性：Critical
- 位置：
  - 路由注册未包裹鉴权：`backend/internal/router/router.go:42`（大量 `mux.HandleFunc`/`mux.Handle` 直接暴露）
  - 仅返回 CORS 中间件：`backend/internal/router/router.go:93`
  - 鉴权中间件存在但未使用：`backend/internal/middleware/auth.go:17`（`Auth(...)` 与 `RequireAuth(...)` 定义）
  - 仅 `/api/auth/me` 自行解析 JWT：`backend/internal/handlers/auth.go:83`
- 证据（节选）：
  - `backend/internal/router/router.go:48` 注册 `POST /api/categories`、`DELETE /api/categories/{id}` 等写接口时未做任何鉴权包装
  - `backend/internal/router/router.go:66` 注册 `POST /api/orders` 未鉴权
  - `backend/internal/router/router.go:72` 注册 `POST /api/payments/initiate` 未鉴权
  - `backend/internal/router/router.go:74` 注册 `POST /api/payments/{paymentId}/refund` 未鉴权
- 影响（一句话）：
  - 任意未授权用户可直接调用写接口篡改业务数据、伪造支付/退款、修改订单状态、管理桌台/门店配置等。
- 修复建议（最小变更路径）：
  1. 明确“公开接口”和“管理员接口”的边界（建议：GET 菜单/门店信息公开；写接口与管理类接口需管理员）。
  2. 在路由层对需要保护的接口统一包裹 `middleware.Auth(cfg)`，并增加基于 JWT `claims.Role` 的角色校验（例如 `RequireAdmin`）。
  3. 对订单查询/状态更新、退款等敏感接口做更细粒度授权（至少管理员；若支持顾客查询则需绑定“订单归属”而非仅凭可猜测的 ID）。
- 缓解措施（无法立刻改代码时）：
  - 先在部署侧限制后端端口仅允许内网/Nginx 访问（并避免把后端端口直接暴露给公网）。
- 可能的误报/需验证：
  - 若你在 API 前还有独立鉴权网关/反向代理强制鉴权，请在运行时配置中确认写接口确实被保护；代码层面目前不可见。

#### F-002：CORS 采用“反射 Origin + Allow-Credentials=true”且缺少 Origin 白名单

- 严重性：Critical（对 Cookie/凭证场景）/ High（对纯 Bearer Token 场景）
- 位置：`backend/internal/middleware/cors.go:6`
- 证据（节选）：
  - `backend/internal/middleware/cors.go:8-13`：只要有 `Origin` 就回写 `Access-Control-Allow-Origin: <Origin>` 并设置 `Access-Control-Allow-Credentials: true`
- 影响：
  - 一旦未来改为 Cookie Session（或浏览器会携带某些凭证），该策略可能允许任意站点跨域发起带凭证请求并读取响应，形成严重跨站风险。
  - 即使当前使用 `Authorization: Bearer ...`，该策略也会使 API 对任意 Origin 可读（取决于浏览器请求与前端实现），不利于最小暴露面。
- 修复建议：
  1. 生产环境：如果前端与后端同源（Nginx 反代 `/api`），**可以直接禁用 CORS** 或仅允许同源。
  2. 需要跨域时：引入 `CORS_ALLOWED_ORIGINS`（逗号分隔）白名单，匹配到才回写 `Access-Control-Allow-Origin`。
  3. 若不使用 Cookie：移除 `Access-Control-Allow-Credentials`，避免“任意 Origin + 凭证”组合。
- 缓解措施：
  - 临时仅允许 `http://localhost:5173`（开发）与生产域名；其余 Origin 不回写 CORS 头。
- 可能的误报/需验证：
  - 若在边缘层已做严格 CORS，请验证后端不会覆盖/扩大边缘层策略。

---

### High

#### F-003：默认管理员账号/密码与前端登录页预填密码，存在误部署风险

- 严重性：High
- 位置：
  - 前端登录页预填：`frontend/src/views/admin/LoginView.vue:13`
  - 后端 seed 固定创建管理员密码：`backend/internal/db/seed.go:15`
  - README 公开默认账号密码：`README.md:42`
- 证据（节选）：
  - `frontend/src/views/admin/LoginView.vue:13-17`：`username: 'admin'`、`password: 'admin123'`
  - `backend/internal/db/seed.go:16`：对 `"admin123"` 做 bcrypt 并写入 `users`
- 影响：
  - 一旦 seed 在生产开启或误将默认凭证暴露到公网，攻击者可直接登录获取管理员权限。
- 修复建议：
  1. 生产环境禁用 `DB_AUTO_SEED`，并在启动时校验 `APP_ENV`（非 `local` 时强制 `AutoSeed=false`）。
  2. 前端去掉默认密码预填（可在 `import.meta.env.DEV` 下才填充，或改为示例提示文字）。
  3. 种子管理员密码改为从环境变量读取（例如 `ADMIN_INITIAL_PASSWORD`），且仅在 `local` 执行；生产改为初始化流程/管理后台创建。

#### F-004：JWT/数据库等敏感配置存在硬编码默认值，生产误配将导致可被伪造

- 严重性：High
- 位置：`backend/internal/config/config.go:46`
- 证据（节选）：
  - `backend/internal/config/config.go:50`：`JWT_SECRET` 存在硬编码回退值
  - `backend/internal/config/config.go:57`：`DB_PASSWORD` 存在硬编码回退值
- 影响：
  - 若生产未正确注入环境变量，攻击者可用已知默认密钥伪造 JWT 或利用默认数据库口令。
- 修复建议：
  - 对 `APP_ENV != local`：若 `JWT_SECRET`/`DB_PASSWORD` 等仍为默认值则拒绝启动，并提示正确配置方式。

#### F-005：HTTP 服务未设置超时（ReadHeaderTimeout/ReadTimeout/WriteTimeout 等）

- 严重性：High
- 位置：`backend/cmd/server/main.go:89`
- 证据（节选）：
  - `backend/cmd/server/main.go:90`：`http.ListenAndServe(...)`（默认 `http.Server` 超时为 0）
- 影响：
  - 易受 Slowloris/连接耗尽类 DoS；攻击者可通过慢速发送请求头/请求体占用连接与 goroutine。
- 修复建议：
  - 改用显式 `http.Server{ ... }` 并设置合理的 `ReadHeaderTimeout/ReadTimeout/WriteTimeout/IdleTimeout/MaxHeaderBytes`；必要时增加全局请求体大小限制。

---

### Medium

#### F-006：限流使用 `X-Forwarded-For`/`X-Real-IP` 作为客户端 IP，易被伪造绕过

- 严重性：Medium
- 位置：
  - 客户端 IP 解析：`backend/internal/middleware/ratelimit.go:67`
  - 后端端口直出：`docker-compose.yaml:41`
- 证据（节选）：
  - `backend/internal/middleware/ratelimit.go:68-77`：优先读取 `X-Forwarded-For` / `X-Real-IP`
  - `docker-compose.yaml:41`：后端 `3000` 映射到宿主机（可绕过 Nginx 直接访问）
- 影响：
  - 攻击者可构造请求头伪造 IP，导致限流 key 可控，从而绕过限流。
- 修复建议：
  - 仅在“确认在可信反向代理后面”时才信任 `X-Forwarded-For`；否则只使用 `RemoteAddr`。
  - 或在 Nginx 层剥离来自客户端的 `X-Forwarded-For` 并重新设置，但后端仍应做“是否信任代理”的显式开关。

#### F-007：JWT 存储于 `localStorage/sessionStorage`，存在 XSS 盗取令牌风险

- 严重性：Medium
- 位置：`frontend/src/utils/authStorage.ts:17`
- 证据（节选）：
  - `frontend/src/utils/authStorage.ts:44-46`：`storage.setItem(AUTH_TOKEN_KEY, token)`
- 影响：
  - 一旦出现 XSS，攻击者可读取 Storage 中的 token 并长期冒用身份（尤其是 `localStorage` 持久化）。
- 修复建议：
  - 推荐改为 HttpOnly Cookie + CSRF 方案（需要后端配合）。
  - 若短期无法改造：优先使用内存存储 + 短有效期 token，减少持久化；并在部署层加 CSP/点击劫持防护等。

#### F-008：Nginx 未设置基础安全响应头（CSP/Clickjacking/Nosniff 等）

- 严重性：Medium
- 位置：`frontend/nginx.conf:1`
- 证据（节选）：
  - 配置中仅包含静态托管与反代逻辑，未见 `add_header ...`（`frontend/nginx.conf:1-19`）
- 影响：
  - 浏览器端缺少防护“兜底”，对 XSS/点击劫持等攻击的抵抗力下降（尤其在未来引入第三方脚本/富文本时）。
- 修复建议：
  - 在 Nginx 添加基础安全头：`X-Content-Type-Options: nosniff`、`X-Frame-Options: DENY`（或 CSP `frame-ancestors`）、`Referrer-Policy`、`Permissions-Policy`。
  - CSP 需结合实际资源加载策略逐步收紧，避免一次性上强 CSP 造成页面不可用。

#### F-009：资源 ID/订单号/支付号大量使用时间戳拼接，可预测并易于枚举

- 严重性：Medium（与 F-001 叠加时可升级为 Critical 的“自动化爬取/篡改”）
- 位置（示例）：
  - 分类 ID：`backend/internal/handlers/categories.go:109`
  - 订单 ID/订单号：`backend/internal/handlers/orders.go:278`
  - 支付 ID：`backend/internal/handlers/payments.go:75`
  - 桌台 ID：`backend/internal/handlers/tables.go:105`
- 影响：
  - 攻击者可推测资源数量与创建时间，并可批量枚举资源（尤其当接口缺少鉴权时）。
- 修复建议：
  - 对外暴露的资源 ID 使用 UUIDv4/随机字符串；订单展示号可保留业务含义但不应作为授权凭据。

#### F-010：部分接口将内部错误信息直接返回给客户端（信息泄露）

- 严重性：Medium
- 位置：
  - `backend/internal/handlers/orders.go:399`
  - `backend/internal/handlers/tables.go:117`
- 证据（节选）：
  - `Fail(..., "database error: "+err.Error())`
- 影响：
  - 可能泄露数据库结构/约束/内部错误细节，为攻击者提供探测信息。
- 修复建议：
  - 对外返回统一错误信息；详细错误仅记录到服务端日志（注意日志脱敏）。

---

### Low

#### F-011：Vite DevTools 插件在构建配置中默认启用，建议仅在开发启用

- 严重性：Low（取决于插件行为）
- 位置：`frontend/vite.config.ts:15`
- 证据（节选）：
  - `frontend/vite.config.ts:15-27`：`vueDevTools()` 在所有模式下加入插件列表
- 影响：
  - 可能引入额外调试入口/元数据暴露或增加依赖面；至少会增加生产包体积与复杂度。
- 修复建议：
  - 仅在 `mode === 'development'`（或 `command === 'serve'`）启用 devtools 插件。

#### F-012：`backend/.gitignore` 规则可能导致审计/检索工具漏扫关键目录

- 严重性：Low
- 位置：`backend/.gitignore:4`
- 证据（节选）：
  - `backend/.gitignore:4`：`server`（会匹配 `backend/cmd/server/` 目录名）
- 影响：
  - 像 `rg` 等默认遵循 gitignore 的工具在全仓扫描时可能跳过 `backend/cmd/server/`，导致安全审计与代码搜索漏报。
- 修复建议：
  - 将规则改为更精确的路径（例如仅忽略根目录二进制 `/server`），避免忽略源码目录。

---

## 三、后续建议（Next Steps）

1. 先处理 F-001/F-002：把后端 API 的“公开/受保护”边界落地，并修正 CORS 策略。  
2. 处理 F-003/F-004：生产启动前强制校验密钥与默认账号策略，禁用自动 seed/migrate。  
3. 处理 F-005/F-006/F-010：补齐超时、请求大小限制、可信代理与错误处理，降低 DoS 与信息泄露风险。  
4. 处理前端存储与部署安全头（F-007/F-008/F-011），做防御纵深。  


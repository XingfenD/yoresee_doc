# TODO / 待办

Known gaps and pending work, ordered by severity. Not a changelog — only the changelog records shipped changes.

已知缺陷与待办项，按严重程度排列。这不是更新日志，只有 changelog 记录已发布变更。

## 1. 鉴权兜底分支为后门 / Auth fallback is a backdoor — `collab-go/auth/auth.go:31`

- 当 `JWT_SECRET` 为空时走 `ParseUnverified`，任何伪造 JWT 都通过；生产已配 secret，不常触发，但是潜伏后门。
- When `JWT_SECRET` is empty, code falls to `ParseUnverified` — any forged JWT passes. Production sets the secret, so rarely triggered, but a latent backdoor.

## 2. 无按用户授权 + 无身份传递 / No per-user authz + no identity propagation — `collab-go/handler/ws.go:84`

- `checkDocumentExists` 只查文档是否存在，不查当前用户是否有权访问；后端 gRPC 不可用时 `return 0, nil` 静默放行，任何人可连任意文档（含不存在的，生成幽灵房间）。
- 网关不解析也不向下游传递用户身份（claims）。
- `checkDocumentExists` only checks existence, not the caller's access right; on backend gRPC failure it returns `0, nil` and silently allows anyone to open any doc (incl. non-existent → phantom rooms). The gateway neither parses nor forwards user identity/claims.

## 3. 代理缺保活/优雅关闭 / Proxy lacks keepalive & graceful close — `collab-go/proxy/proxy.go:50`

- 无 ping-pong / read deadline：半开连接泄漏 goroutine 与 FD。
- 无 `SetMaxMessageSize`：大消息 DoS 风险。
- dial core 失败时直接 `conn.Close()`，不发 WebSocket Close frame，客户端看到突断。
- No ping-pong / read deadline: half-open connections leak goroutines and FDs. No `SetMaxMessageSize` (large-message DoS). On core dial failure it closes the socket without a Close frame, so the client sees an abrupt drop.

## 4. collab-core 零鉴权 + 死代码 / collab-core trusts gateway blindly + dead code — `collab/src/ws-gateway.js:6`

- `bindWebSocketGateway` 零鉴权，完全信任网关；绕过网关可直接开任意房间。
- `saveDocumentYjsSnapshot` 已实现却从未被调用（死代码）。
- `bindWebSocketGateway` performs zero auth and fully trusts the gateway; bypassing the gateway lets anyone open any room. `saveDocumentYjsSnapshot` is implemented but never invoked (dead code).

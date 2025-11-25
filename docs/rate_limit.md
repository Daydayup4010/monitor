# API 接口限流说明

## 限流方式

- **IP 限流**：基于客户端 IP 地址限制，适用于公开接口（无需登录）
- **用户限流**：基于用户 ID 限制，适用于需要认证的接口

## 限流算法

采用 **滑动窗口算法**，基于 Redis 实现，支持分布式部署。

## 响应头

每个请求的响应会包含以下限流信息：

| Header | 说明 |
|--------|------|
| `X-RateLimit-Limit` | 时间窗口内最大请求数 |
| `X-RateLimit-Remaining` | 剩余可用请求数 |
| `X-RateLimit-Reset` | 限流重置时间（Unix 时间戳） |

## 被限流时的响应

HTTP 状态码：`429 Too Many Requests`

```json
{
  "error": "请求过于频繁，请稍后再试",
  "code": "RATE_LIMIT_EXCEEDED",
  "retryAfter": 45
}
```

---

## 接口限流配置

### 用户模块 `/api/v1/user`

| 接口 | 方法 | 限流方式 | 频率限制 | 说明 |
|------|------|----------|----------|------|
| `/register` | POST | IP | 5 次/分钟 | 防止批量注册 |
| `/login` | POST | IP | 10 次/分钟 | 防止暴力破解 |
| `/email-login` | POST | IP | 10 次/分钟 | 防止暴力破解 |
| `/send-email` | POST | IP | 3 次/分钟 | 防止邮件滥发 |
| `/reset-password` | POST | IP | 5 次/小时 | 敏感操作，严格限制 |
| `/email-exist` | POST | - | 无限制 | - |
| `/self` | GET | - | 无限制 | 需认证 |
| `/name` | PUT | - | 无限制 | 需认证 |
| `/refresh` | POST | - | 无限制 | 需认证 |

### VIP 模块 `/api/v1/vip`

| 接口 | 方法 | 限流方式 | 频率限制 | 说明 |
|------|------|----------|----------|------|
| `/goods/data` | GET | 用户 | 30 次/分钟 | 商品数据查询 |
| `/goods/category` | GET | 用户 | 30 次/分钟 | 商品分类查询 |
| `/settings` | GET | - | 无限制 | 获取设置 |
| `/settings` | PUT | - | 无限制 | 更新设置 |

### 管理员模块 `/api/v1/admin`

| 接口 | 方法 | 限流方式 | 频率限制 | 说明 |
|------|------|----------|----------|------|
| `/users` | GET | 用户 | 60 次/分钟 | 获取用户列表 |
| `/user` | DELETE | 用户 | 60 次/分钟 | 删除用户 |
| `/vip-expiry` | POST | 用户 | 60 次/分钟 | 续费 VIP |
| `/tokens/uu` | POST | 用户 | 60 次/分钟 | 更新 UU Token |
| `/tokens/buff` | POST | 用户 | 60 次/分钟 | 更新 Buff Token |
| `/tokens/verify` | POST | 用户 | 60 次/分钟 | 验证 Token |
| `/tokens/verify` | GET | 用户 | 60 次/分钟 | 获取验证状态 |

### 微信模块 `/api/v1/wechat`

| 接口 | 方法 | 限流方式 | 频率限制 | 说明 |
|------|------|----------|----------|------|
| `/login` | POST | IP | 10 次/分钟 | 微信登录 |
| `/bind-email` | POST | 用户 | 5 次/分钟 | 绑定邮箱 |
| `/merge-account` | POST | 用户 | 3 次/小时 | 合并账号，敏感操作 |
| `/bind-wechat` | POST | 用户 | 5 次/分钟 | 绑定微信 |
| `/send-email-code` | POST | 用户 | 3 次/分钟 | 发送验证码 |

---

## 配置修改

限流配置位于 `core/router.go`，可根据实际需求调整：

```go
middleware.RateLimiterByIP(middleware.RateLimiterConfig{
    Window:      60 * time.Second,  // 时间窗口
    MaxRequests: 10,                // 最大请求数
    KeyPrefix:   "user:login",      // Redis 键前缀
})
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `Window` | `time.Duration` | 时间窗口，如 `60 * time.Second` 表示 1 分钟 |
| `MaxRequests` | `int` | 窗口内最大请求次数 |
| `KeyPrefix` | `string` | Redis 键前缀，用于区分不同接口 |


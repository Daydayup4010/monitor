# VIP 会员支付流程

本文档说明 CS Goods VIP 会员支付的完整流程，使用 YunGouOS 微信 Native 支付。

## 支付流程概览

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   前端页面   │────▶│   后端API   │────▶│  YunGouOS   │────▶│   微信支付   │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
      │                   │                                        │
      │                   │◀───────────── 二维码 ──────────────────┘
      │◀── 显示二维码 ────┘
      │
      │  用户扫码支付
      │
      │                   │◀─────────── 支付回调 ──────────────────┘
      │◀── 轮询订单状态 ──┘
      │
      ▼
   支付成功
```

## 详细步骤

### 1. 用户选择套餐

前端页面 (`Settings.vue`) 提供以下 VIP 套餐选项：

| 时长 | 价格 | 折扣 | VIP有效期 |
|------|------|------|----------|
| 1个月 | ￥19.9 | - | +1个月 |
| 3个月 | ￥49.9 | 8.4折 | +3个月 |
| 6个月 | ￥89.9 | 7.5折 | +6个月 |
| 12个月 | ￥169.9 | 7.1折（推荐） | +12个月 |

### 2. 创建支付订单

**前端请求：**
```http
POST /api/v1/vip/payment/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "months": 12
}
```

**后端处理 (`api/payment.go`)：**
1. 验证用户身份
2. 根据 months 获取套餐价格
3. 生成商户订单号：`VIP + 时间戳 + 随机数`
4. 在数据库创建订单记录
5. 调用 YunGouOS Native Pay 接口

**返回数据：**
```json
{
  "code": 1,
  "data": {
    "order_no": "VIP20251218120000001",
    "qrcode_img": "data:image/png;base64,iVBORw0KGgo...",
    "amount": 169.9,
    "created_at": "2025-12-18T12:00:00Z"
  }
}
```

### 3. YunGouOS Native Pay 签名规则

根据 [YunGouOS 签名文档](https://open.pay.yungouos.com/#/api/sign)：

**参与签名的参数：**
- `out_trade_no` - 商户订单号
- `total_fee` - 支付金额
- `mch_id` - 商户号
- `body` - 商品描述
- `attach` - 附加数据（有值时参与）

**不参与签名的参数：**
- `type` - 二维码类型
- `notify_url` - 回调地址
- `sign` - 签名本身

**签名算法：**
1. 将参与签名的参数按 ASCII 码排序
2. 拼接成 `key1=value1&key2=value2` 格式
3. 末尾追加 `&key=商户密钥`
4. 对整个字符串进行 MD5 加密
5. 转换为大写

```go
// 示例：services/payment.go
signParams := map[string]string{
    "out_trade_no": "VIP20251218120000001",
    "total_fee":    "169.90",
    "mch_id":       "YG1234567890",
    "body":         "CS Goods VIP会员12个月",
}
sign := PayGenerateSign(signParams, apiKey)
```

### 4. 显示支付二维码

前端收到 `qrcode_img`（Base64 格式的二维码图片），直接显示给用户扫码。

```vue
<img :src="currentOrder.qrcode_img" alt="微信支付二维码" />
```

### 5. 轮询订单状态

前端每 3 秒查询一次订单状态：

```http
GET /api/v1/vip/payment/query?order_no=VIP20251218120000001
Authorization: Bearer <token>
```

**返回数据：**
```json
{
  "code": 1,
  "data": {
    "order_no": "VIP20251218120000001",
    "amount": 169.9,
    "status": 0,  // 0=待支付, 1=已支付
    "pay_time": null,
    "created_at": "2025-12-18T12:00:00Z"
  }
}
```

### 6. 支付回调通知

用户支付成功后，YunGouOS 会调用配置的 `notify_url`：

```http
POST /api/v1/open/payment/notify
Content-Type: application/x-www-form-urlencoded

code=1&orderNo=YG1234567890&outTradeNo=VIP20251218120000001&payNo=wxpay123456&money=169.90&mchId=YG1234567890&payChannel=wxpay&time=2025-12-18 12:05:00&attach=userid|12&sign=ABC123...
```

**回调参数签名验证：**

根据 [YunGouOS 回调文档](https://open.pay.yungouos.com/#/callback/notify)：

| 参数 | 是否参与签名 |
|------|------------|
| code | ✅ 是 |
| orderNo | ✅ 是 |
| outTradeNo | ✅ 是 |
| payNo | ✅ 是 |
| money | ✅ 是 |
| mchId | ✅ 是 |
| payChannel | ✅ 是 |
| attach | ✅ 是（有值时） |
| time | ❌ 否 |
| sign | ❌ 否 |

**后端处理 (`api/payment.go` - `PayNotify`)：**
1. 验证签名
2. 检查 `code == "1"` 表示支付成功
3. 查询订单，防止重复处理
4. 更新订单状态为已支付
5. 延长用户 VIP 到期时间
6. 返回 `SUCCESS` 给 YunGouOS

### 7. 支付成功

- 前端轮询检测到 `status == 1`
- 停止轮询，关闭支付弹窗
- 显示支付成功提示
- 刷新用户信息，更新 VIP 状态

## 配置说明

### settings.yaml

```yaml
payment:
  mch_id: "YG1234567890"      # YunGouOS 商户号
  api_key: "your_api_key"     # YunGouOS 商户密钥
  notify_url: "https://your-domain.com/api/v1/open/payment/notify"  # 回调地址
```

### 数据库表结构

```sql
CREATE TABLE payment_orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(64) NOT NULL,
    out_trade_no VARCHAR(64) NOT NULL UNIQUE,  -- 商户订单号
    order_no VARCHAR(64),                       -- YunGouOS 订单号
    amount DECIMAL(10,2) NOT NULL,              -- 支付金额
    months INT NOT NULL DEFAULT 1,              -- 购买月数
    status TINYINT NOT NULL DEFAULT 0,          -- 0=待支付, 1=已支付
    pay_time DATETIME,                          -- 支付时间
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_out_trade_no (out_trade_no)
);
```

## 相关文件

| 文件 | 说明 |
|------|------|
| `services/payment.go` | YunGouOS 支付服务封装 |
| `api/payment.go` | 支付相关 API 接口 |
| `models/payment.go` | 订单模型和数据库操作 |
| `FE/src/pages/Settings.vue` | 前端支付页面 |
| `FE/src/api/index.ts` | 前端 API 定义 |

## 参考文档

- [YunGouOS Native Pay](https://open.pay.yungouos.com/#/api/api/pay/wxpay/nativePay)
- [YunGouOS 签名规则](https://open.pay.yungouos.com/#/api/sign)
- [YunGouOS 回调通知](https://open.pay.yungouos.com/#/callback/notify)


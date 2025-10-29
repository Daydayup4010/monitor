# Goods - 微信小程序

## 📱 项目简介

商品差价监控系统的微信小程序版本，支持微信一键登录，查看商品差价数据。

## 🎯 功能特性

### 核心功能
- ✅ 微信一键登录
- ✅ 查看Buff ↔ UU饰品差价
- ✅ 平台切换（Buff→UU / UU→Buff）
- ✅ 实时数据刷新
- ✅ VIP权限管理

### 账号互通
- ✅ Web端用户可绑定微信
- ✅ 小程序用户可绑定邮箱
- ✅ 共享VIP状态和筛选设置

## 📁 项目结构

```
miniprogram/
├── app.js                  # 小程序入口
├── app.json                # 全局配置
├── app.wxss                # 全局样式
├── project.config.json     # 项目配置
├── utils/                  # 工具类
│   ├── request.js          # 网络请求封装
│   └── api.js              # API接口定义
├── pages/                  # 页面
│   ├── login/              # 登录页面
│   ├── home/               # 饰品数据首页
│   ├── my/                 # 我的页面
│   └── bind-email/         # 绑定邮箱页面
└── images/                 # 图片资源
```

## 🚀 快速开始

### 1. 配置后端

修改 `settings.yaml`，填写微信小程序配置：

```yaml
wechat:
  app_id: "你的小程序AppID"
  app_secret: "你的小程序AppSecret"
```

### 2. 数据库迁移

User表已自动添加字段：
- `wechat_openid` - 微信OpenID
- `wechat_unionid` - 微信UnionID（可选）

### 3. 配置小程序

修改 `project.config.json`：
```json
{
  "appid": "你的小程序AppID"
}
```

修改 `app.js` 中的baseURL：
```javascript
baseURL: 'https://your-domain.com/api/v1'  // 你的后端地址
```

### 4. 导入项目

1. 打开微信开发者工具
2. 导入 `miniprogram` 文件夹
3. 填写AppID
4. 编译运行

## 📋 API接口

### 微信登录相关
- `POST /api/v1/wechat/login` - 微信登录
- `POST /api/v1/wechat/bind-email` - 绑定邮箱
- `POST /api/v1/wechat/bind-wechat` - Web用户绑定微信
- `POST /api/v1/wechat/send-email-code` - 发送验证码

### 数据相关（复用现有API）
- `GET /api/v1/vip/goods/data` - 获取饰品数据
- `GET /api/v1/vip/goods/category` - 获取分类列表
- `GET /api/v1/vip/settings` - 获取个人设置

## 🔄 用户互通逻辑

### 场景1：Web用户 → 小程序
1. Web端邮箱注册/登录
2. 小程序微信登录
3. 提示绑定已有账号
4. 输入邮箱+验证码绑定
5. 完成互通

### 场景2：小程序用户 → Web端
1. 小程序微信登录
2. 提示绑定邮箱
3. 输入邮箱+验证码+密码
4. 可用邮箱登录Web端
5. 完成互通

## 🎨 设计风格

- 延续Web端的极简设计
- 淡灰色背景（#f5f7fa）
- 白色卡片
- 蓝色主题色（#1890ff）

## 📝 待完善功能

- [ ] 绑定邮箱页面（已创建文件，待实现）
- [ ] 个人设置页面
- [ ] 搜索和分类筛选
- [ ] 收藏功能
- [ ] 通知提醒

## 🔧 注意事项

1. **AppID和AppSecret**：需要在微信公众平台申请小程序
2. **服务器域名配置**：在小程序后台配置request合法域名
3. **HTTPS**：小程序要求后端必须使用HTTPS
4. **Token过期处理**：已在request.js中实现自动跳转登录

## 📞 技术支持

如有问题，请查看后端日志或前端控制台输出。


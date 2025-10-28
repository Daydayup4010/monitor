# 图标使用指南

## 📁 图标文件位置

所有图标配置在：`FE/src/config/icons.ts`

## 🎯 当前使用的图标列表

### 页面标题图标

| 位置 | 当前图标 | 文件 | 代码位置 |
|------|---------|------|---------|
| 登录页面标题 | 🔐 | `Login.vue` | `<h2 class="auth-title">🔐 登录</h2>` |
| 注册页面标题 | ✨ | `Register.vue` | `<h2 class="auth-title">✨ 注册</h2>` |
| 重置密码标题 | 🔑 | `ResetPassword.vue` | `<h2 class="auth-title">🔑 重置密码</h2>` |
| 饰品数据标题 | 📊 | `Home.vue` | `<div class="card-title">📊 饰品数据</div>` |
| 筛选参数标题 | ⚙️ | `Settings.vue` | `<div class="card-title">⚙️ 筛选参数配置</div>` |
| 用户管理标题 | 👥 | `UserManager.vue` | `<div class="card-title">👥 用户管理</div>` |
| Token管理标题 | 🔑 | `TokenManager.vue` | `<div class="card-title">🔑 Token管理</div>` |

### 左侧菜单图标

| 位置 | 当前图标 | 文件 | 代码位置 |
|------|---------|------|---------|
| 用户管理菜单 | 👥 | `Admin.vue` | `<div class="menu-icon">👥</div>` |
| Token管理菜单 | 🔑 | `Admin.vue` | `<div class="menu-icon">🔑</div>` |

### 统计卡片图标

| 位置 | 当前图标 | 文件 | 代码位置 |
|------|---------|------|---------|
| 总用户数 | 👥 | `UserManager.vue` | `<div class="stat-icon blue">👥</div>` |
| VIP用户 | 👑 | `UserManager.vue` | `<div class="stat-icon green">👑</div>` |

### 用户徽章图标

| 位置 | 当前图标 | 文件 | 代码位置 |
|------|---------|------|---------|
| 管理员 | 👨‍💼 | `Settings.vue` | computed `badgeIcon` |
| VIP会员 | 👑 | `Settings.vue` | computed `badgeIcon` |
| 普通用户 | 👤 | `Settings.vue` | computed `badgeIcon` |

### 功能图标

| 位置 | 当前图标 | 文件 | 代码位置 |
|------|---------|------|---------|
| 刷新按钮 | 🔄 | `Home.vue` | `<span style="font-size: 16px;">🔄</span>` |
| 搜索框 | 🔍 | `UserManager.vue` | placeholder中 |
| VIP提示 | ⚠️ | `Settings.vue` | `<div class="notice-icon">⚠️</div>` |

---

## 🔧 如何替换图标

### 方式1：替换为自定义SVG/PNG图标

1. **将图标文件放到 `FE/src/assets/icons/` 文件夹**
   ```
   FE/src/assets/icons/
     ├─ login.svg
     ├─ register.svg
     ├─ user.svg
     └─ ...
   ```

2. **修改对应组件**

   原代码：
   ```vue
   <h2 class="auth-title">🔐 登录</h2>
   ```

   改为：
   ```vue
   <h2 class="auth-title">
     <img src="@/assets/icons/login.svg" style="width: 24px; height: 24px; vertical-align: middle;" />
     登录
   </h2>
   ```

### 方式2：使用Element Plus图标

原代码：
```vue
<div class="stat-icon blue">👥</div>
```

改为：
```vue
<div class="stat-icon blue">
  <el-icon><User /></el-icon>
</div>
```

### 方式3：直接替换Emoji

在对应文件中直接修改emoji字符：
```vue
🔐 → 🎮  （或其他emoji）
```

---

## 📝 快速替换列表

如果你要替换，请告诉我：

```
login.svg → 替换登录页面的🔐
user.svg → 替换用户管理的👥
...
```

或者直接把图标文件放到 `FE/src/assets/icons/` 文件夹，告诉我文件名和对应位置，我帮你修改代码！

---

## 🎨 推荐的图标尺寸

- 页面标题图标：24-28px
- 菜单图标：20-24px
- 统计卡片图标：24px
- 功能按钮图标：16-20px




# WeUI 组件库集成指南

## 📦 什么是WeUI

WeUI是微信官方设计的UI组件库，提供标准的小程序组件样式。

## 🚀 集成方式

### 方式1：使用npm（推荐）

1. **在miniprogram文件夹下初始化npm**
```bash
cd miniprogram
npm init -y
npm install weui-miniprogram --save
```

2. **在微信开发者工具中构建npm**
- 工具 → 构建npm
- 等待构建完成

3. **在页面json中引入组件**
```json
{
  "usingComponents": {
    "mp-cell": "weui-miniprogram/cell/cell",
    "mp-cells": "weui-miniprogram/cells/cells",
    "mp-slideview": "weui-miniprogram/slideview/slideview"
  }
}
```

### 方式2：直接使用（当前推荐）

我已经为你优化了组件样式，符合微信设计规范，不需要额外安装WeUI。

当前组件特点：
- ✅ 符合微信小程序设计规范
- ✅ 简洁现代的风格
- ✅ 与Web端设计统一
- ✅ 无需额外依赖

## 📋 当前已有的组件

### 登录页
- 登录卡片
- 微信登录按钮

### 饰品数据页
- 筛选卡片（平台、类别、排序）
- 饰品列表卡片
- 空状态
- 加载更多

### 我的页面
- 用户信息卡片
- 菜单列表

### 绑定邮箱页
- 表单组件
- 输入框
- 验证码按钮
- 提交按钮

## 🎨 组件风格

- 圆角：12-24rpx
- 阴影：0 4rpx 16rpx
- 间距：统一的16/24/32rpx体系
- 颜色：与Web端一致

## 💡 建议

当前不需要安装WeUI，现有组件已经足够美观和功能完整。如果将来需要更复杂的组件（如日历、图片上传等），再考虑引入WeUI。


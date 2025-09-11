#!/bin/bash

# 清理并重新安装依赖的脚本
# 适用于 Node.js 18.16.0 环境

echo "清理现有的 node_modules..."
rm -rf node_modules

echo "清理 npm 缓存..."
npm cache clean --force

echo "安装依赖（使用精确版本）..."
npm install

echo "验证 Vite 版本..."
npm list vite

echo "完成！现在可以运行 npm run dev"


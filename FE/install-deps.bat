@echo off
echo 清理现有的 node_modules...
rmdir /s /q node_modules 2>nul

echo 清理 npm 缓存...
npm cache clean --force

echo 安装依赖（使用精确版本）...
npm install

echo 验证 Vite 版本...
npm list vite

echo 完成！现在可以运行 npm run dev
pause


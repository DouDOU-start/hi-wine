# hi-wine

这是一个基于 [uni-app](https://uniapp.dcloud.io/) 的多端项目，支持微信小程序、H5等平台。

## 项目结构
- `pages/`：页面文件夹
- `static/`：静态资源
- `resources/`：资源文件
- `main.js`：入口文件
- `App.vue`：主组件
- `manifest.json`：项目配置
- `package.json`：依赖管理

## 快速开始

### 1. 安装依赖
```bash
npm install
```

### 2. 本地开发
使用 HBuilderX 或命令行工具运行：
```bash
# H5 预览
npm run dev:h5

# 微信小程序预览
npm run dev:mp-weixin
```

### 3. 构建发布
```bash
# 构建 H5
npm run build:h5

# 构建微信小程序
npm run build:mp-weixin
```

## 图片基础地址配置

如需统一管理图片 CDN 或静态资源基础路径，请在 `config.js` 中配置：

```js
export const IMG_BASE_URL = 'https://your-img-base-url.com'; // 替换为你的图片基础地址
```

在项目中可通过 `import { IMG_BASE_URL } from '@/config.js'` 直接使用。

## 后端（backend）服务说明

本项目后端基于 GoFrame 框架开发，目录为 `backend/`。

### 1. 依赖环境
- Go 1.18 及以上
- MySQL 8.0 及以上

### 2. 安装依赖
进入 backend 目录，执行：
```bash
cd backend
go mod tidy
```

### 3. 数据库初始化
使用 `backend/db_schema.sql` 初始化数据库：
```bash
mysql -u root -p your_db < db_schema.sql
```

### 4. 配置
编辑 `backend/config.yml`，填写微信小程序 AppID、Secret 及数据库等信息。

### 5. 启动后端服务
```bash
go run main.go
# 或
make run
```
默认监听 8000 端口。

### 6. 主要接口
- `/api/wechat/login`：微信小程序登录
- 其它接口详见代码

## 其他说明
- 需要 Node.js 环境。
- 推荐使用 [HBuilderX](https://www.dcloud.io/hbuilderx.html) 进行可视化开发和调试。
- 如需本地静态资源服务，可运行 `start-static-server.js`。

## License
MIT 
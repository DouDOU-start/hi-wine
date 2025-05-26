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

## 其他说明
- 需要 Node.js 环境。
- 推荐使用 [HBuilderX](https://www.dcloud.io/hbuilderx.html) 进行可视化开发和调试。
- 如需本地静态资源服务，可运行 `start-static-server.js`。

## License
MIT 
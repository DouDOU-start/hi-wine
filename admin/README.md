# 酒馆后台管理系统

这是一个基于Vue 3、Element Plus和Vite构建的酒馆后台管理系统，用于管理酒馆小程序的用户、商品、分类和订单等数据。

## 功能特性

- 用户管理：查看用户列表、用户详情、修改用户状态等
- 商品分类管理：添加、编辑、删除商品分类
- 商品管理：添加、编辑、删除商品，上下架商品
- 订单管理：查看订单列表、订单详情，处理订单状态
- 数据统计：销售额、订单数、商品数、用户数等数据统计

## 技术栈

- Vue 3：前端框架
- Vue Router：路由管理
- Element Plus：UI组件库
- Axios：HTTP请求
- ECharts：数据可视化图表
- Vite：构建工具

## 开发环境要求

- Node.js >= 14.0.0
- npm >= 6.0.0

## 安装与运行

1. 安装依赖
```bash
npm install
```

2. 开发模式运行
```bash
npm run dev
```

3. 构建生产版本
```bash
npm run build
```

4. 预览生产版本
```bash
npm run preview
```

## 项目结构

```
admin/
├── public/             # 静态资源
├── src/                # 源代码
│   ├── api/            # API接口
│   ├── assets/         # 资源文件
│   ├── components/     # 公共组件
│   ├── router/         # 路由配置
│   ├── utils/          # 工具函数
│   ├── views/          # 页面组件
│   ├── App.vue         # 根组件
│   └── main.js         # 入口文件
├── index.html          # HTML模板
├── vite.config.js      # Vite配置
└── package.json        # 项目依赖
```

## 接口配置

后端接口基础路径在`vite.config.js`中配置：

```js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
  },
},
```

## 开发指南

1. 页面组件放在`src/views/`目录下
2. API接口定义在`src/api/`目录下
3. 公共组件放在`src/components/`目录下
4. 工具函数放在`src/utils/`目录下

## 部署

构建生产版本后，将`dist`目录下的文件部署到Web服务器即可。 
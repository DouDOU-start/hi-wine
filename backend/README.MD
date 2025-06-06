# Hi-Wine 酒水销售后台管理系统

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.16%2B-blue" alt="Go版本">
  <img src="https://img.shields.io/badge/GoFrame-v2-brightgreen" alt="GoFrame版本">
  <img src="https://img.shields.io/badge/MySQL-5.7%2B-orange" alt="MySQL版本">
  <img src="https://img.shields.io/badge/MinIO-最新版-yellow" alt="MinIO版本">
  <img src="https://img.shields.io/badge/微信小程序-接入-brightgreen" alt="微信小程序">
</p>

## 项目介绍

Hi-Wine 是一个基于 GoFrame 框架开发的酒水销售后台管理系统，主要用于酒吧、酒馆等场所的酒水销售和管理。系统分为管理后台和小程序客户端两部分，通过 RESTful API 进行通信。

## 功能特点

- **畅饮套餐**：用户可以购买畅饮套餐，在有效期内享受套餐内指定酒水的免费服务
- **扫码点单**：通过扫描桌号二维码进行点单，提升用户体验
- **微信支付**：集成微信支付功能，便捷完成订单支付
- **库存管理**：实时跟踪酒水库存，自动扣减库存
- **订单管理**：全流程订单处理，包括创建、支付、配送和完成
- **用户管理**：微信小程序用户管理，包括登录、注册和个人信息维护
- **图片存储**：使用 MinIO 对象存储服务管理商品图片

## 技术栈

- **后端框架**：GoFrame v2
- **数据库**：MySQL 5.7+
- **对象存储**：MinIO
- **认证方式**：JWT
- **微信接入**：微信小程序、微信支付

## 快速开始

### 环境要求

- Go 1.16+
- MySQL 5.7+
- MinIO 服务

### 安装步骤

1. 克隆项目
   ```bash
   git clone https://github.com/DouDOU-start/hi-wine
   cd hi-wine/backend
   ```

2. 安装依赖
   ```bash
   go mod tidy
   ```

3. 配置数据库
   - 创建数据库 `hi-wine`
   - 导入 `resource/schema.sql` 初始化表结构

4. 配置应用
   - 复制 `manifest/config/config.example.yaml` 为 `manifest/config/config.yaml`
   - 根据实际情况修改配置参数

5. 启动服务
   ```bash
   go run main.go
   ```

6. 访问接口文档
   ```
   http://localhost:8000/swagger
   ```

## 目录结构

```
├── api/                # API 定义
├── internal/           # 内部实现
│   ├── cmd/            # 命令行入口
│   ├── controller/     # 控制器
│   ├── dao/            # 数据访问对象
│   ├── logic/          # 业务逻辑
│   ├── model/          # 数据模型
│   ├── service/        # 服务接口
│   ├── middleware/     # 中间件
│   └── utility/        # 工具函数
├── manifest/           # 配置文件
│   └── config/         # 应用配置
├── resource/           # 静态资源
└── main.go             # 应用入口
```

## 配置说明

### 数据库配置

在 `manifest/config/config.yaml` 文件中配置数据库连接：

```yaml
# 数据库配置
database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/hi-wine?loc=Local&parseTime=true"
```

### MinIO 配置

在 `manifest/config/config.yaml` 文件中配置 MinIO 相关参数：

```yaml
# MinIO配置
minio:
  endpoint: "127.0.0.1:9000"  # MinIO服务地址
  accessKey: "minioadmin"     # 访问密钥
  secretKey: "minioadmin"     # 密钥
  useSSL: false               # 是否使用SSL
  bucketName: "hi-wine"       # 存储桶名称
  region: ""                  # 区域，可为空
  domain: ""                  # 访问域名，留空使用默认URL格式
```

### JWT 配置

```yaml
# JWT配置
jwt:
  secret: "hi-wine-jwt-secret-key"  # JWT密钥
  expire: 86400                     # 过期时间，单位：秒（24小时）
```

### 微信配置

```yaml
# 微信配置
wechat:
  appid: "your_appid"
  secret: "your_secret"
```

## 部署方法

### 常规部署

1. 编译项目
   ```bash
   go build -o hi-wine
   ```

2. 运行服务
   ```bash
   ./hi-wine
   ```

### Docker 部署

1. 构建镜像
   ```bash
   docker build -t hi-wine:latest .
   ```

2. 使用 Docker Compose 启动服务
   ```bash
   docker-compose up -d
   ```

## API 文档

启动服务后，可通过以下地址访问 API 文档：

- Swagger UI: `http://localhost:8000/swagger`
- OpenAPI JSON: `http://localhost:8000/api.json`

## 主要接口

### 公开接口

- 商品查询：`/api/public/*`
- 套餐查询：`/api/public/*`

### 认证接口

- 微信登录：`/api/v1/auth/wechat-login`
- 管理员登录：`/api/v1/admin/login`

### 用户接口

- 用户信息：`/api/v1/user/profile`
- 用户订单：`/api/v1/user/orders`
- 用户套餐：`/api/v1/user/packages`

### 管理后台接口

- 商品管理：`/api/v1/admin/products/*`
- 分类管理：`/api/v1/admin/categories/*`
- 套餐管理：`/api/v1/admin/packages/*`
- 订单管理：`/api/v1/admin/orders/*`
- 文件上传：`/api/v1/admin/upload/*`

## 注意事项

- 确保 MinIO 服务正常运行
- 确保配置的访问密钥和密钥正确
- 确保配置的存储桶存在或有权限创建
- 微信支付回调需要配置可公网访问的域名

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

[MIT License](LICENSE)

## 联系方式

- 项目负责人：HuangEnJun
- 邮箱：1021217094@qq.com
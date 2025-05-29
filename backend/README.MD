# GoFrame Template For SingleRepo

Quick Start: 
- https://goframe.org/quick

# Hi-Wine 后端服务

基于GoFrame框架的酒水销售后台管理系统

## MinIO 集成说明

本项目使用 MinIO 作为对象存储服务，用于存储商品图片等文件。

### 配置说明

在 `manifest/config/config.yaml` 文件中配置 MinIO 相关参数：

```yaml
# MinIO配置
minio:
  endpoint: "127.0.0.1:9000"  # MinIO服务地址
  accessKey: "minioadmin"     # 访问密钥
  secretKey: "minioadmin"     # 密钥
  useSSL: false               # 是否使用SSL
  bucket: "hi-wine"           # 存储桶名称
  region: ""                  # 区域，可为空
  domain: "http://127.0.0.1:9000"  # 访问域名
```

### 使用方法

1. 安装并启动 MinIO 服务
2. 配置 MinIO 参数
3. 使用管理后台上传图片
4. 图片将自动存储到 MinIO 中并返回可访问的 URL

### 接口说明

- 图片上传接口: `/admin/upload/image`
- 商品管理接口: `/admin/product/*`

### 注意事项

- 确保 MinIO 服务正常运行
- 确保配置的访问密钥和密钥正确
- 确保配置的存储桶存在或有权限创建
# 酒馆小程序故障排除指南

本文档提供了常见问题的排查和解决方法，以帮助您快速解决在使用酒馆小程序过程中可能遇到的问题。

## 数据不显示问题

### 问题：后台接口正常返回数据，但小程序页面没有显示

可能的原因：
1. **响应格式不匹配**：后台返回的数据格式与前端期望的不一致
2. **数据库连接问题**：数据库连接配置错误或数据库表结构有问题
3. **权限问题**：API请求需要授权但未提供有效的token

解决方法：

#### 1. 检查响应格式

小程序前端期望API响应格式为：
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "list": [...] // 或其他数据
  }
}
```

如果后端直接返回原始数据而不是这种包装格式，需要：
- 修改后端中间件，统一包装响应
- 或者修改前端代码，适配后端返回的格式

可以使用提供的测试脚本检查API响应格式：
```bash
chmod +x ./test_api.sh
./test_api.sh
```

#### 2. 使用调试功能

在小程序页面添加调试信息显示：
```vue
<view class="debug-info">
  <text>接口响应: {{ debugInfo }}</text>
</view>
```

在JS代码中记录API响应：
```javascript
console.log('API响应:', JSON.stringify(res));
this.debugInfo = JSON.stringify(res);
```

#### 3. 检查数据库

验证数据库中是否有数据：
```sql
SELECT * FROM category LIMIT 10;
SELECT * FROM product LIMIT 10;
```

如果没有数据，可以导入示例数据：
```bash
./import_sample_data.sh
```

## 登录和授权问题

### 问题：无法登录或提示"请先登录"

可能的原因：
1. **Token失效**：JWT令牌已过期
2. **小程序配置问题**：小程序未正确配置微信登录
3. **后端JWT验证问题**：后端JWT密钥或验证逻辑错误

解决方法：

#### 1. 清除本地缓存

在小程序开发工具中：
- 点击"清缓存"
- 重新登录

#### 2. 检查JWT配置

确保后端JWT密钥配置正确：
```go
// backend/internal/middleware/jwt.go
var JwtSecret = []byte("hi-wine-jwt-secret") // 应与生成token的密钥一致
```

#### 3. 检查网络请求

在微信开发者工具的Network面板中观察请求：
- 请求URL是否正确
- 请求头中是否包含Authorization
- 响应状态码和内容是否正确

## 开发环境配置问题

### 问题：开发环境无法访问后端API

可能的原因：
1. **跨域问题**：后端未配置CORS
2. **网络连接问题**：IP地址或端口配置错误
3. **后端服务未启动**：后端服务未正常运行

解决方法：

#### 1. 检查配置文件

确保小程序中的API地址配置正确：
```javascript
// config.js
export const BASE_URL = 'http://172.16.10.28:8000'; // 修改为实际IP和端口
```

#### 2. 启动后端服务

确保后端服务正常运行：
```bash
cd backend
go run main.go
```

#### 3. 配置CORS（如需要）

如果是在不同域名下开发，后端需要配置CORS支持：
```go
// 在backend/internal/cmd/cmd.go中添加
s.Use(func(r *ghttp.Request) {
    r.Response.CORSDefault()
    r.Middleware.Next()
})
```

## 其他常见问题

### 问题：添加商品到购物车失败

检查本地存储操作：
```javascript
// 确保正确获取现有购物车数据
let cart = uni.getStorageSync('cart') || [];
```

### 问题：订单创建失败

检查订单请求格式是否正确：
```javascript
const items = cart.map(item => ({
    productId: item.id,
    quantity: item.count
}));
```

如果您遇到的问题在本文档中未涵盖，请提交issue或联系技术支持。 
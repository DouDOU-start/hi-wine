本API文档旨在详细阐述后端提供的 RESTful API 接口，供小程序客户端和管理后台调用。
接口通用说明
- 请求方式: 遵循 RESTful 规范，例如 GET、POST、PUT、DELETE。
- 认证: 大部分接口需要认证，通常通过 JWT (JSON Web Token) 或 Session 进行。管理后台接口将实施更严格的权限控制。
- 请求头: 针对 POST 和 PUT 请求，请务必设置 Content-Type: application/json。
- 响应格式: 所有响应均采用 application/json 格式。
- 状态码: 使用标准的 HTTP 状态码表示请求结果： 
  - 200 OK: 请求成功。
  - 201 Created: 资源已成功创建 (针对 POST 请求)。
  - 204 No Content: 请求成功，但无内容返回 (例如，成功的 DELETE 操作)。
  - 400 Bad Request: 请求参数无效或请求格式错误。
  - 401 Unauthorized: 需要认证或认证凭据无效。
  - 403 Forbidden: 已认证，但无权访问该资源。
  - 404 Not Found: 请求的资源不存在。
  - 429 Too Many Requests: 请求过于频繁，触发限流。
  - 500 Internal Server Error: 服务器发生未知错误。
- 错误信息: 错误响应体中将一致包含 code (整数) 和 message (字符串) 字段，提供详细的错误描述。 
  - 示例错误响应:
  - JSON
{
    "code": 40001,
    "message": "微信登录凭证无效。"
}
  - 注：更细致的错误 code 有助于前端进行特定错误场景的处理。

---
一、用户相关接口 (小程序 & 管理后台)
1. 用户登录/注册 (微信授权)
- 接口: POST /api/v1/auth/wechat-login
- 描述: 小程序用户通过微信授权进行登录或注册。
- 请求参数: 
- JSON
{
    "code": "your_wechat_login_code" // (String, 必填) 微信登录凭证 code。
}
- 响应: 
- JSON
{
    "code": 200,
    "message": "success",
    "data": {
        "token": "your_jwt_token", // (String) 用户身份认证 Token。
        "user": {
            "id": 1,
            "openid": "xxxxxxxx",
            "nickname": "张三",
            "avatar_url": "http://...",
            "phone": "13800138000"
         }
    }
}
2. 获取用户个人信息
- 接口: GET /api/v1/user/profile
- 描述: 获取当前登录用户的个人信息。
- 认证: 需要用户Token。
- 请求参数: 无。
- 响应: 用户信息对象。
3. 更新用户个人信息
- 接口: PUT /api/v1/user/profile
- 描述: 更新当前登录用户的个人信息 (例如手机号、昵称、头像)。
- 认证: 需要用户Token。
- 请求参数: 
- JSON
{
    "phone": "13912345678", // (String, 可选) 用户手机号。
    "nickname": "新的昵称", // (String, 可选) 用户昵称。
    "avatar_url": "http://new-avatar-url.com/image.jpg" // (String, 可选) 用户头像URL。
}
- 响应: 更新后的用户信息对象。

---
二、商品相关接口 (小程序 & 管理后台)
1. 获取所有商品分类
- 接口: GET /api/v1/categories
- 描述: 获取所有已激活的商品分类列表。
- 请求参数: 无。
- 响应: 分类对象数组，包含 id、name、sort_order、image_url (可选) 等字段。
2. 获取某分类下的商品列表
- 接口: GET /api/v1/categories/{category_id}/products
- 描述: 获取指定分类下的所有已激活商品。
- 请求参数: 
  - category_id: (Integer, 必填) 分类ID。
- 查询参数: 
  - page: (Integer, 可选, 默认: 1) 当前页码。
  - limit: (Integer, 可选, 默认: 10) 每页商品数量。
- 响应: 商品对象数组，包含 id、name、price、image_url、stock、description (可选) 等字段。
3. 获取商品详情
- 接口: GET /api/v1/products/{product_id}
- 描述: 获取单个商品的详细信息。
- 请求参数: 
  - product_id: (Integer, 必填) 商品ID。
- 响应: 商品详细信息对象。

---
三、订单相关接口 (小程序 & 管理后台)
1. 创建订单
- 接口: POST /api/v1/orders
- 描述: 创建一个新的订单。 
  - 畅饮套餐逻辑: 
    1. 判断有效畅饮套餐: 当用户下单时，后端需首先判断该用户是否有当前有效的畅饮套餐 (即 user_packages 表中 status 为 active 且 end_time 未到的记录)。
    2. 判断是否为套餐内商品: 对于 items 列表中的每个 product_id，检查它是否包含在用户当前有效套餐的 included_products 列表中。
    3. 价格处理: 
      - 如果商品在有效套餐内，则该 order_item 的 item_price 将设置为 0，并标记 is_package_item 为 TRUE，同时关联 user_package_id 到用户的有效套餐记录。
      - 如果商品不在有效套餐内，则按正常价格计算 item_price。
    4. 库存扣减: 无论商品是否在套餐内，下单时都需扣减其对应的库存。
- 认证: 需要用户Token。
- 请求参数: 
- JSON
{
    "table_qrcode_id": 101, // (Integer, 必填) 桌号二维码ID。
    "items": [ // (Array, 必填) 订单商品列表。
        {
            "product_id": 1, // (Integer, 必填) 商品ID。
            "quantity": 2, // (Integer, 必填) 购买数量。
            "notes": "加冰" // (String, 可选) 单个商品备注。
        },
        {
            "product_id": 5,"quantity": 1
        }
    ],
    "total_notes": "整体订单备注：不要辣" // (String, 可选) 订单的整体备注。
}
- 响应: 
- JSON
{
    "code": 200,
    "message": "订单创建成功。",
    "data": {
        "order_id": 1001, // (Integer) 新创建的订单ID。
        "order_sn": "20250603123456789", // (String) 订单号。
        "total_amount": 58.00, // (Decimal) 订单总金额 (扣除套餐优惠后的实际支付金额)。
        "prepay_id": "wx1234567890abcdef" // (String) 微信支付预支付ID (用于小程序调起支付)。
    }
}
2. 微信支付回调 (异步通知)
- 接口: POST /api/v1/wechat-pay/notify
- 描述: 微信支付成功后的异步回调通知。后台接收到此通知后，将更新订单的支付状态。此接口由微信支付服务器调用。
- 请求参数: 微信支付提供的 XML 数据。
- 响应: 微信支付要求的 XML 响应。
- 注意: 此接口不对外暴露，仅供微信支付服务器调用。实现时需严格进行签名验证和幂等性处理，以确保数据安全和一致性。
3. 获取用户订单列表
- 接口: GET /api/v1/user/orders
- 描述: 获取当前登录用户的所有订单列表。
- 认证: 需要用户Token。
- 查询参数: 
  - status: (String, 可选) 筛选订单状态 (new, processing, completed, cancelled, paid, pending_payment, refunded)。
  - page: (Integer, 可选, 默认: 1) 页码。
  - limit: (Integer, 可选, 默认: 10) 每页数量。
- 响应: 订单列表 (数组)。
4. 获取订单详情
- 接口: GET /api/v1/orders/{order_id}
- 描述: 获取指定订单的详细信息，包含订单项。
- 认证: 需要用户Token (或管理后台Token)。
- 请求参数: 
  - order_id: (Integer, 必填) 订单ID。
- 响应: 订单详细信息及订单项列表。

---
四、管理后台接口 (酒馆内部管理)
1. 管理员登录
- 接口: POST /api/v1/admin/login
- 描述: 后台管理员登录接口。
- 请求参数: 
- JSON
{
    "username": "admin_user", // (String, 必填) 管理员用户名。
    "password": "admin_password" // (String, 必填) 管理员密码。
}
- 响应: 
- JSON
{
    "code": 200,
    "message": "登录成功。",
    "data": {
        "token": "your_admin_jwt_token", // (String) 管理员身份认证 Token。
        "admin_user": {
            "id": 1,
            "username": "admin_user",
            "role": "店长" // (String) 管理员角色 (例如：'店长', '服务员', '厨师')，用于权限控制。
        }
    }
}
2. 获取所有订单 (管理后台)
- 接口: GET /api/v1/admin/orders
- 描述: 管理员获取所有订单，支持筛选和分页。
- 认证: 需要管理员Token。
- 查询参数: 
  - status: (String, 可选) 筛选订单状态。
  - page: (Integer, 可选) 页码。
  - limit: (Integer, 可选) 每页数量。
  - order_sn: (String, 可选) 按订单号搜索。
  - start_date: (Date, 可选, 格式: YYYY-MM-DD) 订单创建起始日期。
  - end_date: (Date, 可选, 格式: YYYY-MM-DD) 订单创建结束日期。
  - table_number: (String, 可选) 按桌号搜索。
  - user_id: (Integer, 可选) 按用户ID搜索。
- 响应: 订单列表 (数组)。
3. 更新订单状态 (管理后台)
- 接口: PUT /api/v1/admin/orders/{order_id}/status
- 描述: 管理员更新订单状态 (例如：标记为 processing 处理中, completed 已完成, cancelled 已取消)。
- 认证: 需要管理员Token，并根据管理员角色进行权限校验。
- 请求参数: 
- JSON
{
    "status": "processing", // (String, 必填) 新的订单状态。
    "reason": "顾客要求取消" // (String, 可选) 状态变更原因 (例如取消订单时)。
}
- 响应: 更新后的订单信息。
4. 商品管理 (增删改查)
- 认证: 需要管理员Token，并根据管理员角色进行权限校验。
- 创建商品: 
  - 接口: POST /api/v1/admin/products
  - 描述: 新增商品。
  - 请求参数: name (String, 必填), category_id (Integer, 必填), price (Decimal, 必填), stock (Integer, 必填), image_url (String, 可选), description (String, 可选), is_active (Boolean, 可选, 默认true)。
  - 响应: 新创建的商品信息。
- 更新商品: 
  - 接口: PUT /api/v1/admin/products/{product_id}
  - 描述: 更新指定商品信息。
  - 请求参数: (部分或全部商品字段，参考创建商品请求参数)。
  - 响应: 更新后的商品信息。
- 删除商品: 
  - 接口: DELETE /api/v1/admin/products/{product_id}
  - 描述: 删除指定商品。
  - 响应: 204 No Content 或成功消息。
- 获取商品列表: 
  - 接口: GET /api/v1/admin/products
  - 描述: 获取商品列表，支持分页、搜索、按分类筛选。
  - 查询参数: page, limit, name (String, 模糊搜索), category_id (Integer, 精确筛选), is_active (Boolean, 筛选激活状态)。
  - 响应: 商品列表 (数组)。
5. 分类管理 (增删改查)
- 认证: 需要管理员Token，并根据管理员角色进行权限校验。
- 创建分类: POST /api/v1/admin/categories
  - 请求参数: name (String, 必填), sort_order (Integer, 可选), image_url (String, 可选)。
- 更新分类: PUT /api/v1/admin/categories/{category_id}
  - 请求参数: (部分或全部分类字段)。
- 删除分类: DELETE /api/v1/admin/categories/{category_id}
- 获取分类列表: GET /api/v1/admin/categories
6. 畅饮套餐管理 (增删改查)
- 认证: 需要管理员Token，并根据管理员角色进行权限校验。
- 创建套餐: POST /api/v1/admin/packages
  - 请求参数: name (String, 必填), price (Decimal, 必填), duration_hours (Integer, 必填, 有效时长), description (String, 可选), is_active (Boolean, 可选)。
- 更新套餐: PUT /api/v1/admin/packages/{package_id}
- 删除套餐: DELETE /api/v1/admin/packages/{package_id}
- 获取套餐列表: GET /api/v1/admin/packages
  - 查询参数: 支持 page, limit, name 模糊搜索。
7. 管理套餐包含酒水
- 认证: 需要管理员Token，并根据管理员角色进行权限校验。
- 添加套餐酒水: 
  - 接口: POST /api/v1/admin/packages/{package_id}/products
  - 描述: 为指定套餐添加可畅饮的酒水。
  - 请求参数: 
  - JSON
{
    "product_ids": [1, 5, 8] // (Array of Integer, 必填) 要添加到套餐的商品ID列表。
}
- 移除套餐酒水: 
  - 接口: DELETE /api/v1/admin/packages/{package_id}/products/{product_id}
  - 描述: 从指定套餐中移除某个酒水。
8. 用户畅饮套餐记录查询
- 接口: GET /api/v1/admin/user-packages
- 描述: 管理员查询所有用户的畅饮套餐购买和使用记录。
- 认证: 需要管理员Token。
- 查询参数: 支持按 user_id (Integer), package_id (Integer), status (String, 例如 active, expired, pending), start_date (Date), end_date (Date), page, limit 等筛选。
- 响应: UserPackages 列表 (数组)。

---
五、打印服务接口 (管理后台内部调用或独立服务)
1. 打印订单
- 接口: POST /api/v1/print/order (建议作为内部服务或独立的打印服务接口)
- 描述: 接收订单数据，连接打印机并打印订单小票。
- 请求参数: 
- JSON
{
    "order_id": 1001, // (Integer, 必填) 需要打印的订单ID。
    "print_type": "kitchen_ticket", // (String, 可选, 默认: "customer_bill") 打印类型 (例如：`kitchen_ticket` 厨房单, `customer_bill` 客户账单)。
    "notes": "加急打印" // (String, 可选) 打印备注。
}
- 认证: 内部调用或需要管理员Token (视具体实现方式)。
- 响应: 
- JSON
{
    "code": 200,
    "message": "打印任务已发送。",
    "data": {
        "status": "success", // (String) 打印状态 (`success`, `failed`, `pending`)。
        "print_job_id": "uuid_12345" // (String, 可选) 打印任务ID，用于后续查询打印状态。
    }
}
- 实现说明: 
  - 此接口通常由后台系统内部调用，例如当新订单生成或订单状态更新时触发。
  - 实际的打印逻辑 (连接打印机、生成打印内容、发送打印指令) 建议封装在独立的打印服务模块中，以提高可扩展性和解耦。
  - 打印内容应包含：订单号、桌号、商品名称、数量、单价、总价、下单时间、备注等关键信息。
  - 考虑打印队列和重试机制，确保打印任务的可靠性。

---
六、二维码管理接口 (管理后台)
1. 生成桌号二维码
- 接口: POST /api/v1/admin/table-qrcodes
- 描述: 为指定桌号生成小程序点餐二维码。二维码通常会包含一个指向小程序页面的链接，并带上 table_qrcode_id 参数。
- 认证: 需要管理员Token。
- 请求参数: 
- JSON
{
    "table_number": "A01" // (String, 必填) 桌号。
}
- 响应: 
- JSON
{
    "code": 200,"message": "二维码生成成功。",
    "data": {"id": 1, // (Integer) 新生成的二维码ID。
        "table_number": "A01", // (String) 对应的桌号。
        "qrcode_url": "http://your_cdn.com/qrcodes/A01_1.png", // (String) 生成的二维码图片URL。
        "created_at": "2025-06-03T10:00:00Z"
    }
}
2. 获取桌号二维码列表
- 接口: GET /api/v1/admin/table-qrcodes
- 描述: 获取所有桌号二维码列表。
- 认证: 需要管理员Token。
- 查询参数: 
  - table_number: (String, 可选) 按桌号模糊搜索。
  - page: (Integer, 可选) 页码。
  - limit: (Integer, 可选) 每页数量。
- 响应: 二维码列表 (数组)。

---
完善建议总结与扩展
1. 权限管理 (RBAC):
  - 在管理后台接口层面，需要实现基于角色的访问控制 (RBAC)。定义不同的角色（例如：店长、服务员、厨师、财务），并为每个角色配置其可访问的接口和操作权限（增、删、改、查）。
  - 在接口认证后，根据用户的角色来动态判断其是否有权执行当前操作。
2. 日志记录:
  - 所有关键操作（如用户登录、订单创建/更新、商品/套餐增删改、管理员操作等）都应有详细的日志记录。
  - 日志应包含：请求时间、请求IP、操作用户ID、操作类型、请求参数、响应结果、耗时等信息。这对于问题排查、审计和安全监控至关重要。
3. 异常处理:
  - 建立统一的全局异常处理机制。针对不同类型的错误（例如：参数校验失败、数据库操作异常、业务逻辑错误、权限不足等），返回清晰、友好的错误信息和对应的 HTTP 状态码。
  - 使用自定义错误码 (如 40001、50002 等) 来区分具体错误类型，方便前端进行精准处理和用户提示。
4. 数据校验:
  - 所有接收到的请求参数都必须进行严格的数据类型、格式、长度、范围等校验。
  - 可以使用框架提供的校验工具或手动实现校验逻辑，确保数据的合法性，防止非法输入和安全漏洞。
5. 前端交互:
  - 小程序端: 需配套开发购物车、订单确认、支付流程、个人中心、我的订单、套餐购买等页面。
  - 管理后台: 需开发订单列表、商品管理、分类管理、套餐管理、用户管理、报表统计、桌号二维码管理等功能界面。
6. 营销功能:
  - 优惠券: 增加优惠券创建、发放、领取、使用的接口和逻辑。
  - 满减/满赠: 实现基于订单金额或商品数量的满减/满赠活动。
  - 会员积分: 建立积分体系，用户消费获取积分，积分可兑换商品或抵扣金额。
7. 消息通知:
  - 新订单通知: 当有新订单生成时，除了打印，还可通过微信小程序订阅消息、短信、后台消息弹窗等方式即时通知服务员或厨房。
  - 订单状态变更通知: 当订单状态更新（如已支付、已完成）时，可通知用户。
8. 库存管理:
  - 实时库存扣减: 下单时严格扣减库存。
  - 库存预警: 设置库存阈值，当商品库存低于阈值时触发预警通知管理员。
  - 库存回滚: 在订单支付失败或取消时，应将库存回滚。
9. 报表统计:
  - 后台提供丰富的统计报表接口，如： 
    - 销售额统计: 按日/周/月/年统计总销售额、不同商品分类的销售额。
    - 商品销量榜: 统计畅销商品。
    - 订单量统计: 统计不同状态的订单数量。
    - 畅饮套餐使用情况: 统计套餐销售量、使用频率、套餐内商品消耗量等。
    - 用户消费行为分析: 高价值用户、消费频率等。
10. 安全增强:
  - HTTPS: 所有接口强制使用 HTTPS，确保数据传输加密。
  - API 限流: 防止恶意请求和DDoS攻击。
  - 输入验证: 再次强调对所有用户输入进行严格验证，防止SQL注入、XSS攻击等。
  - 敏感数据加密: 数据库中存储的敏感信息（如用户手机号、部分支付信息）应进行加密。
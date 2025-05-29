# 酒水示例数据导入指南

本目录包含了酒馆小程序的示例数据，包括各种类型的酒水分类和商品信息。这些数据可以帮助您快速测试和展示小程序的功能。

## 数据内容

示例数据包括以下内容：

- 7个酒水分类：红葡萄酒、白葡萄酒、啤酒、威士忌、鸡尾酒、清酒和白酒
- 每个分类3个示例商品，共21个酒水商品
- 每个商品包含名称、价格、库存、图片URL和详细描述

## 导入方法

### 方法一：使用脚本导入（推荐）

如果您在Linux或macOS系统上，可以使用提供的Shell脚本自动导入数据：

1. 确保您已安装MySQL客户端
2. 进入resource目录
3. 执行以下命令赋予脚本执行权限：
   ```bash
   chmod +x import_sample_data.sh
   ```
4. 运行脚本：
   ```bash
   ./import_sample_data.sh
   ```

### 方法二：手动导入

1. 使用MySQL命令行客户端连接到您的数据库：
   ```bash
   mysql -u用户名 -p密码 数据库名
   ```

2. 在MySQL提示符下，执行以下命令导入数据：
   ```sql
   source /path/to/sample_data.sql
   ```
   请将`/path/to/sample_data.sql`替换为实际的文件路径。

3. 或者，您可以直接复制`sample_data.sql`文件中的SQL语句并在MySQL客户端中执行。

## 验证导入

导入完成后，您可以执行以下SQL查询验证数据是否导入成功：

```sql
-- 查询分类数量
SELECT COUNT(*) FROM category;

-- 查询商品数量
SELECT COUNT(*) FROM product;

-- 按分类查询商品
SELECT c.name AS category_name, COUNT(p.id) AS product_count
FROM category c
LEFT JOIN product p ON c.id = p.category_id
GROUP BY c.id;
```

## 注意事项

- 导入数据前请确保您的数据库中已创建了相应的表结构
- 如果数据库中已有数据，导入操作可能会导致主键冲突
- 导入前建议备份您的数据库 
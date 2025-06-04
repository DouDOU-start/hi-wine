-- 创建数据库（如果尚未创建）
-- 根据您的数据库管理系统和配置，这部分可能需要手动执行或调整。
-- 例如，对于 MySQL:
-- CREATE DATABASE IF NOT EXISTS `pub_ordering_system` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- USE `pub_ordering_system`;

-- 1. 用户表 (Users)
CREATE TABLE `users` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    `openid` VARCHAR(50) NOT NULL UNIQUE COMMENT '微信openid',
    `nickname` VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '微信昵称',
    `avatar_url` VARCHAR(255) NULL COMMENT '微信头像URL',
    `phone` VARCHAR(20) NULL COMMENT '手机号',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT='用户表';

-- 2. 商品分类表 (Categories)
CREATE TABLE `categories` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL UNIQUE  COMMENT '分类名称',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序，数字越小越靠前',
    `is_active` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否激活（是否显示）',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT='商品分类表';

-- 3. 商品表 (Products)
CREATE TABLE `products` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '商品ID',
    `category_id` INT NOT NULL COMMENT '所属分类ID',
    `name` VARCHAR(100) NOT NULL UNIQUE COMMENT '商品名称',
    `description` TEXT NULL COMMENT '商品描述',
    `price` DECIMAL(10, 2) NOT NULL COMMENT '商品价格',
    `image_url` VARCHAR(255) NULL COMMENT '商品图片URL',
    `stock` INT NOT NULL DEFAULT 0 COMMENT '库存数量',
    `is_active` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否上架',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    CONSTRAINT `fk_products_category_id` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) COMMENT='商品表';

-- 4. 桌号二维码表 (TableQRCodes)
CREATE TABLE `table_qrcodes` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    `table_number` VARCHAR(20) NOT NULL UNIQUE COMMENT '桌号，例如"A1", "吧台"',
    `qrcode_url` VARCHAR(255) NULL COMMENT '生成的二维码图片URL',
    `status` ENUM('idle', 'occupied') NOT NULL DEFAULT 'idle' COMMENT '桌位状态',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT='桌号二维码表';

-- 5. 畅饮套餐表 (DrinkAllYouCanPackages) - 新增
CREATE TABLE `drink_all_you_can_packages` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '套餐ID',
    `name` VARCHAR(100) NOT NULL COMMENT '套餐名称',
    `description` TEXT NULL COMMENT '套餐描述',
    `price` DECIMAL(10, 2) NOT NULL COMMENT '套餐价格',
    `duration_minutes` INT NOT NULL DEFAULT 0 COMMENT '有效时长（分钟），0表示无时间限制',
    `is_active` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否激活（是否可售）',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT='畅饮套餐表';

-- 6. 套餐包含商品关联表 (PackageProducts) - 新增
CREATE TABLE `package_products` (
    `package_id` INT NOT NULL COMMENT '关联畅饮套餐ID',
    `product_id` INT NOT NULL COMMENT '关联商品ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`package_id`, `product_id`),
    CONSTRAINT `fk_package_products_package_id` FOREIGN KEY (`package_id`) REFERENCES `drink_all_you_can_packages` (`id`),
    CONSTRAINT `fk_package_products_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) COMMENT='套餐包含商品关联表';

-- 7. 订单表 (Orders)
CREATE TABLE `orders` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '订单ID',
    `order_sn` VARCHAR(32) NOT NULL UNIQUE COMMENT '订单号，唯一',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `table_qrcode_id` INT NULL COMMENT '关联的桌号二维码ID',
    `total_amount` DECIMAL(10, 2) NOT NULL COMMENT '订单总金额',
    `payment_status` ENUM('pending', 'paid', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT '支付状态',
    `order_status` ENUM('new', 'processing', 'completed', 'cancelled') NOT NULL DEFAULT 'new' COMMENT '订单状态',
    `payment_method` VARCHAR(20) NULL COMMENT '支付方式（例如：wechat_pay）',
    `transaction_id` VARCHAR(50) NULL COMMENT '微信支付交易ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `paid_at` DATETIME NULL COMMENT '支付时间',
    CONSTRAINT `fk_orders_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_orders_table_qrcode_id` FOREIGN KEY (`table_qrcode_id`) REFERENCES `table_qrcodes` (`id`)
) COMMENT='订单表';

-- 8. 用户套餐购买记录表 (UserPackages) - 新增
CREATE TABLE `user_packages` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
    `user_id` INT NOT NULL COMMENT '关联用户ID',
    `package_id` INT NOT NULL COMMENT '关联畅饮套餐ID',
    `order_id` INT NOT NULL COMMENT '关联购买此套餐的订单ID',
    `start_time` DATETIME NULL COMMENT '套餐开始时间（首次使用时激活）',
    `end_time` DATETIME NULL COMMENT '套餐结束时间（根据duration_minutes计算）',
    `status` ENUM('active', 'expired', 'refunded', 'pending') NOT NULL DEFAULT 'pending' COMMENT '套餐状态',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    CONSTRAINT `fk_user_packages_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_user_packages_package_id` FOREIGN KEY (`package_id`) REFERENCES `drink_all_you_can_packages` (`id`),
    CONSTRAINT `fk_user_packages_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) COMMENT='用户套餐购买记录表';

-- 9. 订单详情表 (OrderItems) - 修改
CREATE TABLE `order_items` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '订单项ID',
    `order_id` INT NOT NULL COMMENT '订单ID',
    `product_id` INT NOT NULL COMMENT '商品ID',
    `product_name` VARCHAR(100) NOT NULL COMMENT '冗余商品名称',
    `price` DECIMAL(10, 2) NOT NULL COMMENT '下单时商品单价',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '购买数量',
    `subtotal` DECIMAL(10, 2) NOT NULL COMMENT '小计',
    `is_package_item` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否为畅饮套餐内商品',
    `user_package_id` INT NULL COMMENT '关联的用户套餐购买记录ID（如果为套餐商品）',
    `item_price` DECIMAL(10, 2) NOT NULL COMMENT '该订单项的实际结算价格（畅饮则为0）',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    CONSTRAINT `fk_order_items_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
    CONSTRAINT `fk_order_items_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
    CONSTRAINT `fk_order_items_user_package_id` FOREIGN KEY (`user_package_id`) REFERENCES `user_packages` (`id`)
) COMMENT='订单详情表';
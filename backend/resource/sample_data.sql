-- 插入酒水分类
INSERT INTO `category` (`name`, `sort`) VALUES
('红葡萄酒', 1),
('白葡萄酒', 2),
('啤酒', 3),
('威士忌', 4),
('鸡尾酒', 5),
('清酒', 6),
('白酒', 7);

-- 插入红葡萄酒
INSERT INTO `product` (`name`, `category_id`, `price`, `stock`, `image`, `status`, `description`) VALUES
('法国波尔多干红', 1, 198.00, 50, 'https://img.alicdn.com/imgextra/i2/3081122511/O1CN01IgGhD31ybQEwkO4hE_!!3081122511.jpg', 1, '来自法国波尔多产区的经典干红葡萄酒，口感丰富，单宁适中，带有黑醋栗、黑莓和淡淡的橡木香气。'),
('意大利基安蒂经典红葡萄酒', 1, 168.00, 30, 'https://img.alicdn.com/imgextra/i3/2206997287603/O1CN01XhOCXr25ULV0XVfkm_!!2206997287603.jpg', 1, '采用桑娇维塞葡萄酿制，带有樱桃、李子和香料的风味，酸度适中，单宁柔和。'),
('澳大利亚西拉子红葡萄酒', 1, 258.00, 25, 'https://img.alicdn.com/imgextra/i3/2212384246732/O1CN01gp75S41uQToxzg8U9_!!2212384246732.jpg', 1, '来自澳大利亚的浓郁西拉子红葡萄酒，带有黑莓、黑胡椒和巧克力的风味，单宁强劲，余味悠长。');

-- 插入白葡萄酒
INSERT INTO `product` (`name`, `category_id`, `price`, `stock`, `image`, `status`, `description`) VALUES
('法国霞多丽干白', 2, 178.00, 40, 'https://img.alicdn.com/imgextra/i1/2206997287603/O1CN01PYKHZ525ULUwzrQ1g_!!2206997287603.jpg', 1, '法国经典霞多丽干白葡萄酒，带有青苹果、柑橘和矿物质的香气，酸度清爽，口感圆润。'),
('德国雷司令半甜白', 2, 168.00, 35, 'https://img.alicdn.com/imgextra/i2/90699028/O1CN01FG5v4w1zq83XLqiYq_!!90699028.jpg', 1, '来自德国的半甜型雷司令，带有蜜桃、杏子和柑橘的风味，酸甜平衡，余味悠长。'),
('新西兰长相思干白', 2, 188.00, 30, 'https://img.alicdn.com/imgextra/i1/2208023209815/O1CN01mLdPOl1aGkNqN15J1_!!2208023209815.jpg', 1, '新西兰马尔堡产区的长相思白葡萄酒，带有浓郁的青草、百香果和青柠檬的香气，酸度清新，口感活泼。');

-- 插入啤酒
INSERT INTO `product` (`name`, `category_id`, `price`, `stock`, `image`, `status`, `description`) VALUES
('比利时修道院啤酒', 3, 68.00, 100, 'https://img.alicdn.com/imgextra/i3/2208023209815/O1CN01vl7gZK1aGkNtSNsBF_!!2208023209815.jpg', 1, '比利时传统修道院酿造的浓郁啤酒，带有焦糖、香料和果味，酒体饱满，余味悠长。'),
('德国小麦啤酒', 3, 48.00, 120, 'https://img.alicdn.com/imgextra/i1/2212384246732/O1CN0188Tphp1uQTp1Dapqo_!!2212384246732.jpg', 1, '德国传统小麦啤酒，色泽金黄，带有香蕉、丁香和柑橘的香气，口感顺滑清爽。'),
('英国印度淡色艾尔', 3, 58.00, 80, 'https://img.alicdn.com/imgextra/i3/2209301165602/O1CN019XiILK1oU9oZTsazJ_!!2209301165602.jpg', 1, 'IPA风格啤酒，具有浓郁的花香和柑橘风味，苦味适中，回味持久。');

-- 插入威士忌
INSERT INTO `product` (`name`, `category_id`, `price`, `stock`, `image`, `status`, `description`) VALUES
('苏格兰单一麦芽威士忌12年', 4, 588.00, 20, 'https://img.alicdn.com/imgextra/i1/3470203725/O1CN01MYVYSt1fEOv1qDT3n_!!3470203725.jpg', 1, '经典苏格兰单一麦芽威士忌，陈年12年，带有麦芽、蜂蜜和淡淡的烟熏香气，口感圆润，余味悠长。'),
('日本调和威士忌', 4, 498.00, 25, 'https://img.alicdn.com/imgextra/i3/2212384246732/O1CN01lB0XY11uQToxyvd5w_!!2212384246732.jpg', 1, '精致的日本调和威士忌，带有蜂蜜、水果和轻微的烟熏味，口感柔顺，平衡感极佳。'),
('美国波本威士忌', 4, 398.00, 30, 'https://img.alicdn.com/imgextra/i4/2212384246732/O1CN01bJJP0p1uQTp1FDZhh_!!2212384246732.jpg', 1, '传统美国波本威士忌，带有焦糖、香草和橡木的风味，口感浓郁，余味带有淡淡的甜味。');

-- 插入鸡尾酒
INSERT INTO `product` (`name`, `category_id`, `price`, `stock`, `image`, `status`, `description`) VALUES
('经典马提尼', 5, 68.00, 50, 'https://img.alicdn.com/imgextra/i2/2212343599606/O1CN01q0JHFG1LxSDzj0tqB_!!2212343599606.jpg', 1, '用优质金酒和干味美思调制的经典鸡尾酒，干爽清冽，带有淡淡的草本香气。'),
('莫吉托', 5, 58.00, 60, 'https://img.alicdn.com/imgextra/i4/3469659703/O1CN01EoBN141NI7IM6XeWA_!!3469659703.jpg', 1, '清新的朗姆酒鸡尾酒，加入新鲜薄荷、青柠和苏打水，口感清爽，风味层次丰富。'),
('血腥玛丽', 5, 58.00, 45, 'https://img.alicdn.com/imgextra/i1/3469659703/O1CN01WdYT461NI7IJp5PmK_!!3469659703.jpg', 1, '以伏特加和番茄汁为基础，加入柠檬汁、辣椒和香料，口感浓郁，略带辛辣。');

-- 插入清酒
INSERT INTO `product` (`name`, `category_id`, `price`, `stock`, `image`, `status`, `description`) VALUES
('大吟酿清酒', 6, 298.00, 25, 'https://img.alicdn.com/imgextra/i2/2212384246732/O1CN01TqCxvF1uQTown6Jtp_!!2212384246732.jpg', 1, '日本顶级大吟酿清酒，精米率50%以下，带有水果、花香和米香，口感精致，余味悠长。'),
('纯米清酒', 6, 198.00, 30, 'https://img.alicdn.com/imgextra/i3/2209301165602/O1CN01rQ5EBr1oU9ogTgQiD_!!2209301165602.jpg', 1, '不添加酿造酒精的纯米清酒，米香浓郁，口感醇厚，带有淡淡的甜味和果香。'),
('本酿造清酒', 6, 128.00, 40, 'https://img.alicdn.com/imgextra/i4/2211984646497/O1CN01hBgcGO29JiNuPKBTj_!!2211984646497.jpg', 1, '传统酿造的清酒，米香和酒香平衡，口感顺滑，适合常温或温饮。');

-- 插入白酒
INSERT INTO `product` (`name`, `category_id`, `price`, `stock`, `image`, `status`, `description`) VALUES
('茅台飞天53度', 7, 3688.00, 10, 'https://img.alicdn.com/imgextra/i1/2209301165602/O1CN01vwW2Zr1oU9ofPlkDr_!!2209301165602.jpg', 1, '中国传统名酒，酱香型白酒代表，带有酱香、窖香和坚果香，口感丰富，回味悠长。'),
('五粮液52度', 7, 1288.00, 15, 'https://img.alicdn.com/imgextra/i4/3470203725/O1CN01UNxDTu1fEOv0tJV6w_!!3470203725.jpg', 1, '浓香型白酒代表，带有浓郁的粮食香和水果香，口感甘甜醇厚，余味绵长。'),
('洋河蓝色经典梦之蓝', 7, 988.00, 20, 'https://img.alicdn.com/imgextra/i3/3469659703/O1CN019SWBuM1NI7ILOVrn5_!!3469659703.jpg', 1, '浓香型白酒，口感绵柔，带有淡雅的花香和果香，余味甘甜爽净。'); 
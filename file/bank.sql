CREATE TABLE `micro_bank` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk',
  `cn_name` varchar(64) NOT NULL COMMENT '用户中文姓名',
  `tel` char(11) NOT NULL,
  `latitude` decimal(10,4) NOT NULL DEFAULT '0.0000',
  `longtitude` decimal(10,4) NOT NULL DEFAULT '0.0000',
  `address` varchar(1024) DEFAULT '',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
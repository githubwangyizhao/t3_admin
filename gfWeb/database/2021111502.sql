CREATE TABLE `region`  (
                           `id` int(11) NOT NULL AUTO_INCREMENT,
                           `region` varchar(32) NOT NULL DEFAULT '台灣' COMMENT '国家/地区',
                           `currency` varchar(32) NOT NULL DEFAULT 'TWD' COMMENT '货币单位',
                           `area_code` varchar(8) NOT NULL DEFAULT '886' COMMENT '电话区号',
                           `created_by` int(11) NOT NULL DEFAULT 1 COMMENT '创建人',
                           `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
                           `updated_by` int(11) NOT NULL DEFAULT 1 COMMENT '编辑人',
                           `updated_at` int(11) NOT NULL DEFAULT 0 COMMENT '编辑时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `region_region1` (`region`),
                           UNIQUE KEY `region_currency1` (`currency`),
                           UNIQUE KEY `region_area_code1` (`area_code`),
                           KEY `region_region` (`id`) USING BTREE,
                           KEY `region_currency` (`id`) USING BTREE,
                           KEY `region_area_code` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
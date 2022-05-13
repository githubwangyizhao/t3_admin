CREATE TABLE `player_infos`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `server_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '区服id',
    `player_id` int(11) NOT NULL COMMENT '玩家id',
     `platform_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '平台',
     `pay_times` int(11) NOT NULL DEFAULT 99999 COMMENT '支付成功指定次数后展示第三方支付',
     `created_at` datetime(0) NOT NULL,
     `updated_at` datetime(0) NOT NULL,
     PRIMARY KEY (`id`) USING BTREE,
     UNIQUE INDEX `unique_key`(`server_id`, `player_id`, `platform_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户支付一定次数后，展示第三方次数' ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;
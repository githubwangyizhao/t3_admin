CREATE TABLE `client_heartbeat_verify`  (
                                           `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
                                           `platform` varchar(125) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT "local" COMMENT '平台',
                                           `server_id` varchar(125) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT "s1" COMMENT '区服',
                                           `start_date` datetime(0) NOT NULL DEFAULT NOW() COMMENT '匹配起始时间',
                                           `interval` smallint(4) UNSIGNED NOT NULL DEFAULT 100 COMMENT '有效期(单位:s)',
                                           `status` tinyint(1) UNSIGNED NOT NULL DEFAULT 2 COMMENT '状态(1为禁用,2为启用,默认2)',
                                           `created_by` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建者',
                                           `created_at` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
                                           `updated_by` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最后编辑者',
                                           `updated_at` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最后编辑时间',
                                           PRIMARY KEY (`id`),
                                           UNIQUE INDEX `unique`(`platform`, `server_id`) USING BTREE,
                                           INDEX `platform`(`platform`) USING BTREE,
                                           INDEX `serverId`(`server_id`) USING BTREE
);
CREATE TABLE `app_notice`  (
    `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `app_id` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'app_id',
    `type` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '类型(0为登录页公告,默认0)',
    `version` varchar(16) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '版本号',
    `notice` varchar(10240) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '通告内容',
    `stats` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态(1为启用,2为禁用,默认2)',
    `created_at` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间戳',
    `created_by` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建者',
    `updated_at` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '编辑时间戳',
    `updated_by` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '编辑者',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `app_version_unique`(`app_id`, `version`, `type`) USING BTREE COMMENT 'app_version_unique'
) COMMENT = 'app公告';

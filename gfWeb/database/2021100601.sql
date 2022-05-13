ALTER TABLE `platform_client_info`
    ADD COLUMN `reviewing_versions` varchar(12) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '当app在审核时前端js资源版本号' AFTER `versions`;
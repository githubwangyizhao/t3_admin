ALTER TABLE `platform_client_info`
    ADD COLUMN `package_size` float(10, 2) UNSIGNED NOT NULL DEFAULT 100.00 COMMENT '前端资源大小(单位mb)' AFTER `region`;
ALTER TABLE `platform_client_info`
    ADD COLUMN `area_code` varchar(16) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '886' COMMENT '隶属国家/区号' AFTER `region`;
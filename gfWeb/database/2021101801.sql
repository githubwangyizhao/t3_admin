ALTER TABLE `db_t3_admin`.`server_rel`
    CHANGE COLUMN `s_account` `remark` varchar(15) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL AFTER `id`,
    ADD COLUMN `s_account` varchar(255) NOT NULL COMMENT '客服帐号myadmin_user.id' AFTER `uid`,
    DROP PRIMARY KEY,
    ADD PRIMARY KEY (`id`) USING BTREE,
    DROP INDEX `server_rel_UN`,
    DROP INDEX `server_rel_id_IDX`,
    ADD UNIQUE INDEX `u_s`(`uid`, `s_account`) USING BTREE;
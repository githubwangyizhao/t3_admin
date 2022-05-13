CREATE TABLE `statistic_res`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '包名',
  `url` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL COMMENT '.zip结尾的地址',
  `version` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL COMMENT '静态资源ID',
  `add_time` int(64) NOT NULL,
  `uptime` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `statistic_res_UN`(`url`, `app_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 33 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '静态资源管理表' ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;

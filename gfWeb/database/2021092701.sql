-- db_t3_admin.server_rel definition

CREATE TABLE `server_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `s_account` varchar(15) NOT NULL COMMENT '客服帐号',
  `uid` int(11) NOT NULL,
  `plat_id` varchar(100) NOT NULL,
  `server_id` varchar(100) NOT NULL,
  `create_time` int(11) DEFAULT NULL,
  `uptime` varchar(100) NOT NULL,
  PRIMARY KEY (`s_account`,`uid`),
  UNIQUE KEY `server_rel_UN` (`uid`),
  KEY `server_rel_id_IDX` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8;
# noinspection SqlNoDataSourceInspectionForFile

# DROP TABLE IF EXISTS `myadmin_role`;
CREATE TABLE `myadmin_role`
(
    `id`   INTEGER      NOT NULL AUTO_INCREMENT COMMENT '角色id',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '角色名称',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `channel`;
CREATE TABLE `channel`
(
    `id`          INTEGER      NOT NULL AUTO_INCREMENT,
    `platform_id` varchar(64)  NOT NULL COMMENT '平台',
    `channel`     varchar(64)  NOT NULL COMMENT '渠道',
    `name`        varchar(256) NOT NULL COMMENT '渠道名称',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='渠道';

# DROP TABLE IF EXISTS `myadmin_role_channel_rel`;
CREATE TABLE `myadmin_role_channel_rel`
(
    `id`         INTEGER  NOT NULL AUTO_INCREMENT,
    `role_id`    INTEGER  NOT NULL,
    `channel_id` INTEGER  NOT NULL,
    `created`    datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `platform`;
CREATE TABLE `platform`
(
    `id`                       varchar(64)  NOT NULL COMMENT '平台',
    `name`                     varchar(256) NOT NULL COMMENT '平台名称',
    `inventory_database_id`    INTEGER      NOT NULL COMMENT '数据库',
    `zone_inventory_server_id` INTEGER      NOT NULL COMMENT '跨服服务器',
    `is_auto_open_server`      INTEGER      NOT NULL COMMENT '是否自动开服',
    `create_role_limit`        INTEGER      NOT NULL COMMENT '创角人数开服',
    `open_server_take_time`    INTEGER      NOT NULL DEFAULT 0 COMMENT '定时开服时间',
    `interval_init_time`       INTEGER      NOT NULL DEFAULT 0 COMMENT '间隔开服初始时间',
    `interval_day`             INTEGER      NOT NULL DEFAULT 1 COMMENT '间隔几天开服',
    `open_server_time_scope`   varchar(64)  NOT NULL DEFAULT '' COMMENT '开服时间段24小时制',
    `server_alias_str`         varchar(128) NOT NULL DEFAULT '' COMMENT '区服别名xx%dxx',
    `open_server_change_time`  INTEGER      NOT NULL DEFAULT 0 COMMENT '开服管理操作时间',
    `cron_update_time`         INTEGER      NOT NULL DEFAULT 0 COMMENT '订时更新时间',
    `cron_update_type`         TINYINT      NOT NULL DEFAULT 0 COMMENT '更新类型(冷热更)',
    `enter_state_type`         TINYINT      NOT NULL DEFAULT 0 COMMENT '入口状态类型',
    `enter_state`              TINYINT      NOT NULL DEFAULT 0 COMMENT '入口状态',
    `update_state`             TINYINT      NOT NULL DEFAULT 0 COMMENT '更新状态',
    `cron_update_user_id`      INTEGER      NOT NULL DEFAULT 0 COMMENT '更新操作者',
    `time`                     INTEGER      NOT NULL COMMENT '操作时间',
    `version`                  varchar(64)  NOT NULL DEFAULT '' COMMENT '版本库',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='平台';

# DROP TABLE IF EXISTS `background_charge_log`;
CREATE TABLE `background_charge_log`
(
    `id`           INTEGER          NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `platform_id`  varchar(64)      NOT NULL COMMENT '平台id',
    `server_id`    varchar(64)      NOT NULL COMMENT '区服id',
    `player_id`    int(10) unsigned NOT NULL COMMENT '玩家id',
    `charge_value` INTEGER                   DEFAULT NULL COMMENT '充值元宝',
    `charge_type`  varchar(32)      NOT NULL DEFAULT '0' COMMENT '充值类型',
    `item_id`      INTEGER          NOT NULL DEFAULT '0' COMMENT '充值物品',
    `ingot`        INTEGER          NOT NULL DEFAULT '0' COMMENT '元宝',
    `time`         INTEGER          NOT NULL DEFAULT '0' COMMENT '充值时间',
    `user_id`      INTEGER                   DEFAULT NULL COMMENT '用户id',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='后台充值日志';

# DROP TABLE IF EXISTS `forbid_log`;
CREATE TABLE `forbid_log`
(
    `platform_id` varchar(64) NOT NULL COMMENT '平台',
    `server_id`   varchar(64) NOT NULL COMMENT '区服',
    `player_id`   INTEGER     NOT NULL COMMENT '玩家ID',
    `forbid_type` tinyint(4)  NOT NULL COMMENT '封禁类型[0:禁言 1:封号]',
    `forbid_time` INTEGER     NOT NULL COMMENT '封禁时间',
    `time`        INTEGER     NOT NULL COMMENT '操作时间',
    `user_id`     INTEGER     NOT NULL COMMENT '操作用户id',
    PRIMARY KEY (`platform_id`, `server_id`, `player_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='封禁日志';

# DROP TABLE IF EXISTS `login_notice`;
CREATE TABLE `login_notice`
(
    `id`          varchar(128)   NOT NULL COMMENT '编号',
    `platform_id` varchar(64)    NOT NULL COMMENT '平台',
    `channel_id`  varchar(64)    NOT NULL COMMENT '渠道',
    `notice`      varchar(10240) NOT NULL COMMENT '公告内容',
    `time`        INTEGER        NOT NULL COMMENT '操作时间',
    `user_id`     INTEGER        NOT NULL COMMENT '用户id',
    PRIMARY KEY (`id`) USING BTREE,
    index (`platform_id`, `channel_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='登录公告';

# DROP TABLE IF EXISTS `mail_log`;
CREATE TABLE `mail_log`
(
    `id`               INTEGER        NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `platform_id`      varchar(64)    NOT NULL COMMENT '平台',
    `node_list`        varchar(10240) NOT NULL COMMENT '游戏节点列表',
    `title`            varchar(64)    NOT NULL COMMENT '标题',
    `content`          varchar(526)   NOT NULL COMMENT '内容',
    `item_list`        varchar(526)   NOT NULL COMMENT '道具列表',
    `time`             INTEGER        NOT NULL COMMENT '发送时间',
    `user_id`          INTEGER        NOT NULL COMMENT '用户id',
    `player_name_list` varchar(526) DEFAULT NULL COMMENT '发送玩家列表',
    `status`           int(255)     DEFAULT NULL COMMENT '状态',
    `type`             varchar(255) DEFAULT NULL COMMENT '类型(1: 发给玩家 2: 发给多个服 3:全服)',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='邮件日志';


# DROP TABLE IF EXISTS `myadmin_inventory_database`;
CREATE TABLE `myadmin_inventory_database`
(
    `id`          INTEGER     NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name`        varchar(64) NOT NULL COMMENT '名称',
    `host`        varchar(64) NOT NULL COMMENT '地址',
    `user`        varchar(64) NOT NULL COMMENT '用户',
    `port`        INTEGER DEFAULT NULL COMMENT '端口',
    `add_time`    INTEGER     NOT NULL COMMENT '添加时间',
    `update_time` INTEGER     NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='资产-数据库';

# DROP TABLE IF EXISTS `myadmin_inventory_server`;
CREATE TABLE `myadmin_inventory_server`
(
    `id`             INTEGER     NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name`           varchar(64) NOT NULL COMMENT '名称',
    `out_ip`         varchar(64) NOT NULL COMMENT '外网ip',
    `host`           varchar(64) NOT NULL COMMENT '域名',
    `inner_ip`       varchar(64) NOT NULL COMMENT '内网ip',
    `type`           INTEGER     NOT NULL COMMENT '类型 1:控制服务器 2:公共服务器 3:跨服服务器 4:游戏服务器',
    `add_time`       INTEGER     NOT NULL COMMENT '添加时间',
    `update_time`    INTEGER     NOT NULL COMMENT '更新时间',
    `max_node_count` INTEGER     NOT NULL DEFAULT '0' COMMENT '最大节点数量',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='资产-服务器';

# DROP TABLE IF EXISTS `myadmin_menu`;
CREATE TABLE `myadmin_menu`
(
    `id`        INTEGER     NOT NULL AUTO_INCREMENT COMMENT '菜单id',
    `name`      varchar(64) NOT NULL DEFAULT '' COMMENT '标识',
    `title`     varchar(64) NOT NULL DEFAULT '' COMMENT '标题',
    `parent_id` INTEGER              DEFAULT NULL COMMENT '父菜单id',
    `seq`       INTEGER     NOT NULL DEFAULT '0' COMMENT '序号',
    `is_show`   INTEGER              DEFAULT NULL COMMENT '是否显示',
    `icon`      varchar(32) NOT NULL DEFAULT '' COMMENT '图标',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `myadmin_resource`;
CREATE TABLE `myadmin_resource`
(
    `id`        INTEGER      NOT NULL AUTO_INCREMENT COMMENT '资源id',
    `name`      varchar(64)  NOT NULL DEFAULT '' COMMENT '资源名称',
    `parent_id` INTEGER               DEFAULT NULL COMMENT '父资源id',
    `url_for`   varchar(256) NOT NULL DEFAULT '' COMMENT '地址',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `myadmin_role_menu_rel`;
CREATE TABLE `myadmin_role_menu_rel`
(
    `id`      INTEGER  NOT NULL AUTO_INCREMENT,
    `role_id` INTEGER  NOT NULL,
    `menu_id` INTEGER  NOT NULL,
    `created` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `myadmin_role_platform_rel`;
CREATE TABLE `myadmin_role_platform_rel`
(
    `id`          INTEGER      NOT NULL AUTO_INCREMENT,
    `role_id`     INTEGER      NOT NULL,
    `platform_id` varchar(256) NOT NULL,
    `created`     datetime     NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `myadmin_role_resource_rel`;
CREATE TABLE `myadmin_role_resource_rel`
(
    `id`          INTEGER  NOT NULL AUTO_INCREMENT,
    `role_id`     INTEGER  NOT NULL,
    `resource_id` INTEGER  NOT NULL,
    `created`     datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `myadmin_role_user_rel`;
CREATE TABLE `myadmin_role_user_rel`
(
    `id`      INTEGER  NOT NULL AUTO_INCREMENT,
    `role_id` INTEGER  NOT NULL,
    `user_id` INTEGER  NOT NULL,
    `created` datetime NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `myadmin_user`;
CREATE TABLE `myadmin_user`
(
    `id`                         INTEGER      NOT NULL AUTO_INCREMENT COMMENT '帐号id',
    `name`                       varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
    `account`                    varchar(255) NOT NULL DEFAULT '' COMMENT '登录帐号',
    `password`                   varchar(255) NOT NULL DEFAULT '' COMMENT '登录密码',
    `status`                     INTEGER      NOT NULL DEFAULT '0' COMMENT '状态',
    `mobile`                     varchar(16)  NOT NULL DEFAULT '手机号',
    `login_times`                INTEGER      NOT NULL DEFAULT '0' COMMENT '登录次数',
    `last_login_time`            INTEGER      NOT NULL DEFAULT '0' COMMENT '最近登录时间',
    `last_login_ip`              varchar(64)  NOT NULL DEFAULT '0' COMMENT '最近登录',
    `mail_str`                   varchar(128) NOT NULL DEFAULT '' COMMENT '用户邮箱',
    `is_super`                   INTEGER               DEFAULT NULL COMMENT '是否超级管理员',
    `can_login_time`             INTEGER               DEFAULT '0' COMMENT '允许登录的时间',
    `continue_login_error_times` INTEGER               DEFAULT '0' COMMENT '连续登录失败次数',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='账号信息';

# DROP TABLE IF EXISTS `myadmin_user_myadmin_roles`;
CREATE TABLE `myadmin_user_myadmin_roles`
(
    `id`              bigint(20) NOT NULL AUTO_INCREMENT,
    `myadmin_user_id` INTEGER    NOT NULL,
    `myadmin_role_id` INTEGER    NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;

# DROP TABLE IF EXISTS `notice_log`;
CREATE TABLE `notice_log`
(
    `id`               INTEGER        NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `platform_id`      varchar(64)    NOT NULL COMMENT '平台',
    `server_id_list`   varchar(10240) NOT NULL COMMENT '区服列表',
    `content`          varchar(526)   NOT NULL COMMENT '公告内容',
    `notice_type`      tinyint(4)     NOT NULL COMMENT '类型[0:立即发送 1:定时发送 2:循环发送]',
    `notice_time`      INTEGER        NOT NULL COMMENT '公告执行时间',
    `create_cron_time` INTEGER        NOT NULL default 0 COMMENT '创建定时时间',
    `create_user_id`   INTEGER        NOT NULL COMMENT '创建用户id',
    `cron_time_str`    varchar(64)    NOT NULL default '' COMMENT '定时时间内容',
    `cron_times`       INTEGER        NOT NULL default 1 COMMENT '定时发送次数',
    `send_times`       INTEGER        NOT NULL default 0 COMMENT '已发送次数',
    `time`             INTEGER        NOT NULL COMMENT '操作时间',
    `user_id`          INTEGER        NOT NULL COMMENT '操作用户id',
    `status`           tinyint(4)              DEFAULT 0 COMMENT '状态[1:已执行]',
    `last_send_time`   INTEGER unsigned        DEFAULT NULL COMMENT '上次发送时间',
    `is_all_server`    INTEGER                 DEFAULT 0 COMMENT '是否全服发送',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='公告日志';

# DROP TABLE IF EXISTS `platform_inventory_server_rel`;
CREATE TABLE `platform_inventory_server_rel`
(
    `id`                  INTEGER      NOT NULL AUTO_INCREMENT,
    `platform_id`         varchar(256) NOT NULL,
    `inventory_server_id` INTEGER      NOT NULL,
    `created`             datetime     NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT;


# DROP TABLE IF EXISTS `remain_active`;
CREATE TABLE `remain_active`
(
    `node`        varchar(64)      NOT NULL COMMENT '节点',
    `time`        int(10) unsigned NOT NULL COMMENT '日期',
    `remain2`     INTEGER          NOT NULL DEFAULT '-1' COMMENT '次日留存',
    `remain3`     INTEGER          NOT NULL DEFAULT '-1' COMMENT '3日留存',
    `remain4`     INTEGER          NOT NULL DEFAULT '-1' COMMENT '4留存',
    `remain5`     INTEGER          NOT NULL DEFAULT '-1' COMMENT '5日留存',
    `remain6`     INTEGER          NOT NULL DEFAULT '-1' COMMENT '6日留存',
    `remain7`     INTEGER          NOT NULL DEFAULT '-1' COMMENT '7日留存',
    `remain8`     INTEGER          NOT NULL DEFAULT '-1' COMMENT '8日留存',
    `remain9`     INTEGER          NOT NULL DEFAULT '-1' COMMENT '9日留存',
    `remain10`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '10日留存',
    `remain11`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '11日留存',
    `remain12`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '12日留存',
    `remain13`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '13日留存',
    `remain14`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '14日留存',
    `remain15`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '15日留存',
    `remain16`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '16日留存',
    `remain17`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '17日留存',
    `remain18`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '18日留存',
    `remain19`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '19日留存',
    `remain20`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '20日留存',
    `remain21`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '21日留存',
    `remain22`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '22日留存',
    `remain23`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '23日留存',
    `remain24`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '24日留存',
    `remain25`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '25日留存',
    `remain26`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '26日留存',
    `remain27`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '27日留存',
    `remain28`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '28日留存',
    `remain29`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '29日留存',
    `remain30`    INTEGER          NOT NULL DEFAULT '-1' COMMENT '30日留存',
    `platform_id` varchar(64)      NOT NULL COMMENT '平台',
    `server_id`   varchar(64)      NOT NULL COMMENT '区服id',
    `channel`     varchar(64)      NOT NULL COMMENT '渠道',
    PRIMARY KEY (`node`, `time`, `platform_id`, `server_id`, `channel`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='活跃留存';

# DROP TABLE IF EXISTS `remain_charge`;
CREATE TABLE `remain_charge`
(
    `platform_id` varchar(64) NOT NULL DEFAULT '' COMMENT '平台id',
    `server_id`   varchar(64) NOT NULL DEFAULT '' COMMENT '区服id',
    `channel`     varchar(64) NOT NULL DEFAULT '' COMMENT '渠道',
    `time`        INTEGER     NOT NULL COMMENT '时间戳',
    `charge_num`  INTEGER     NOT NULL COMMENT '充值人数',
    `remain2`     INTEGER     NOT NULL COMMENT 'remain2',
    `remain3`     INTEGER     NOT NULL COMMENT 'remain3',
    `remain4`     INTEGER     NOT NULL COMMENT 'remain4',
    `remain7`     INTEGER     NOT NULL COMMENT 'remain7',
    `remain14`    INTEGER     NOT NULL COMMENT 'remain14',
    `remain30`    INTEGER     NOT NULL COMMENT 'remain30',
    PRIMARY KEY (`time`, `platform_id`, `server_id`, `channel`) USING BTREE,
    KEY `idx_remain_charge_1` (`platform_id`, `server_id`, `channel`, `time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='付费次留';


# DROP TABLE IF EXISTS `ten_minute_statistics`;
CREATE TABLE `ten_minute_statistics`
(
    `channel`             varchar(64) NOT NULL DEFAULT '' COMMENT '渠道',
    `time`                INTEGER     NOT NULL COMMENT '时间戳',
    `online_count`        INTEGER     NOT NULL COMMENT '在线人数',
    `register_count`      INTEGER     NOT NULL COMMENT '注册人数',
    `create_role_count`   INTEGER     NOT NULL DEFAULT '0' COMMENT '创角人数',
    `charge_count`        INTEGER     NOT NULL DEFAULT '0' COMMENT '付费金额',
    `charge_player_count` INTEGER     NOT NULL DEFAULT '0' COMMENT '付费人数',
    `platform_id`         varchar(64) NOT NULL DEFAULT '' COMMENT '平台id',
    `server_id`           varchar(64) NOT NULL DEFAULT '' COMMENT '区服id',
    PRIMARY KEY (`time`, `platform_id`, `server_id`, `channel`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='每10分钟统计';

# DROP TABLE IF EXISTS `daily_statistics`;
CREATE TABLE `daily_statistics`
(
    `platform_id`               varchar(64) NOT NULL DEFAULT '' COMMENT '平台id',
    `server_id`                 varchar(64) NOT NULL DEFAULT '' COMMENT '区服id',
    `channel`                   varchar(64) NOT NULL DEFAULT '' COMMENT '渠道',
    `time`                      INTEGER     NOT NULL COMMENT '时间戳',
    `charge_money`              INTEGER     NOT NULL COMMENT '当日充值人民币',
    `new_charge_money`          INTEGER     NOT NULL COMMENT '当日新增付费',
    `total_charge_money`        INTEGER     NOT NULL COMMENT '累计充值人民币',
    `charge_player_count`       INTEGER     NOT NULL COMMENT '当日充值人数',
    `total_charge_player_count` INTEGER     NOT NULL COMMENT '累计充值人数',
    `new_charge_player_count`   INTEGER              DEFAULT NULL COMMENT '当日首充人数',
    `login_times`               INTEGER              DEFAULT NULL COMMENT '当日登录次数',
    `login_player_count`        INTEGER              DEFAULT NULL COMMENT '当日登录用户数',
    `active_player_count`       INTEGER              DEFAULT NULL COMMENT '当日活跃用户',
    `create_role_count`         INTEGER              DEFAULT NULL COMMENT '当日创角',
    `share_create_role_count`   INTEGER              DEFAULT NULL COMMENT '当日分享创角',
    `total_create_role_count`   INTEGER              DEFAULT '0' COMMENT '累计创角',
    `max_online_count`          INTEGER              DEFAULT NULL COMMENT '当日最高在线人数',
    `min_online_count`          INTEGER              DEFAULT NULL COMMENT '当日最低在线人数',
    `avg_online_count`          INTEGER              DEFAULT NULL COMMENT '平均在线人数',
    `avg_online_time`           INTEGER              DEFAULT NULL COMMENT '平均在线时长',
    `register_count`            INTEGER              DEFAULT NULL COMMENT '当日注册',
    `total_register_count`      INTEGER              DEFAULT NULL COMMENT '累计注册',
    `valid_role_count`          INTEGER              DEFAULT NULL COMMENT '当日有效创角',
    `first_charge_player_count` INTEGER     NOT NULL DEFAULT '0' COMMENT '首充人数',
    `source`                    INTEGER     NOT NULL DEFAULT '0' COMMENT '3为谷歌支付，1为第三方支付通道（装备交易平台支付），2为苹果支付',
    PRIMARY KEY (`time`, `platform_id`, `server_id`, `channel`, `source`) USING BTREE,
    KEY `idx_daily_statistics_1` (`platform_id`, `server_id`, `channel`, `time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='每日汇总';

# DROP TABLE IF EXISTS `remain_total`;
CREATE TABLE `remain_total`
(
    `time`          int(10) unsigned NOT NULL COMMENT '日期',
    `remain2`       INTEGER          NOT NULL DEFAULT '-1' COMMENT '次日留存',
    `remain3`       INTEGER          NOT NULL DEFAULT '-1' COMMENT '3日留存',
    `remain4`       INTEGER          NOT NULL DEFAULT '-1' COMMENT '4留存',
    `remain5`       INTEGER          NOT NULL DEFAULT '-1' COMMENT '5日留存',
    `remain6`       INTEGER          NOT NULL DEFAULT '-1' COMMENT '6日留存',
    `remain7`       INTEGER          NOT NULL DEFAULT '-1' COMMENT '7日留存',
    `remain8`       INTEGER          NOT NULL DEFAULT '-1' COMMENT '8日留存',
    `remain9`       INTEGER          NOT NULL DEFAULT '-1' COMMENT '9日留存',
    `remain10`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '10日留存',
    `remain11`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '11日留存',
    `remain12`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '12日留存',
    `remain13`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '13日留存',
    `remain14`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '14日留存',
    `remain15`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '15日留存',
    `remain16`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '16日留存',
    `remain17`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '17日留存',
    `remain18`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '18日留存',
    `remain19`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '19日留存',
    `remain20`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '20日留存',
    `remain21`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '21日留存',
    `remain22`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '22日留存',
    `remain23`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '23日留存',
    `remain24`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '24日留存',
    `remain25`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '25日留存',
    `remain26`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '26日留存',
    `remain27`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '27日留存',
    `remain28`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '28日留存',
    `remain29`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '29日留存',
    `remain30`      INTEGER          NOT NULL DEFAULT '-1' COMMENT '30日留存',
    `platform_id`   varchar(64)      NOT NULL COMMENT '平台',
    `server_id`     varchar(64)      NOT NULL COMMENT '区服id',
    `channel`       varchar(64)      NOT NULL COMMENT '渠道',
    `create_role`   INTEGER          NOT NULL DEFAULT '0' COMMENT '创角人数',
    `register_role` INTEGER          NOT NULL DEFAULT '0' COMMENT '注册人数',
    PRIMARY KEY (`time`, `platform_id`, `server_id`, `channel`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='总体留存';

# DROP TABLE IF EXISTS `daily_ltv`;
CREATE TABLE `daily_ltv`
(
    `platform_id` varchar(64) NOT NULL DEFAULT '' COMMENT '平台id',
    `server_id`   varchar(64) NOT NULL DEFAULT '' COMMENT '区服id',
    `channel`     varchar(64) NOT NULL DEFAULT '' COMMENT '渠道',
    `time`        INTEGER     NOT NULL COMMENT '时间戳',
    `c1`          INTEGER     NOT NULL COMMENT 'c1',
    `c2`          INTEGER     NOT NULL COMMENT 'c2',
    `c3`          INTEGER     NOT NULL COMMENT 'c3',
    `c7`          INTEGER     NOT NULL COMMENT 'c7',
    `c14`         INTEGER     NOT NULL COMMENT 'c14',
    `c30`         INTEGER     NOT NULL COMMENT 'c30',
    `c60`         INTEGER     NOT NULL COMMENT 'c60',
    `c90`         INTEGER     NOT NULL COMMENT 'c90',
    `c120`        INTEGER     NOT NULL COMMENT 'c120',
    PRIMARY KEY (`time`, `platform_id`, `server_id`, `channel`) USING BTREE,
    KEY `idx_daily_ltv_1` (`platform_id`, `server_id`, `channel`, `time`),
    KEY `idx_daily_ltv_2` (`platform_id`, `channel`, `time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='每日LTV';

# DROP TABLE IF EXISTS `platform_merge_server_data`;
CREATE TABLE `platform_merge_server_data`
(
    `platform_id`    varchar(64)   NOT NULL DEFAULT '' COMMENT '平台id',
    `merge_time`     INTEGER       NOT NULL COMMENT '合服时间',
    `merge_state`    TINYINT       NOT NULL DEFAULT 0 COMMENT '合服状态1:请求合服,3:审核通过,4:失败,5:合服中,9:合服完成',
    `merge_str`      varchar(2048) NOT NULL DEFAULT '' COMMENT '合服数据',
    `fail_msg`       varchar(1024) NOT NULL DEFAULT '' COMMENT '失败内容',
    `request_id`     INTEGER       NOT NULL DEFAULT 0 COMMENT '申请者id',
    `audit_id`       INTEGER       NOT NULL DEFAULT 0 COMMENT '审核id',
    `request_time`   INTEGER       NOT NULL DEFAULT 0 COMMENT '申请者时间',
    `audit_time`     INTEGER       NOT NULL DEFAULT 0 COMMENT '审核时间',
    `merge_use_time` INTEGER       NOT NULL DEFAULT 0 COMMENT '合服用时',
    PRIMARY KEY (`platform_id`, `merge_time`) USING BTREE,
    KEY `idx_merger_1` (`platform_id`, `merge_state`),
    KEY `idx_merge_2` (`merge_state`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='平台合服数据';

# DROP TABLE IF EXISTS `platform_ding_yue`;
CREATE TABLE `platform_ding_yue`
(
    `platform_id`   varchar(64) NOT NULL DEFAULT '' COMMENT '平台id',
    `ding_yue_time` INTEGER     NOT NULL COMMENT '订阅时间',
    `ding_yue_num`  INTEGER     NOT NULL DEFAULT 0 COMMENT '订阅数量',
    PRIMARY KEY (`platform_id`, `ding_yue_time`) USING BTREE,
    KEY `idx_merger_1` (`platform_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='平台订阅数据';

# DROP TABLE IF EXISTS `platform_acc_statistics`;
CREATE TABLE `platform_acc_statistics`
(
    `platform_id`   varchar(64) NOT NULL DEFAULT '' COMMENT '平台id',
    `acc_id`        varchar(64) NOT NULL DEFAULT '' COMMENT '账号',
    `ding_yue_time` INTEGER     NOT NULL COMMENT '订阅时间',
    PRIMARY KEY (`platform_id`, `acc_id`) USING BTREE,
    KEY `idx_merger_1` (`platform_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='平台账号统计数据';


# DROP TABLE IF EXISTS `mail_smtp_data`;
CREATE TABLE `mail_smtp_data`
(
    `user`  varchar(255) NOT NULL DEFAULT '' COMMENT '发送者mail',
    `pass`  varchar(255) NOT NULL DEFAULT '' COMMENT '发送者mail密码',
    `host`  varchar(64)  NOT NULL DEFAULT '' COMMENT '中间host',
    `port`  INTEGER      NOT NULL DEFAULT 0 COMMENT '端口',
    `state` TINYINT      NOT NULL DEFAULT 0 COMMENT '开启状态',
    PRIMARY KEY (`user`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='邮件smtp配置';

# DROP TABLE IF EXISTS `setting_data`;
CREATE TABLE `setting_data`
(
    `id`          INTEGER     NOT NULL DEFAULT 0 COMMENT '设置id',
    `name`        varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
    `state`       TINYINT     NOT NULL DEFAULT 0 COMMENT '状态1:开启',
    `change_time` INTEGER     NOT NULL DEFAULT 0 COMMENT '操作时间',
    `user_id`     INTEGER     NOT NULL DEFAULT 0 COMMENT '操作用户id',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='后台设置数据';

# DROP TABLE IF EXISTS `page_change_auth`;
CREATE TABLE `page_change_auth`
(
    `id`          INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `sign`        varchar(64)      NOT NULL DEFAULT '' COMMENT '标识',
    `name`        varchar(64)      NOT NULL DEFAULT '' COMMENT '名称',
    `state`       TINYINT          NOT NULL DEFAULT 0 COMMENT '状态1:开启',
    `change_time` INTEGER          NOT NULL DEFAULT 0 COMMENT '操作时间',
    `user_id`     INTEGER          NOT NULL DEFAULT 0 COMMENT '操作用户id',
    PRIMARY KEY (`id`) USING BTREE,
    index (`sign`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='界面操作权限';

# DROP TABLE IF EXISTS `background_msg_template`;
CREATE TABLE `background_msg_template`
(
    `id`             INTEGER     NOT NULL DEFAULT 0 COMMENT '设置id',
    `name`           varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
    `state`          TINYINT     NOT NULL DEFAULT 0 COMMENT '状态1:开启',
    `phone_msg_code` varchar(64) NOT NULL DEFAULT '' COMMENT '手机信息模板code',
    `change_time`    INTEGER     NOT NULL DEFAULT 0 COMMENT '操作时间',
    `user_id`        INTEGER     NOT NULL DEFAULT 0 COMMENT '操作用户id',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='后台信息数据模板';

# DROP TABLE IF EXISTS `user_rel_data`;
CREATE TABLE `user_rel_data`
(
    `type`        INTEGER NOT NULL DEFAULT 0 COMMENT '类型',
    `id`          INTEGER NOT NULL DEFAULT 0 COMMENT '类型对应数据id',
    `user_id`     INTEGER NOT NULL DEFAULT 0 COMMENT '用户id',
    `create_time` INTEGER NOT NULL DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY (`type`, `id`, `user_id`) USING BTREE,
    index (`type`, `id`) USING BTREE,
    index (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='用户关联数据';


# DROP TABLE IF EXISTS `request_admin_log`;
CREATE TABLE `request_admin_log`
(
    `id`          INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `user_id`     INTEGER          NOT NULL DEFAULT 0 COMMENT '用户id',
    `use_time`    DOUBLE           NOT NULL DEFAULT 0 COMMENT '用时',
    `url`         VARCHAR(256)     NOT NULL DEFAULT '' COMMENT '请求url',
    `param_str`   VARCHAR(20480)   NOT NULL DEFAULT '' COMMENT '参数内容',
    `login_type`  TINYINT          NOT NULL DEFAULT 0 COMMENT '登录类型1:登录2:未登录',
    `state`       TINYINT          NOT NULL DEFAULT 0 COMMENT '状态',
    `fail_msg`    VARCHAR(1024)    NOT NULL DEFAULT '' COMMENT '失败内容',
    `change_time` INTEGER          NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE,
    index (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='用户请求日志';

# DROP TABLE IF EXISTS `login_admin_log`;
CREATE TABLE `login_admin_log`
(
    `id`          INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `account`     varchar(255)     NOT NULL DEFAULT '' COMMENT '登录帐号',
    `user_name`   varchar(255)     NOT NULL DEFAULT '' COMMENT '登录名称',
    `ip`          varchar(64)      NOT NULL COMMENT 'ip',
    `login_type`  TINYINT          NOT NULL DEFAULT 0 COMMENT '登录类型1:登录2:刷新',
    `explorer`    varchar(64)               DEFAULT '' COMMENT '浏览器类型',
    `os`          varchar(64)               DEFAULT '' COMMENT '操作系统',
    `state`       TINYINT          NOT NULL DEFAULT 0 COMMENT '状态',
    `fail_msg`    VARCHAR(1024)    NOT NULL DEFAULT '' COMMENT '失败内容',
    `change_time` INTEGER          NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE,
    index (`account`) USING BTREE,
    index (`ip`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='用户登录日志';

# DROP TABLE IF EXISTS `open_server_manage_log`;
CREATE TABLE `open_server_manage_log`
(
    `id`          INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `platform_id` varchar(64)      NOT NULL DEFAULT '' COMMENT '平台id',
    `user_id`     INTEGER          NOT NULL DEFAULT 0 COMMENT '用户id',
    `change_str`  VARCHAR(512)     NOT NULL DEFAULT '' COMMENT '操作内容',
    `change_time` INTEGER          NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE,
    index (`user_id`) USING BTREE,
    index (`platform_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='开服管理日志';


# DROP TABLE IF EXISTS `update_platform_version_log`;
CREATE TABLE `update_platform_version_log`
(
    `id`          INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `platform_id` varchar(64)      NOT NULL DEFAULT '' COMMENT '平台id',
    `user_id`     INTEGER          NOT NULL DEFAULT 0 COMMENT '用户id',
    `state`       TINYINT          NOT NULL DEFAULT 0 COMMENT '状态',
    `type`        TINYINT          NOT NULL DEFAULT 0 COMMENT '更新类型0:手动1:订时',
    `use_time`    DOUBLE           NOT NULL DEFAULT 0 COMMENT '用时',
    `fail_msg`    VARCHAR(1024)    NOT NULL DEFAULT '' COMMENT '失败内容',
    `change_time` INTEGER          NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE,
    index (`user_id`) USING BTREE,
    index (`platform_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='平台版本更新日志';


# DROP TABLE IF EXISTS `branch_path`;
CREATE TABLE `branch_path`
(
    `id`          INTEGER       NOT NULL AUTO_INCREMENT COMMENT '版本编号id',
    `name`        varchar(64)   NOT NULL DEFAULT '' COMMENT '名字',
    `path`        VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '路径',
    `user_id`     INTEGER       NOT NULL DEFAULT 0 COMMENT '用户id',
    `change_time` INTEGER       NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='版本路径';

# DROP TABLE IF EXISTS `platform_version_path`;
CREATE TABLE `platform_version_path`
(
    `id`          INTEGER       NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `type`        TINYINT       NOT NULL DEFAULT 0 COMMENT '类型1:客户端 2:服务器端',
    `branch_id`   INTEGER       NOT NULL DEFAULT 0 COMMENT '版本编号',
    `change_type` TINYINT       NOT NULL DEFAULT 0 COMMENT '操作类型1:打包 2:同步 3更新',
    `state`       TINYINT       NOT NULL DEFAULT 0 COMMENT '就否显示',
    `name`        varchar(64)   NOT NULL DEFAULT '' COMMENT '名字',
    `path`        VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '路径',
    `user_id`     INTEGER       NOT NULL DEFAULT 0 COMMENT '用户id',
    `change_time` INTEGER       NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE,
    index (`type`, `change_type`, `branch_id`) USING BTREE,
    index (`type`, `change_type`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='平台版本路径';

# DROP TABLE IF EXISTS `platform_version_change`;
CREATE TABLE `version_tool_change`
(
    `id`                  INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `type`                TINYINT          NOT NULL DEFAULT 0 COMMENT '渠道操作类型0:测试服1:正式服资源2:正式服生效',
    `platform_version_id` INTEGER          NOT NULL DEFAULT 0 COMMENT '平台版本路径id',
    `state`               TINYINT          NOT NULL DEFAULT 0 COMMENT '就否显示',
    `name`                varchar(64)      NOT NULL DEFAULT '' COMMENT '名字',
    `sh_path`             VARCHAR(512)     NOT NULL DEFAULT '' COMMENT 'sh路径',
    `user_id`             INTEGER          NOT NULL DEFAULT 0 COMMENT '用户id',
    `change_time`         INTEGER          NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='版本工具操作';

# DROP TABLE IF EXISTS `version_tool_change_log`;
CREATE TABLE `version_tool_change_log`
(
    `id`          INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `change_id`   INTEGER          NOT NULL DEFAULT 0 COMMENT '操作id',
    `state`       TINYINT          NOT NULL DEFAULT 0 COMMENT '状态',
    `update_type` TINYINT          NOT NULL DEFAULT 0 COMMENT '更新类型1:手动,2:订时',
    `use_time`    DOUBLE           NOT NULL DEFAULT 0 COMMENT '用时',
    `sh_name`     VARCHAR(256)     NOT NULL DEFAULT '' COMMENT '执行脚本',
    `fail_msg`    VARCHAR(1024)    NOT NULL DEFAULT '' COMMENT '失败内容',
    `body`        VARCHAR(1024)    NOT NULL DEFAULT '' COMMENT '执行内容 ',
    `user_id`     INTEGER          NOT NULL DEFAULT 0 COMMENT '用户id',
    `change_time` INTEGER          NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='版本工具操作日志';

# DROP TABLE IF EXISTS `version_tool_change_cron`;
CREATE TABLE `version_tool_change_cron`
(
    `id`             INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `change_id_str`  varchar(512)     NOT NULL DEFAULT '' COMMENT '版本操作编号id列表',
    `state`          TINYINT                   DEFAULT 1 COMMENT '状态1:可执行2:已完成',
    `robot_type`     TINYINT          NOT NULL DEFAULT 0 COMMENT '类型1:客户端 2:服务器端',
    `change_type`    TINYINT          NOT NULL DEFAULT 0 COMMENT '操作类型1:打包 2:同步 3更新',
    `cron_time_str`  varchar(64)      NOT NULL DEFAULT '' COMMENT '定时时间内容',
    `send_times`     INTEGER          NOT NULL default 0 COMMENT '已发送次数',
    `last_send_time` INTEGER                   DEFAULT 0 COMMENT '上次发送时间',
    `user_id`        INTEGER          NOT NULL DEFAULT 0 COMMENT '用户id',
    `change_time`    INTEGER          NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE,
    index (`state`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='订时版本工具操作';

CREATE TABLE `background_version`
(
    `version`     varchar(128) NOT NULL COMMENT '版本号',
    `change_time` INTEGER      NOT NULL DEFAULT 0 COMMENT '操作时间',
    PRIMARY KEY (`version`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='后台版本';

CREATE TABLE `promote_data`
(
    `id`         INTEGER unsigned NOT NULL AUTO_INCREMENT COMMENT '编号id',
    `promote`    varchar(32)      NOT NULL COMMENT '推广员标识',
    `name`       varchar(32)      NOT NULL DEFAULT '' COMMENT '推广员备注',
    `state`      TINYINT          NOT NULL DEFAULT 1 COMMENT '0:禁用;1:启用',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` INTEGER                   DEFAULT NULL COMMENT '创建者',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` INTEGER          NOT NULl DEFAULT 0 COMMENT '编辑者',
    PRIMARY KEY (`id`) USING BTREE,
    index (`promote`, `created_by`, `updated_by`, `state`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='推广员';

CREATE TABLE `client_versions`
(
    `platform_id`          varchar(32)  NOT NULL COMMENT '平台id',
    `platform_name`        varchar(16)  NOT NULL DEFAULT '' comment '平台中文标识',
    `android_download_url` varchar(255) NOT NULL DEFAULT 'https://goldenmaster1.s3-ap-southeast-1.amazonaws.com/goldenmaster.apk' comment '安卓冷更包下载地址',
    `ios_download_url`     varchar(255) NOT NULL DEFAULT 'https://iosdownloads.site/install/3dvhvzpa1gcd-goldmaster-292' comment 'iOS冷更包下载地址',
    `first_versions`       INTEGER      NOT NULL DEFAULT 0 COMMENT '初始版本',
    `client_version`       varchar(16)  NOT NULL DEFAULT '1.0' COMMENT '客户端版本号(格式:大版本.小版本)',
    `versions`             INTEGER      NOT NULL DEFAULT 0 COMMENT '版本',
    `ip`                   varchar(15)  NOT NULL DEFAULT '' COMMENT 'ip地址',
    `created_at`           datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by`           INTEGER      NOT NULl DEFAULT 0 COMMENT '创建者',
    `updated_at`           datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by`           INTEGER      NOT NULl DEFAULT 0 COMMENT '编辑者',
    `is_close_charge`      TINYINT      NOT NULl DEFAULT 0 COMMENT '是否关闭充值',
    PRIMARY KEY (`platform_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='客户端版本';

# alter table `client_versions` add column `is_close_charge` TINYINT not null default 0 comment '是否关闭充值' after updated_by;
# alter table `client_versions` add column `client_version` varchar(16) not null default '1.0' comment '客户端版本号(格式:大版本.小版本)' after `versions`;
# alter table `client_versions` add column `platform_name` varchar(16) not null default "" comment "平台中文标识" after `platform_id`;
# alter table `client_versions` add column `android_download_url` varchar(255) not null default "https://goldenmaster1.s3-ap-southeast-1.amazonaws.com/goldenmaster.apk" comment "安卓冷更包下载地址" after `platform_name`;
# alter table `client_versions` add column `ios_download_url` varchar(255) not null default "https://iosdownloads.site/install/3dvhvzpa1gcd-goldmaster-292" comment "iOS冷更包下载地址" after `android_download_url`;


-- ----------------------------
DROP TABLE IF EXISTS `platform_client_info`;
CREATE TABLE `platform_client_info`
(
    `id`                  int(10) unsigned NOT NULL AUTO_INCREMENT,
    `app_id`              varchar(50)      NOT NULL DEFAULT '',
    `facebook_app_id`     varchar(255)     NOT NULL DEFAULT '',
    `platform`            varchar(50)      NOT NULL DEFAULT 'indonesia' COMMENT '平台名(platform)',
    `platform_remark`     varchar(100)     NOT NULL DEFAULT '印度尼西亚' COMMENT '完整国家名',
    `client_version`      decimal(10, 2)   NOT NULL DEFAULT '4.16' COMMENT '客户端版本号',
    `first_versions`      varchar(12)      NOT NULL DEFAULT '2021030601' COMMENT '前端js资源初始版本号',
    `versions`            varchar(12)      NOT NULL DEFAULT '2021032402' COMMENT '前端js资源当前版本号',
    `is_charge_open`      tinyint(1) unsigned       DEFAULT '1' COMMENT '是否开启充值(0为否,1为是,默认1)',
    `native_pay`          tinyint(1) unsigned       DEFAULT '1' COMMENT '是否启用商店充值(0为否,1为是,默认1)',
    `upgrade_ios_url`     varchar(255)     NOT NULL DEFAULT '' COMMENT 'iOS冷更地址',
    `upgrade_android_url` varchar(255)     NOT NULL DEFAULT '' COMMENT '安卓冷更地址',
    `reload_url`          varchar(255)              DEFAULT NULL COMMENT '热更包地址',
    `stats`               tinyint(1) unsigned       DEFAULT '1' COMMENT '是否启用(0为否,1为是,默认1)',
    `created_by`          int              NOT NULL DEFAULT 1 COMMENT '创建者',
    `pay_times`           int(11)          NOT NULL DEFAULT '99999' COMMENT '支付成功指定次数后展示第三方支付',
    `created_at`          datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_by`          int              NOT NULL DEFAULT 1 COMMENT '最近一次编辑者',
    `updated_at`          datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`, `app_id`),
    UNIQUE KEY `app_id` (`app_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  ROW_FORMAT = COMPACT COMMENT ='渠道客户端信息';

# ALTER TABLE `db_t1_admin`.`platform_client_info` ADD COLUMN `pay_times` int(11) NOT NULL DEFAULT 99999 COMMENT '支付成功指定次数后展示第三方支付' AFTER `stats`;

DROP TABLE IF EXISTS `player_infos`;
CREATE TABLE `player_infos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `server_id` varchar(255) NOT NULL COMMENT '区服id',
  `player_id` int(11) NOT NULL COMMENT '玩家id',
  `platform_id` varchar(255) NOT NULL COMMENT '平台',
  `pay_times` int(11) NOT NULL DEFAULT '99999' COMMENT '支付成功指定次数后展示第三方支付',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_key` (`server_id`,`player_id`,`platform_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户支付一定次数后，展示第三方次数';

/**
 为platform_client_info表新增region字段 2021-08-18
 */
alter table platform_client_info
    add region varchar(50) default 'tw' not null comment '隶属国家/地区' after pay_times;
alter table platform_client_info alter column platform set default 'local';
alter table platform_client_info
    add channel varchar(50) default 'local' not null comment '渠道' after platform_remark;

CREATE TABLE `adjust` (
                                           `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
                                           `platform_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 0 COMMENT '平台',
                                           `server` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 0 COMMENT '区服',
                                           `type` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '类型(1为场景,2为玩家,默认1)',
                                           `ref_id` int(11) NOT NULL COMMENT '关联id(type为0时此处为场景id;type为1时此处为玩家id)',
                                           `value` int(11) NOT NULL DEFAULT 10000 COMMENT '修正值',
                                           `status` tinyint(1) UNSIGNED NOT NULL COMMENT '状态(1为禁用,2为启用,默认1)',
                                           `created_at` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
                                           `created_by` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建者',
                                           `updated_at` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最近更新者',
                                           `updated_by` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最近更新时间',
                                           PRIMARY KEY (`id`),
                                           UNIQUE INDEX `type_ref_id`(`platform_id`, `server`, `type`, `ref_id`) USING BTREE,
                                           INDEX `status`(`status`) USING BTREE,
                                           INDEX `platform_id`(`platform_id`) USING BTREE,
                                           INDEX `server`(`server`) USING BTREE,
                                           INDEX `type`(`type`) USING BTREE
) COMMENT = '修正值';
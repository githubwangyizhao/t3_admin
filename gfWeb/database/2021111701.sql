ALTER TABLE `myadmin_role_menu_rel`
    MODIFY COLUMN `created` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0);
ALTER TABLE `myadmin_role_channel_rel`
    MODIFY COLUMN `created` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0);
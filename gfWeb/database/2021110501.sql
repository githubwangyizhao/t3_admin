ALTER TABLE `platform_inventory_server_rel`
    MODIFY COLUMN `created` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) AFTER `inventory_server_id`;
ALTER TABLE `platform_inventory_server_rel`
    MODIFY COLUMN `created` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) AFTER `inventory_server_id`;
ALTER TABLE `channel` ADD COLUMN `campaign` varchar(50) not null default '0' COMMENT 'adjust的campaign' AFTER `name`
ALTER TABLE `channel` DROP COLUMN `campaign`
ALTER TABLE `channel` ADD COLUMN `tracker_token` varchar(50) not null default '0' COMMENT 'campaign的tracker_token'
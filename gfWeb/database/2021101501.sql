ALTER TABLE app_notice ADD COLUMN repeated tinyint(1) DEFAULT 1 NOT NULL COMMENT '每次打开app都显示。1为是,2为否,默认1' AFTER `notice`;
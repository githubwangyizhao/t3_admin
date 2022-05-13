ALTER TABLE daily_statistics MODIFY COLUMN charge_money decimal(20,2) NOT NULL COMMENT '当日充值人民币';
ALTER TABLE daily_statistics MODIFY COLUMN new_charge_money decimal(20,2) NOT NULL COMMENT '当日新增付费';
ALTER TABLE daily_statistics MODIFY COLUMN total_charge_money decimal(20,2) NOT NULL COMMENT '累计充值人民币';
ALTER TABLE daily_statistics MODIFY COLUMN first_charge_total_money decimal(20,2) DEFAULT 0 NOT NULL COMMENT '首充金额';
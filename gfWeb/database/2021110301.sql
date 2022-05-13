alter table `channel` add column currency varchar(20) not null default 'TWD' comment '该渠道隶属地区货币单位';
alter table `channel` add column area_code varchar(20) not null default '886' comment '该渠道隶属地区区号';
alter table `channel` add column region varchar(20) not null default '台灣' comment '渠道隶属国家'
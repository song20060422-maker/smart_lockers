-- ============================================================
-- 智能快寄柜 建表脚本 V1.0
-- 按工单归集；每表注明来源工单编号
-- ============================================================
CREATE DATABASE IF NOT EXISTS smart_lockers DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE smart_lockers;

-- ---------- 工单03 首页 ----------
CREATE TABLE IF NOT EXISTS city (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL COMMENT '城市名 @ticket 03-首页',
  pinyin VARCHAR(128),
  initial CHAR(2),
  is_hot TINYINT DEFAULT 0,
  sort INT DEFAULT 0,
  location VARCHAR(128),
  INDEX idx_initial (initial), INDEX idx_hot (is_hot)
) ENGINE=InnoDB COMMENT='城市 @ticket 03-首页';

CREATE TABLE IF NOT EXISTS locker_site (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  city_id BIGINT NOT NULL,
  name VARCHAR(128) NOT NULL,
  address VARCHAR(256),
  lat DECIMAL(10,7), lng DECIMAL(10,7),
  open_time VARCHAR(16), close_time VARCHAR(16),
  status TINYINT DEFAULT 1,
  INDEX idx_city (city_id), INDEX idx_lat_lng (lat, lng)
) ENGINE=InnoDB COMMENT='寄存点 @ticket 03-首页';

CREATE TABLE IF NOT EXISTS guide (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  version VARCHAR(32) NOT NULL,
  lang VARCHAR(16),
  content TEXT,
  is_current TINYINT DEFAULT 0,
  UNIQUE KEY uk_version (version)
) ENGINE=InnoDB COMMENT='寄存指南 @ticket 03-首页';

-- ---------- 工单04 物品寄存 / 05 订单 ----------
CREATE TABLE IF NOT EXISTS locker_cabinet (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  site_id BIGINT NOT NULL,
  code VARCHAR(32), size VARCHAR(16),
  total INT DEFAULT 0, available INT DEFAULT 0,
  status TINYINT DEFAULT 1,
  INDEX idx_site (site_id)
) ENGINE=InnoDB COMMENT='柜组 @ticket 04-物品寄存';

CREATE TABLE IF NOT EXISTS locker_door (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  cabinet_id BIGINT NOT NULL,
  code VARCHAR(32), size VARCHAR(16),
  status TINYINT DEFAULT 0 COMMENT '0空闲 1占用 2故障 @ticket 04-物品寄存',
  version INT DEFAULT 0,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_cabinet (cabinet_id), INDEX idx_status (status)
) ENGINE=InnoDB COMMENT='柜门 @ticket 04-物品寄存';

CREATE TABLE IF NOT EXISTS pricing_rule (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  size VARCHAR(16) NOT NULL,
  first_hours INT DEFAULT 1,
  first_fee BIGINT DEFAULT 0 COMMENT '分 @ticket 04-物品寄存',
  extra_hour_fee BIGINT DEFAULT 0,
  daily_cap BIGINT DEFAULT 0,
  deposit BIGINT DEFAULT 0,
  priority INT DEFAULT 0,
  enabled TINYINT DEFAULT 1,
  INDEX idx_size (size)
) ENGINE=InnoDB COMMENT='收费规则 @ticket 04-物品寄存';

CREATE TABLE IF NOT EXISTS deposit_order (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(64) NOT NULL,
  user_id BIGINT NOT NULL,
  door_id BIGINT NOT NULL,
  site_id BIGINT NOT NULL,
  start_time DATETIME, end_time DATETIME,
  amount BIGINT DEFAULT 0, deposit BIGINT DEFAULT 0,
  status TINYINT COMMENT '1进行中 2已完成 3已取消 @ticket 04/05',
  idempotency_key VARCHAR(128),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_order_no (order_no), UNIQUE KEY uk_idem (idempotency_key),
  INDEX idx_user_status (user_id, status), INDEX idx_door (door_id)
) ENGINE=InnoDB COMMENT='寄存订单 @ticket 04/05';

CREATE TABLE IF NOT EXISTS temp_open_auth (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(64) NOT NULL,
  token VARCHAR(128) NOT NULL,
  max_times INT DEFAULT 1, used_times INT DEFAULT 0,
  expire_at DATETIME NOT NULL, used TINYINT DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_token (token), INDEX idx_order (order_no)
) ENGINE=InnoDB COMMENT='临时开柜授权 @ticket 05-订单查看与取出';

CREATE TABLE IF NOT EXISTS makeup_record (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(64) NOT NULL,
  amount BIGINT NOT NULL, status TINYINT,
  idempotency_key VARCHAR(128),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_idem (idempotency_key), INDEX idx_order (order_no)
) ENGINE=InnoDB COMMENT='超时补缴记录 @ticket 05-订单查看与取出';

-- 历史订单归档分区建议 @ticket 05：按月 RANGE 分区（数据量上来后启用）
-- ALTER TABLE deposit_order PARTITION BY RANGE COLUMNS(created_at) (...);

-- ---------- 工单02 支付宝支付 ----------
CREATE TABLE IF NOT EXISTS account (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  balance BIGINT DEFAULT 0 COMMENT '分，以DB为准 @ticket 02-支付宝支付',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user (user_id)
) ENGINE=InnoDB COMMENT='账户余额 @ticket 02-支付宝支付';

CREATE TABLE IF NOT EXISTS balance_log (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  change_amount BIGINT NOT NULL,
  balance_after BIGINT NOT NULL,
  biz_type VARCHAR(32), biz_no VARCHAR(64),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_user (user_id), UNIQUE KEY uk_biz (biz_type, biz_no)
) ENGINE=InnoDB COMMENT='余额明细 @ticket 02-支付宝支付';

CREATE TABLE IF NOT EXISTS recharge_order (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  charge_no VARCHAR(64) NOT NULL,
  user_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  status TINYINT DEFAULT 1 COMMENT '1待支付 2已支付 3失败 @ticket 02',
  trade_no VARCHAR(128),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  paid_at DATETIME,
  UNIQUE KEY uk_charge (charge_no), INDEX idx_user (user_id)
) ENGINE=InnoDB COMMENT='充值订单 @ticket 02-支付宝支付';

-- ---------- 工单06 个人中心 ----------
CREATE TABLE IF NOT EXISTS user_profile (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  phone VARCHAR(32) NOT NULL,
  nickname VARCHAR(64), avatar VARCHAR(256), role VARCHAR(32),
  status TINYINT DEFAULT 1,
  token_ver INT DEFAULT 0 COMMENT '退出登录自增使旧token失效 @ticket 06',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_phone (phone)
) ENGINE=InnoDB COMMENT='用户资料 @ticket 06-个人中心';

CREATE TABLE IF NOT EXISTS coupon (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(64) NOT NULL,
  user_id BIGINT NOT NULL,
  discount BIGINT, scope VARCHAR(64),
  expire_at DATETIME NOT NULL, used TINYINT DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_code (code), INDEX idx_user (user_id)
) ENGINE=InnoDB COMMENT='优惠卡券 @ticket 06-个人中心';

CREATE TABLE IF NOT EXISTS my_service (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(64), url VARCHAR(256), perm VARCHAR(64), roles VARCHAR(128),
  enabled TINYINT DEFAULT 1,
  UNIQUE KEY uk_code (code)
) ENGINE=InnoDB COMMENT='我的服务 @ticket 06-个人中心';

-- ---------- 工单01 消息通知 ----------
CREATE TABLE IF NOT EXISTS notify_template (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(64) NOT NULL, channel VARCHAR(32),
  content TEXT, variables TEXT, enabled TINYINT DEFAULT 1,
  UNIQUE KEY uk_code (code)
) ENGINE=InnoDB COMMENT='通知模板 @ticket 01-消息通知';

CREATE TABLE IF NOT EXISTS notify_message (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tpl_code VARCHAR(64), receiver_type VARCHAR(32), receiver_id BIGINT,
  channel VARCHAR(32), content TEXT,
  status TINYINT DEFAULT 0 COMMENT '0待发 1成功 2失败 3重试中 @ticket 01',
  retry_count INT DEFAULT 0, dedupe_key VARCHAR(128),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_receiver (receiver_type, receiver_id), INDEX idx_dedupe (dedupe_key)
) ENGINE=InnoDB COMMENT='通知消息 @ticket 01-消息通知';

-- ---------- 初始化基础数据 ----------
INSERT IGNORE INTO pricing_rule (size, first_hours, first_fee, extra_hour_fee, daily_cap, deposit, priority, enabled) VALUES
  ('small', 1, 200, 100, 2000, 0, 10, 1),
  ('medium', 1, 300, 150, 3000, 0, 10, 1),
  ('large', 1, 500, 200, 4000, 0, 10, 1);

INSERT IGNORE INTO city (name, pinyin, initial, is_hot, sort, location) VALUES
  ('郑州', 'zhengzhou', 'Z', 1, 1, '34.7466,113.6254'),
  ('北京', 'beijing', 'B', 1, 2, '39.9042,116.4074'),
  ('上海', 'shanghai', 'S', 1, 3, '31.2304,121.4737');

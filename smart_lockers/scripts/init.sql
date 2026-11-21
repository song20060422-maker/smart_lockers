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

-- ============================================================
-- 补齐种子数据 @ticket 00-project：使 17 张表联调可跑（幂等 INSERT IGNORE）
-- 所有 seed 使用固定 ID/唯一键，重复执行不产生重复数据
-- ============================================================
INSERT IGNORE INTO guide (id, version, lang, content, is_current) VALUES
  (1, 'v1.0', 'zh-CN', '寄存流程：选择寄存点 → 选择柜门 → 支付下单 → 存入物品 → 取件时输入订单号。', 1);

INSERT IGNORE INTO user_profile (id, phone, nickname, avatar, role, status, token_ver) VALUES
  (10001, '13800138000', '演示用户', '', 'user', 1, 0);

INSERT IGNORE INTO locker_site (id, city_id, name, address, lat, lng, open_time, close_time, status) VALUES
  (1, 1, '郑州东站寄存点', '郑州市金水区郑州东站候车厅', 34.7521, 113.7762, '06:00', '23:00', 1),
  (2, 1, '二七广场寄存点', '郑州市二七区二七广场地铁站A口', 34.7466, 113.6254, '07:00', '22:00', 1);

INSERT IGNORE INTO locker_cabinet (id, site_id, code, size, total, available, status) VALUES
  (1, 1, 'ZDA-A', 'medium', 12, 8, 1),
  (2, 2, 'EQQ-A', 'small', 20, 15, 1);

INSERT IGNORE INTO locker_door (id, cabinet_id, code, size, status, version) VALUES
  (1, 1, 'ZDA-A-01', 'medium', 0, 0),
  (2, 1, 'ZDA-A-02', 'medium', 0, 0),
  (3, 2, 'EQQ-A-01', 'small', 0, 0);

INSERT IGNORE INTO deposit_order (id, order_no, user_id, door_id, site_id, start_time, end_time, amount, deposit, status, idempotency_key) VALUES
  (1, 'DEP202611200001', 10001, 1, 1, NOW(), DATE_ADD(NOW(), INTERVAL 24 HOUR), 300, 0, 1, 'seed-dep-001');

INSERT IGNORE INTO temp_open_auth (id, order_no, token, max_times, used_times, expire_at, used) VALUES
  (1, 'DEP202611200001', 'seed-open-token-001', 3, 0, DATE_ADD(NOW(), INTERVAL 2 HOUR), 0);

INSERT IGNORE INTO account (id, user_id, balance) VALUES
  (1, 10001, 10000);

INSERT IGNORE INTO balance_log (id, user_id, change_amount, balance_after, biz_type, biz_no) VALUES
  (1, 10001, 10000, 10000, 'recharge', 'SEED-RC-001');

INSERT IGNORE INTO recharge_order (id, charge_no, user_id, amount, status, trade_no) VALUES
  (1, 'RC202611200001', 10001, 10000, 2, 'SEED-TRADE-001');

INSERT IGNORE INTO coupon (id, code, user_id, discount, scope, expire_at, used) VALUES
  (1, 'NEWUSER5', 10001, 500, 'deposit', DATE_ADD(NOW(), INTERVAL 30 DAY), 0),
  (2, 'WELCOME2', 10001, 200, 'deposit', DATE_ADD(NOW(), INTERVAL 30 DAY), 0);

INSERT IGNORE INTO my_service (id, code, name, url, perm, roles, enabled) VALUES
  (1, 'my_orders', '我的订单', '/orders', 'order:list', 'user', 1),
  (2, 'my_coupons', '我的卡券', '/coupons', 'coupon:list', 'user', 1),
  (3, 'my_service', '我的服务', '/services', 'service:list', 'user', 1);

INSERT IGNORE INTO notify_template (id, code, channel, content, variables, enabled) VALUES
  (1, 'order_created', 'site', '您的寄存订单已创建，订单号：{{order_no}}', 'order_no', 1),
  (2, 'order_open', 'site', '取件授权已生成，请及时取件。', 'order_no', 1);

INSERT IGNORE INTO notify_message (id, tpl_code, receiver_type, receiver_id, channel, content, status, retry_count, dedupe_key) VALUES
  (1, 'order_created', 'user', 10001, 'site', '您的寄存订单已创建，订单号：DEP202611200001', 1, 0, 'seed-notify-001');


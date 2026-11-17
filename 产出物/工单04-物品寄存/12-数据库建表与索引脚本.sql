-- 工单04 物品寄存 建表与索引脚本 @ticket 04-物品寄存
CREATE TABLE IF NOT EXISTS locker_cabinet (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, site_id BIGINT NOT NULL,
  code VARCHAR(32), size VARCHAR(16), total INT DEFAULT 0, available INT DEFAULT 0, status TINYINT DEFAULT 1,
  INDEX idx_site (site_id)
) ENGINE=InnoDB COMMENT='柜组 @ticket 04-物品寄存';

CREATE TABLE IF NOT EXISTS locker_door (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, cabinet_id BIGINT NOT NULL,
  code VARCHAR(32), size VARCHAR(16),
  status TINYINT DEFAULT 0 COMMENT '0空闲 1占用 2故障 @ticket 04',
  version INT DEFAULT 0,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_cabinet (cabinet_id), INDEX idx_status (status)
) ENGINE=InnoDB COMMENT='柜门 @ticket 04-物品寄存';

CREATE TABLE IF NOT EXISTS pricing_rule (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, size VARCHAR(16) NOT NULL,
  first_hours INT DEFAULT 1, first_fee BIGINT DEFAULT 0 COMMENT '分 @ticket 04',
  extra_hour_fee BIGINT DEFAULT 0, daily_cap BIGINT DEFAULT 0,
  deposit BIGINT DEFAULT 0, priority INT DEFAULT 0, enabled TINYINT DEFAULT 1,
  INDEX idx_size (size)
) ENGINE=InnoDB COMMENT='收费规则 @ticket 04-物品寄存';

CREATE TABLE IF NOT EXISTS deposit_order (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, order_no VARCHAR(64) NOT NULL,
  user_id BIGINT NOT NULL, door_id BIGINT NOT NULL, site_id BIGINT NOT NULL,
  start_time DATETIME, end_time DATETIME, amount BIGINT DEFAULT 0, deposit BIGINT DEFAULT 0,
  status TINYINT COMMENT '1进行中 2已完成 3已取消 @ticket 04/05',
  idempotency_key VARCHAR(128),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_order_no (order_no), UNIQUE KEY uk_idem (idempotency_key),
  INDEX idx_user_status (user_id, status), INDEX idx_door (door_id)
) ENGINE=InnoDB COMMENT='寄存订单 @ticket 04/05';

INSERT IGNORE INTO pricing_rule (size,first_hours,first_fee,extra_hour_fee,daily_cap,deposit,priority,enabled) VALUES
  ('small',1,200,100,2000,0,10,1),('medium',1,300,150,3000,0,10,1),('large',1,500,200,4000,0,10,1);

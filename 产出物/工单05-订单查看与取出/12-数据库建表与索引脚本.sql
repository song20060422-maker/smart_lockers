-- 工单05 订单查看与取出 建表与索引脚本 @ticket 05-订单查看与取出
CREATE TABLE IF NOT EXISTS deposit_order (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, order_no VARCHAR(64) NOT NULL,
  user_id BIGINT NOT NULL, door_id BIGINT NOT NULL, site_id BIGINT NOT NULL,
  start_time DATETIME, end_time DATETIME, amount BIGINT DEFAULT 0, deposit BIGINT DEFAULT 0,
  status TINYINT COMMENT '1进行中 2已完成 3已取消 @ticket 05',
  idempotency_key VARCHAR(128),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_order_no (order_no), UNIQUE KEY uk_idem (idempotency_key),
  INDEX idx_user_status (user_id, status), INDEX idx_door (door_id)
) ENGINE=InnoDB COMMENT='寄存订单 @ticket 04/05';

CREATE TABLE IF NOT EXISTS temp_open_auth (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, order_no VARCHAR(64) NOT NULL, token VARCHAR(128) NOT NULL,
  max_times INT DEFAULT 1, used_times INT DEFAULT 0,
  expire_at DATETIME NOT NULL, used TINYINT DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_token (token), INDEX idx_order (order_no)
) ENGINE=InnoDB COMMENT='临时开柜授权 @ticket 05-订单查看与取出';

CREATE TABLE IF NOT EXISTS makeup_record (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, order_no VARCHAR(64) NOT NULL,
  amount BIGINT NOT NULL, status TINYINT, idempotency_key VARCHAR(128),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_idem (idempotency_key), INDEX idx_order (order_no)
) ENGINE=InnoDB COMMENT='超时补缴记录 @ticket 05-订单查看与取出';

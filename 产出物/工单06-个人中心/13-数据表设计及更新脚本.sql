-- 工单06 个人中心 数据表设计及更新脚本 @ticket 06-个人中心
CREATE TABLE IF NOT EXISTS user_profile (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, phone VARCHAR(32) NOT NULL,
  nickname VARCHAR(64), avatar VARCHAR(256), role VARCHAR(32),
  status TINYINT DEFAULT 1,
  token_ver INT DEFAULT 0 COMMENT '退出登录自增使旧token失效 @ticket 06',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_phone (phone)
) ENGINE=InnoDB COMMENT='用户资料 @ticket 06-个人中心';

CREATE TABLE IF NOT EXISTS coupon (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, code VARCHAR(64) NOT NULL, user_id BIGINT NOT NULL,
  discount BIGINT, scope VARCHAR(64), expire_at DATETIME NOT NULL, used TINYINT DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_code (code), INDEX idx_user (user_id)
) ENGINE=InnoDB COMMENT='优惠卡券 @ticket 06-个人中心';

CREATE TABLE IF NOT EXISTS my_service (
  id BIGINT PRIMARY KEY AUTO_INCREMENT, code VARCHAR(64) NOT NULL,
  name VARCHAR(64), url VARCHAR(256), perm VARCHAR(64), roles VARCHAR(128),
  enabled TINYINT DEFAULT 1,
  UNIQUE KEY uk_code (code)
) ENGINE=InnoDB COMMENT='我的服务 @ticket 06-个人中心';

-- 更新脚本示例：新增头像字段（幂等）
-- ALTER TABLE user_profile ADD COLUMN avatar VARCHAR(256) AFTER nickname;

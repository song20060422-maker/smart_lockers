-- 工单03 首页 建表与索引脚本 @ticket 03-首页
CREATE TABLE IF NOT EXISTS city (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL COMMENT '城市名 @ticket 03',
  pinyin VARCHAR(128), initial CHAR(2),
  is_hot TINYINT DEFAULT 0, sort INT DEFAULT 0, location VARCHAR(128),
  INDEX idx_initial (initial), INDEX idx_hot (is_hot)
) ENGINE=InnoDB COMMENT='城市 @ticket 03-首页';

CREATE TABLE IF NOT EXISTS locker_site (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  city_id BIGINT NOT NULL, name VARCHAR(128) NOT NULL, address VARCHAR(256),
  lat DECIMAL(10,7), lng DECIMAL(10,7),
  open_time VARCHAR(16), close_time VARCHAR(16), status TINYINT DEFAULT 1,
  INDEX idx_city (city_id), INDEX idx_lat_lng (lat, lng)
) ENGINE=InnoDB COMMENT='寄存点 @ticket 03-首页';

CREATE TABLE IF NOT EXISTS guide (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  version VARCHAR(32) NOT NULL, lang VARCHAR(16), content TEXT, is_current TINYINT DEFAULT 0,
  UNIQUE KEY uk_version (version)
) ENGINE=InnoDB COMMENT='寄存指南 @ticket 03-首页';

-- 初始化基础数据 @ticket 03
INSERT IGNORE INTO city (name,pinyin,initial,is_hot,sort,location) VALUES
  ('郑州','zhengzhou','Z',1,1,'34.7466,113.6254'),
  ('北京','beijing','B',1,2,'39.9042,116.4074'),
  ('上海','shanghai','S',1,3,'31.2304,121.4737');

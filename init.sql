-- 创建数据库
CREATE DATABASE test /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
  id int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  name varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户登录名，必须唯一，不可以重复。但是可以修改。',
  password varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  email varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '电子邮件',
  language varchar(20) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户语言',
  role int(11) NOT NULL COMMENT '角色',
  noticeEnable tinyint(4) NOT NULL COMMENT '是否提示',
  noticeLevel int(11) NOT NULL COMMENT '邮件提示级别',
  createTime int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (id),
  UNIQUE KEY id_UNIQUE (id),
  UNIQUE KEY name_UNIQUE (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_config (
  id bigint(20) NOT NULL AUTO_INCREMENT,
  op_log_expired int(11) NOT NULL DEFAULT '30' COMMENT '操作日志过期时间，单位: 天，默认30天',
  PRIMARY KEY (id),
  UNIQUE KEY op_log_expired_UNIQUE (op_log_expired)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 初始化系统配置
INSERT INTO system_config(op_log_expired) VALUES (30);
-- 提交
commit;
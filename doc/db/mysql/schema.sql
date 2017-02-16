CREATE DATABASE im DEFAULT CHARACTER SET utf8;
 
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(19) NOT NULL AUTO_INCREMENT,
  `uid` bigint(19) NOT NULL,
  `user_name` varchar(255) NOT NULL,
  `password`  varchar(255) NOT NULL,
  `gmt_create` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `gmt_modify` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `id` bigint(19) NOT NULL AUTO_INCREMENT,
  `gid` bigint(19) NOT NULL,
  `group_name` varchar(255) NOT NULL,
  `owner_id`   bigint(19) NOT NULL,
  `gmt_create` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `gmt_modify` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `user_group`;
CREATE TABLE `user_group` (
  `id` bigint(19) NOT NULL AUTO_INCREMENT,
  `gid` bigint(19) NOT NULL,
  `uid` varchar(255) NOT NULL,
  `gmt_create` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `gmt_modify` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `user_msg_id`;
CREATE TABLE `user_msg_id` (
  `id` bigint(19) NOT NULL AUTO_INCREMENT,
  `uid` bigint(19) NOT NULL,
  `current_msg_id` bigint(19) NOT NULL,
  `total_msg_id` bigint(19) NOT NULL,
  `gmt_create` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `gmt_modify` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
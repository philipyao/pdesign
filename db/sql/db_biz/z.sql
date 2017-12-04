create database IF NOT EXISTS db_hgame;

use db_hgame;
create table IF NOT EXISTS t_user
(
  `charid` bigint(20) unsigned NOT NULL,
  `accid` int(10) unsigned NOT NULL,
  `ulevel` int(10) unsigned NOT NULL default '1',
  `uexp` int(10) unsigned NOT NULL default '0',
  `main` blob,
  `rare` blob,
  PRIMARY KEY  (`charid`),
  KEY (`accid`),
  UNIQUE KEY `update_or_insert` (`charid`)  
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame;
create table IF NOT EXISTS t_freeid
(
  `typeid` int(10) unsigned NOT NULL,
  `freeid` bigint(20) unsigned NOT NULL default '0',
  PRIMARY KEY  (`typeid`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame;
create table IF NOT EXISTS t_sys_mail
(
    `seqno`     int(10) unsigned NOT NULL AUTO_INCREMENT,
    `createtm`  int(10) unsigned NOT NULL,
    `starttm`   int(10) unsigned NOT NULL,
    `expiretm`  int(10) unsigned NOT NULL,
    `param`     blob,
    PRIMARY KEY (`seqno`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame;
create table IF NOT EXISTS t_guild
(
    `id`        int(10) unsigned NOT NULL,
    `name`      varchar(64) NOT NULL default '',
    `leader`    int(10) unsigned NOT NULL,
    `leadercharid` bigint(20) unsigned NOT NULL,
    `base`      blob,
    `members`   blob,
    PRIMARY KEY (`id`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame;
create table IF NOT EXISTS t_param
(
    `id`        int(10) unsigned NOT NULL COMMENT '参数对应的id',
    `name`      varchar(64) NOT NULL default '',
    `val`       mediumblob COMMENT "参数的值",
    `updatetm`  int(10) unsigned NOT NULL,
    `ver`       int(10) unsigned NOT NULL default '0' COMMENT '参数版本',
    PRIMARY KEY (`id`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame;
create table IF NOT EXISTS t_scene
(
    `id`        int(10) unsigned NOT NULL,
    `npcs`      blob,
    PRIMARY KEY (`id`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame;
create table IF NOT EXISTS t_user_profile
(
    `charid`    bigint(20) unsigned NOT NULL,
    `info`      blob,
    PRIMARY KEY (`charid`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame;
create table IF NOT EXISTS t_friendship
(
    `charidstr` varchar(64) NOT NULL,
    `charid1`    bigint(20) unsigned NOT NULL,
    `accid1`     int(10) unsigned NOT NULL,
    `charid2`    bigint(20) unsigned NOT NULL,
    `accid2`     int(10) unsigned NOT NULL,
    `addtime`     int(10) unsigned NOT NULL,
    PRIMARY KEY (`charidstr`),
    UNIQUE KEY `UNIQ_IDSTR`(`charid1`, `charid2`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame;
create table IF NOT EXISTS t_top_report
(
    `id`        bigint(20) unsigned NOT NULL COMMENT '战报唯一id',
    `data`      blob COMMENT '战报pb数据',
    PRIMARY KEY (`id`)
) engine = InnoDB DEFAULT CHARSET=utf8;


use db_hgame;
create table IF NOT EXISTS t_activity
(
    `id`        bigint(20) unsigned NOT NULL COMMENT '活动唯一id',
    `kind`      int(10) unsigned default '0' COMMENT '活动类型',
    `tid`       int(10) unsigned default '0' COMMENT '活动配置tid',
    `createtm`  int(10) unsigned NOT NULL,
    `starttm`   int(10) unsigned NOT NULL,
    `expiretm`  int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`)
) engine = InnoDB DEFAULT CHARSET=utf8;



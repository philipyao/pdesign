create database IF NOT EXISTS db_hgame_global;

use db_hgame_global;
create table IF NOT EXISTS t_account
(
  `accid` int(10) unsigned NOT NULL,
  `chanid` int(10) unsigned NOT NULL,
  `openid` varchar(128) NOT NULL,
  `flags` int(10) unsigned NOT NULL,
  `bitmap` int(10) unsigned NOT NULL,
  `expire_time` int(10) unsigned NOT NULL,
  `active_time` int(10) unsigned NOT NULL,
  `create_time` int(10) unsigned NOT NULL,
  `main` blob,
  PRIMARY KEY  (`accid`),
  UNIQUE KEY `openid` (`openid`)
)engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame_global;
create table IF NOT EXISTS t_gift
(
  `id` int(10) unsigned NOT NULL,
  `typ` int(10) unsigned NOT NULL default '0',
  `tm_start` int(10) unsigned NOT NULL default '0',
  `tm_end` int(10) unsigned NOT NULL default '0',
  `tm_gen` int(10) unsigned NOT NULL default '0',
  `muid` int(10) unsigned NOT NULL default '0',
  `num` int(10) unsigned NOT NULL default '0',
  `main` blob,
  `keys` blob,
  PRIMARY KEY  (`id`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame_global;
create table IF NOT EXISTS t_exchange
(
  `uid` bigint(20) unsigned NOT NULL,
  `cdkey` varchar(15) NOT NULL default '',
  `gift` int(10) unsigned NOT NULL default '0',
  `muid` int(10) unsigned NOT NULL default '0',
  `zone` int(10) unsigned NOT NULL default '0',
  `accid` int(10) unsigned NOT NULL default '0',
  `tm_use` int(10) unsigned NOT NULL default '0',
  `name` varchar(64) NOT NULL default '',
  PRIMARY KEY  (`uid`, `cdkey`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame_global;
create table IF NOT EXISTS t_keyvalue
(
  `id` bigint(20) unsigned NOT NULL,
  `typ` int(10) unsigned NOT NULL default '0',
  `ver` int(10) unsigned NOT NULL default '0',
  `val` blob,
  PRIMARY KEY  (`id`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame_global;
create table IF NOT EXISTS t_pay
(
  `pid` bigint(20) unsigned NOT NULL  default '0',
  `orderid` varchar(64) NOT NULL,
  `tm` int(10) unsigned NOT NULL default '0',
  `state` int(10) unsigned NOT NULL default '0',
  `zone` int(10) unsigned NOT NULL default '0',
  `accid` int(10) unsigned NOT NULL default '0',
  `charid` bigint(20) unsigned NOT NULL default '0',
  `good` int(10) unsigned NOT NULL default '0',
  `count` int(10) unsigned NOT NULL default '0',
  `price` int(10) unsigned NOT NULL default '0',
  `csdkid` varchar(32) NOT NULL,
  `sdkid` varchar(32) NOT NULL,
  `chanid` int(10) unsigned NOT NULL default '0',
  `paytype` varchar(32) NOT NULL,
  `paytime` varchar(32) NOT NULL,
  `payresult` varchar(32) NOT NULL,
  `mu` int(10) unsigned NOT NULL default '0',
  `tmpay` int(10) unsigned NOT NULL default '0',
  `fake` int(10) unsigned NOT NULL default '0',
  `ext` blob,
  PRIMARY KEY  (`pid`)
) engine = InnoDB DEFAULT CHARSET=utf8;


use db_hgame_global;
create table IF NOT EXISTS t_fanli
(
  `accid`   int(10) unsigned NOT NULL,
  `rmb`     int(10) unsigned NOT NULL default '0',
  `gotzone` int(10) unsigned NOT NULL default '0' COMMENT '领取的zone, 0表示未领取',
  `tm`      int(10) unsigned NOT NULL default '0' COMMENT '领取时间',
  `mux`     int(10) unsigned NOT NULL default '0' COMMENT '互斥操作计数',
  PRIMARY KEY  (`accid`)
) engine = InnoDB DEFAULT CHARSET=utf8;

/* 充值校验订单数据 */
use db_hgame_global;
create table IF NOT EXISTS t_billing
(
  `billid` bigint(20) unsigned NOT NULL  default 0,
  `accid`  int(10) unsigned NOT NULL default 0,
  `chanid` int(10) unsigned NOT NULL default 0,
  `state`    varchar(32) NOT NULL,
  `paytype`  varchar(64) NOT NULL,
  `serverid` varchar(64) NOT NULL,
  `game_orderid` varchar(64) NOT NULL,
  `gooidsid`  varchar(64) NOT NULL,
  `price`  int(10) unsigned NOT NULL default 0,
  `count`  int(10) unsigned NOT NULL default 0,
  `tmlaunch` int(10) unsigned NOT NULL default 0,
  `tmcomplete` int(10) unsigned NOT NULL default 0,
  `result`      varchar(64) NOT NULL,
  `push_result` varchar(64) NOT NULL,
  `push_times`  int(10) unsigned NOT NULL default 0,
  `fake` int(10) unsigned NOT NULL default 0,
  `param` text NOT NULL,
  `reply` text NOT NULL,
  `ext` varchar(256) NOT NULL,
  `mux` int(10) unsigned NOT NULL default 0,
  PRIMARY KEY  (`billid`),
  KEY (`accid`, `tmlaunch`)
) engine = InnoDB DEFAULT CHARSET=utf8;

/* 苹果票据记录 */
use db_hgame_global;
create table IF NOT EXISTS t_iap_receipt
(
  `receipt` varchar(64) NOT NULL,
  `status`  int(10)  NOT NULL default 0,
  `billid`  bigint(20) unsigned NOT NULL  default 0, 

  PRIMARY KEY  (`receipt`)
) engine = InnoDB DEFAULT CHARSET=utf8;

use db_hgame_global;
create table IF NOT EXISTS t_invite_code
(
  `code`    varchar(11)  not null,
  `zoneid`  int(10)     unsigned not null,
  `accid`   int(10)     unsigned not null,
  PRIMARY KEY  (`code`),
  UNIQUE INDEX `code` (`code`)
)engine = InnoDB DEFAULT CHARSET=utf8;

/* 瀚趣账号 */
use db_hgame_global;
create table IF NOT EXISTS t_hanqu_account
(
  `huin`      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `hacc`      varchar(64) NOT NULL UNIQUE,
  `hpwd`      varchar(64) DEFAULT NULL,
  `tmreg`     int(10) unsigned NOT NULL, 
  `tmlast`    int(10) unsigned NOT NULL, 
  `mux`       int(10) unsigned NOT NULL default 0,

  PRIMARY KEY  (`huin`),
  KEY (`hacc`)
)engine = InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8;


/* 瀚趣绑定账号 */
use db_hgame_global;
create table IF NOT EXISTS t_hanqu_bind_acc
(
  `uid`     varchar(64) NOT NULL,
  `typ`     varchar(16) NOT NULL,
  `huin`    bigint(20) unsigned NOT NULL,
  `tmreg`   int(10) unsigned NOT NULL,
  `mux`     int(10) unsigned NOT NULL default 0,

  PRIMARY KEY  (`uid`),
  KEY (`huin`)
)engine = InnoDB  DEFAULT CHARSET=utf8;

/* 客户端版本md5码 */
use db_hgame_global;
create table if not exists t_cli_md5
(
    `ver`   varchar(64) not null,
    `sub`   varchar(32) not null,
    `md5`   varchar(128) not null,
    PRIMARY KEY (`ver`, `sub`)
)engine = InnoDB  DEFAULT CHARSET=utf8;

/* 活动display配置 */
use db_hgame_global;
create table IF NOT EXISTS t_act_display_cfg
(
  `tid`             int(10) unsigned NOT NULL,
  `description`     varchar(256) NOT NULL COMMENT '备注',
  `category`        int(10) unsigned NOT NULL COMMENT '分类',
  `weight`          int(10) unsigned NOT NULL COMMENT '显示排序权重',
  `widget_type`     int(10) unsigned NOT NULL COMMENT '空间类型',
  `name`            varchar(256) NOT NULL COMMENT '名称',
  `title`           varchar(256) NOT NULL COMMENT '标题',
  `icon_image`      varchar(256) NOT NULL COMMENT '图标image',
  `back_image`      varchar(256) NOT NULL COMMENT '背景image',
  `custom`          blob COMMENT '自定义显示',

  PRIMARY KEY  (`tid`),
  KEY (`tid`)
)engine = InnoDB  DEFAULT CHARSET=utf8;

/* func 类型活动配置 */
use db_hgame_global;
create table IF NOT EXISTS t_func_act_cfg
(
  `tid`             int(10) unsigned NOT NULL,
  `acttp`           int(10) unsigned NOT NULL COMMENT '活动类型',
  `acttp_para`      int(10) unsigned NOT NULL COMMENT '类型参数',
  `entries`         blob COMMENT '多条目，包括条件值和奖励',
  `name`            varchar(256) NOT NULL COMMENT '活动名称',
  `display`         int(10) unsigned NOT NULL COMMENT '显示id',

  PRIMARY KEY  (`tid`),
  KEY (`tid`)
)engine = InnoDB  DEFAULT CHARSET=utf8;

/* goal 类型活动配置 */
use db_hgame_global;
create table IF NOT EXISTS t_goal_act_cfg
(
  `tid`             int(10) unsigned NOT NULL,
  `target`          int(10) unsigned NOT NULL COMMENT '条件类型',
  `target_para`     int(10) unsigned NOT NULL COMMENT '条件参数',
  `entries`         blob COMMENT '多条目，包括条件值和奖励',
  `name`            varchar(256) NOT NULL COMMENT '活动名称',
  `cost`            int(10) unsigned NOT NULL COMMENT '领取奖励的消耗钻石',
  `mailid`          int(10) unsigned NOT NULL COMMENT '完成后的邮件通知',
  `display`         int(10) unsigned NOT NULL COMMENT '显示id',

  PRIMARY KEY  (`tid`),
  KEY (`tid`)
)engine = InnoDB  DEFAULT CHARSET=utf8;

/* 静态(全服共享)开启的活动 */
use db_hgame_global;
create table IF NOT EXISTS t_static_activity
(
 `id`        bigint(20) unsigned NOT NULL COMMENT '活动唯一id',
 `kind`      int(10) unsigned default '0' COMMENT '活动类型',
 `tid`       int(10) unsigned default '0' COMMENT '活动配置tid',
 `createtm`  int(10) unsigned NOT NULL,
 `opentm`   int(10) unsigned NOT NULL COMMENT '绝对时间开启：开启时间',
 `closetm`  int(10) unsigned NOT NULL COMMENT '绝对时间开启：关闭时间',
 `nexttm`   int(10) unsigned NOT NULL COMMENT '绝对时间开启：下次开启时间',
 `offset_day` int(10) unsigned NOT NULL COMMENT '相对时间开启: 开启时间相对开服时间的offset',
 `lasting` int(10) unsigned NOT NULL COMMENT '相对时间开启: 持续时间',
 `effrange`    blob COMMENT '起效范围',
 PRIMARY KEY (`id`)
)engine = InnoDB  DEFAULT CHARSET=utf8;


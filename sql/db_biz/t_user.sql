create database IF NOT EXISTS db_biz;

use db_biz;
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


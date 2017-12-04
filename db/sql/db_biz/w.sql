create database IF NOT EXISTS db_hgame_world;

use db_hgame_world;
create table IF NOT EXISTS t_name
(
  `name` varchar(64) NOT NULL,
  `nid`  bigint(20) unsigned NOT NULL default '0' UNIQUE,
  `accid`  int(10) unsigned NOT NULL default '0',
  `zoneid` int(10) unsigned NOT NULL default '0',
  `ntyp` int(10) unsigned NOT NULL default '0',
  `used` int(10) unsigned NOT NULL default '0',
  `flag` int(10) unsigned NOT NULL default '0',
  `tm` int(10) unsigned NOT NULL default '0',
  PRIMARY KEY  (`name`, `ntyp`)
) engine = InnoDB DEFAULT CHARSET=utf8;

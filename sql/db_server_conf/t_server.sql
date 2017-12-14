create database IF NOT EXISTS db_server_conf;

use db_server_conf;
create table IF NOT EXISTS t_server
(
  `typename` varchar(32) NOT NULL,
  `typeid` int(10) unsigned NOT NULL,
  `params` varchar(512) default '',
  PRIMARY KEY  (`typename`)
) engine = InnoDB DEFAULT CHARSET=utf8;

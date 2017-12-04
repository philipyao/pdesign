create database IF NOT EXISTS db_biz;

use db_biz;
create table IF NOT EXISTS t_db_route
(
  `tablestr` varchar(32) NOT NULL,
  `dbname` varchar(32) NOT NULL,
  `dbaddr` varchar(32) NOT NULL,
  PRIMARY KEY  (`tablestr`)
) engine = InnoDB DEFAULT CHARSET=utf8;

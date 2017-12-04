create database IF NOT EXISTS db_biz;

use db_biz;
create table IF NOT EXISTS t_user_profile
(
    `charid`    bigint(20) unsigned NOT NULL,
    `info`      blob,
    PRIMARY KEY (`charid`)
) engine = InnoDB DEFAULT CHARSET=utf8;

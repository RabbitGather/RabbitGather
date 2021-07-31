show tables;


drop table if exists `user`;

CREATE TABLE `user`
(
    `id` int unsigned primary key auto_increment,
    `name` varchar(24)  not null unique,
    `password` char(60) not null ,
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `api_permission_bitmask` int unsigned not null
);





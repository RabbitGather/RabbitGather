show tables;


drop table if exists `user`;

CREATE TABLE `user`
(
    `id` int unsigned primary key auto_increment,
    `name` varchar(24)  not null unique,
    `password` varchar(128) not null ,
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `api_permission_bitmask` int unsigned not null
);

insert into `user` (name,password,api_permission_bitmask)
values
('Name1','Password1',1)
;


# CREATE TABLE `api_permission_bitmask`
# (
#     `id` int unsigned primary key auto_increment,
#     `name` varchar(24)  not null unique,
#     `bitmask` int unsigned not null,
#     `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ,
#     `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
# );




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



CREATE TABLE `user_article_setting`
(
    `user` int unsigned primary key,
    `setting` JSON  not null,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    foreign key (`user`) references user(`id`)
);

show tables;

insert into `user_article_setting` (user, setting)
values
(1,'{"max_radius":100,"min_radius":1}')
on duplicate  key update user = 1 , setting = '{"max_radius":100,"min_radius":1}';

SELECT * FROM `user_article_setting` WHERE setting = CAST('{"max_radius":100,"min_radius":1}' as JSON);

select setting from `user_article_setting` where setting->'max_radius' = 100;
SELECT * FROM `user_article_setting` WHERE setting->'$.max_radius' = 100;
SELECT * FROM `user_article_setting` WHERE JSON_CONTAINS(setting, '100', '$.max_radius');
UPDATE `user_article_setting` SET setting = JSON_INSERT(setting, '$.area', 'china') WHERE user = 1;
UPDATE `user_article_setting` SET setting = JSON_SET(setting, '$.area1', 'china', '$.max_radius', '200') WHERE user = 1;
UPDATE `user_article_setting` SET setting = JSON_REPLACE(setting, '$.area1', 'chinasss', '$.max_radius', '400') WHERE user = 1;
UPDATE `user_article_setting` SET setting = JSON_REMOVE(setting, '$.area') WHERE user = 1;


SELECT * FROM `user_article_setting`;


select setting->'$.max_radius' from `user_article_setting`;

select setting->'$.max_radius' from `user_article_setting`;


select * from `user`;
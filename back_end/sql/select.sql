select * from user;
select * from user;
show tables;
select password from user where id = '' limit 1;
select id from user where name = '' limit 1;
insert into user ( name, password, api_permission_bitmask) value ('','','');


insert into `article` (title,content)
    value('A test title','some content')
;

insert into `article_tag` (article_id,tag_name)
    value(1,'DELETED');

select * from `article`;
select * from `article_tag`;



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

select  * from `user`;
select  * from `user_information`;
insert into `user_information`(`user`,`email`)
    value (?,?);
insert into `user_information`(`user`,`phone`)
    value (?,?);

select `user` from `user_information`where email ='a.meowalien@gmail.com';
select `user` from `user_information`where phone =?;

select a.title , a.content , b.coords
from `article` as a left join `article_details` as b on a.id = b.article
where a.id = 1;


select ST_AsText(coords) from `article_details`;
select coords from `article_details`;
# POINT(25.040056717110396 121.51187490970621)

SELECT ST_AsBinary(coords) FROM `article_details`

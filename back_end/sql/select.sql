select *
from user;
select *
from user;
show tables;
select password
from user
where id = ''
limit 1;
select id
from user
where name = ''
limit 1;
insert into user (name, password, api_permission_bitmask) value ('', '', '');


insert into `article` (title, content)
    value ('A test title', 'some content')
;

insert into `article_tag` (article_id, tag_name)
    value (1, 'DELETED');

select *
from `article`;
select *
from `article_tag`;



show tables;

insert into `user_article_setting` (user, setting)
values (1, '{"max_radius":100,"min_radius":1}')
on duplicate key update user    = 1,
                        setting = '{"max_radius":100,"min_radius":1}';

SELECT *
FROM `user_article_setting`
WHERE setting = CAST('{"max_radius":100,"min_radius":1}' as JSON);

select setting
from `user_article_setting`
where setting -> 'max_radius' = 100;
SELECT *
FROM `user_article_setting`
WHERE setting -> '$.max_radius' = 100;
SELECT *
FROM `user_article_setting`
WHERE JSON_CONTAINS(setting, '100', '$.max_radius');
UPDATE `user_article_setting`
SET setting = JSON_INSERT(setting, '$.area', 'china')
WHERE user = 1;
UPDATE `user_article_setting`
SET setting = JSON_SET(setting, '$.area1', 'china', '$.max_radius', '200')
WHERE user = 1;
UPDATE `user_article_setting`
SET setting = JSON_REPLACE(setting, '$.area1', 'chinasss', '$.max_radius', '400')
WHERE user = 1;
UPDATE `user_article_setting`
SET setting = JSON_REMOVE(setting, '$.area')
WHERE user = 1;


SELECT *
FROM `user_article_setting`;


select setting -> '$.max_radius'
from `user_article_setting`;

select setting -> '$.max_radius'
from `user_article_setting`;


select *
from `user`;

select *
from `user`;
select *
from `user_information`;
insert into `user_information`(`user`, `email`)
    value (?, ?);
insert into `user_information`(`user`, `phone`)
    value (?, ?);

select `user`
from `user_information`
where email = 'a.meowalien@gmail.com';
select `user`
from `user_information`
where phone = ?;

select a.title, a.content, b.coords
from `article` as a
         left join `article_details` as b on a.id = b.article
where a.id = 1;


select ST_AsText(coords)
from `article_details`;
select coords
from `article_details`;
# POINT(25.040056717110396 121.51187490970621)

SELECT ST_AsBinary(coords)
FROM `article_details`;


BEGIN;
insert into `article` (title, content)
    value ('rjfiondfsdf', 'sdasdasjdlkasd')
;
insert into `article_details` (article, coords)
    value (LAST_INSERT_ID(), Point(25.040056717110396, 121.51187490970621))
;
COMMIT;


BEGIN;
insert into `article` (title, content)
    value (?, ?)
;
insert into `article_details` (article, coords)
    value (LAST_INSERT_ID(), Point(?, ?))
;
COMMIT;

select a.title, a.content, b.coords
from `article` as a
         left join `article_details` as b
         left join `article_tag` as c
                   on a.id = b.article
                   on b.article = c.article_id and c.tag_id != 1
where a.id = 1;

select a.title, a.content
from `article` as a
         left join `article_details` as b on a.id = b.article
         left join `article_tag` as c on a.id = c.article_id and c.tag_name = 1
where c.tag_name is null
  and a.id = 1;



insert into `article_tag` (article_id, tag_name, tag_type)
    value (1, 2, 1);

delete
from article_tag
where article_id = 1
  and tag_name = 1;

select *
from article;

update article set title = ?, content = ? where id = ?;

select * from article_user_setting;

begin ;
set @max_radius = 0;
set @min_radius = 0;


# begin ;
# set @user_id = 0;
# set @key = 0;
# set @value = 0;
# '{"max_radius":400}'
# "$.max_radius"
# 400
insert into `article_user_setting` (user, setting)
value
( ?,?)
on duplicate key update setting = JSON_SET(setting ,?,?);
# commit;

insert into `article_user_setting` (user, setting)
    value
    (1,'{"min_radius":8}')
on duplicate key update setting = JSON_SET(setting ,'$.min_radius',8);


select setting from article_user_setting;

delete from article_user_setting where user = 2;
insert into `article_user_setting` (user, setting)
    value
    (?,?)
on duplicate key update setting = JSON_SET(setting ,?,?);





insert into user ( name, password,randomSalt, api_permission_bitmask) value (?,?,?,?);
select password,randomSalt from user where id = ? limit 1;


select * from `user`;


insert into `article` (title, content)
    value (?, ?)
;
select  * from article;

insert into `article_details` (article, coords)
    value (?, Point(?, ?))
;


select title, content,UNIX_TIMESTAMP(update_time) from article where id = 13;

select a.title , a.content , ST_AsBinary( b.coords)
from `article` as a
         left join `article_details` as b on a.id = b.article
         left join  `article_tag` as c
                    on a.id = c.article_id and c.tag_name = 1
where c.tag_name is null and  a.id = ?;


select * from `article` left join `article_details`
    on `article`.id = `article_details`.article
where  `article`.id = 13;

insert into `article_tag` (article_id,tag_id) value(?,?);

select * from article_tag;

select * from article;


select title, content,UNIX_TIMESTAMP(update_time) from article as a left join article_tag t on a.id = t.article_id where id = 2 and t.tag_id != 1;


select a.title , a.content , ST_AsBinary( b.coords)
from `article` as a
         left join `article_details` as b on a.id = b.article
         left join  `article_tag` as c
                    on a.id = c.article_id and c.tag_id = 1
where c.tag_id is null and  a.id = ?;


update article left join article_tag a on article.id = a.article_id set title = '?', content = '?' where id = 2 and a.tag_id != 1;

select * from article;

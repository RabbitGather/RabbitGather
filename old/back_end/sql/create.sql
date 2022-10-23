CREATE TABLE `user`
(
    `id`                     int unsigned primary key auto_increment,
    `name`                   varchar(24) not null unique,
    `password`               char(60)    not null,
    `randomSalt`             char(24)    not null,
    `create_time`            timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`            timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `api_permission_bitmask` int unsigned not null
);


insert into `user` (name, password, randomSalt, api_permission_bitmask)
    value ('test name','testpassword','testsalt',8)
;
CREATE TABLE `user_information`
(
    `user`  int unsigned primary key,
    `email` varchar(254) unique,
    `phone` varchar(30) unique,
    foreign key (`user`) references user (`id`)
);
#使用者存取文章的權限設定
CREATE TABLE `article_user_setting`
(
    `user`        int unsigned primary key,
    `setting`     JSON not null,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    foreign key (`user`) references user (`id`)
);


CREATE TABLE `article`
(
    `id`          int unsigned primary key auto_increment,
    `title`       varchar(48) not null,
    `content`     mediumtext  not null,
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)
;


insert into `article` (title, content)
    value ('A test title', 'A test content')
;


CREATE TABLE `article_details`
(
    `article` int unsigned primary key,
    `coords`  point not null,
    foreign key (`article`) references article (`id`)
) ENGINE = MyISAM
;



insert into `article_details` (article, coords)
    value (1, Point(25.040056717110396, 121.51187490970621))
;

create table `tag_type`
(
    `id`   int unsigned primary key auto_increment,
    `name` char(24) not null unique
);

insert into `tag_type` (name)
    value ('SYSTEM_TAG')
;

create table `tags`
(
    `id`   int unsigned primary key auto_increment,
    `name` char(24) not null unique,
    `type` int unsigned  null unique,
        foreign key (`type`) references `tag_type` (`id`)

);

insert into `tags` (name,type)
values ('DELETE',1)
;

create table `article_tag`
(
    `tag_id`     int unsigned not null ,
    `article_id` int unsigned not null,
    primary key (`article_id`,`tag_id`),
    foreign key (`article_id`) references `article` (`id`),
    foreign key (`tag_id`) references `tags` (`id`)
);


insert into `article_tag` (article_id, tag_id)
    value (1, 1);



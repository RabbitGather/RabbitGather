show tables;

CREATE TABLE `user`
(
    `id` int unsigned primary key auto_increment,
    `name` varchar(24)  not null unique,
    `password` char(60) not null ,
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `api_permission_bitmask` int unsigned not null
);

CREATE TABLE `user_information`
(
    `user` int unsigned primary key,
    `email` varchar (254) unique ,
    `phone` varchar (30) unique ,
    foreign key (`user`) references user(`id`)
);

# 使用者存取文章的權限設定
CREATE TABLE `article_user_setting`
(
    `user` int unsigned primary key,
    `setting` JSON  not null,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    foreign key (`user`) references user(`id`)
);


CREATE TABLE `article`
(
    `id` int unsigned primary key auto_increment,
    `title` varchar(48)  not null ,
    `content` mediumtext not null ,
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)
;
insert into `article` (title,content)
    value('A test title', 'A test content')
;


CREATE TABLE `article_details`
(
    `article` int unsigned primary key,
    `coords` point not null,
    foreign key (`article`) references article(`id`)
)ENGINE=MyISAM
;

insert into `article_details` (article,coords)
    value(1, Point( 25.040056717110396,121.51187490970621))
;

create table `tag_type`
(
    `id` int unsigned primary key auto_increment,
    `type`  char(24) not null unique
);

insert into `tag_type` (type)
    value('SYSTEM_TAG')
;

create table `tag_name`
(
    `id` int unsigned primary key auto_increment,
    `name`  char(24) not null unique
);

insert into `tag_name` (name)
    value('DELETE')
;

create table `article_tag`
(
    `tag_id` int unsigned primary key auto_increment,
    `article_id`  int unsigned not null ,
    `tag_name` int unsigned not null,
    `tag_type`int unsigned not null,
    unique (`tag_name` , `tag_type`),
    unique (`article_id` , `tag_id`),
    foreign key (`article_id`) references `article`(`id`),
    foreign key (`tag_name`) references `tag_name`(`id`),
    foreign key (`tag_type`) references `tag_type`(`id`)
);


insert into `article_tag` (article_id,tag_name,tag_type)
    value(1,1,1);

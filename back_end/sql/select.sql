select * from user;
select * from user;
show tables;
select password from user where id = '' limit 1;
select id from user where name = '' limit 1;
insert into user ( name, password, api_permission_bitmask) value ('','','');
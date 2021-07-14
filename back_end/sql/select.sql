select * from user;

select password from user where id = '' limit 1;

insert into user ( name, password, api_permission_bitmask) value ('','','');
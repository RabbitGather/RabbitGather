select id,api_permission_bitmask from user where name = '' limit 1;

select password from user where id = '' limit 1;

insert into user ( name, password, api_permission_bitmask) value ('','','');
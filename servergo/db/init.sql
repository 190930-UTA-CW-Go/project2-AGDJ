create table users(
    id serial primary key,
    username varchar NOT NULL UNIQUE,
    password varchar NOT NULL
);

create table ips(
    ip varchar
);

create table installed(
    appname varchar
);

insert into users (username, password) values ('godfrey', 'hello');
insert into ips values ('52.176.60.129');
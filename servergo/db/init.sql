create table users(
    id serial primary key,
    username varchar NOT NULL UNIQUE,
    password varchar NOT NULL
);

create table ips(
    ip varchar
);

create table apps(
    name varchar,
    descrip varchar
)


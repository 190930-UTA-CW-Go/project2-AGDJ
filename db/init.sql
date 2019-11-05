create table users(
    id serial primary key,
    username varchar NOT NULL UNIQUE,
    password varchar NOT NULL
);

create table ipaddr(
    ip varchar NOT NULL,
    hostname varchar NOT NULL,
    username varchar REFERENCES users(username) NOT NULL
);

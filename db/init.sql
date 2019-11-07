create table users(
    id serial primary key,
    username varchar NOT NULL UNIQUE,
    password varchar NOT NULL
);

create sequence numcons;
select setval('numcons', 99);

create table running(
    name text PRIMARY KEY CHECK (name ~ '^container[0-9]+$') DEFAULT 'container' || nextval('numcons'),
    port int NOT NULL,
    username varchar REFERENCES users(username) NOT NULL
);

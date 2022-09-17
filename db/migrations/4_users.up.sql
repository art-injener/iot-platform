create table if not exists users
(
    id           serial  not null
        constraint users_id primary key,
    created_at   timestamp default CURRENT_TIMESTAMP,
    login        varchar not null,
    email        varchar,
    passwordHash varchar,
    salt         varchar
);

create index user_login_idx on users (login);
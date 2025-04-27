-- +goose Up
create type user_role as enum ('admin', 'user');
create table users (
  id serial primary key,
  name text not null,
  email text not null,
  password text not null,
  role user_role not null default 'user'
);

-- +goose Down
drop table users;
drop type user_role;

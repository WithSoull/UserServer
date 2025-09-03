-- +goose Up
drop type user_role;

-- +goose Down
create type user_role as enum ('ADMIN', 'USER');

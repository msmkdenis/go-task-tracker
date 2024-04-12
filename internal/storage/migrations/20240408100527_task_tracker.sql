-- +goose Up
create table if not exists scheduler (
    id integer primary key,
    date text not null,
    title text not null,
    comment text not null,
    repeat char(128) not null
);
create index if not exists scheduler_date_idx on scheduler (date);

-- +goose Down
drop index if exists scheduler_date_idx;
drop table if exists scheduler;



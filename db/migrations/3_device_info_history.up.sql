create table if not exists device_info_history
(
    id     serial not null constraint device_info_history_id primary key,
    device_id     varchar not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP,
    device jsonb
);

create index device_history_id_idx on device_info_history(device_id);
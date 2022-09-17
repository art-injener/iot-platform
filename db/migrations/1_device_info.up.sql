create table if not exists device_info
(
    id          serial  not null
        constraint device_info_id primary key,
    device_id   varchar not null
        constraint device_id unique,
    created_at  timestamp default CURRENT_TIMESTAMP,
    updated_at  timestamp default CURRENT_TIMESTAMP,
    device      jsonb
);

create index device_id_idx on device_info(device_id);
create table if not exists device_settings
(
    id         serial not null constraint device_settings_id primary key,
    device_id  varchar not null constraint device_settings_device_id unique,
    version    varchar,
    timezone   int       default 220,
    wui        int       default 10,
    gpst       int       default 3,
    created_at timestamp default current_timestamp
);

comment on column device_settings.timezone is '220 - Москва. [-14, +14] - смещение от UTC';
comment on column device_settings.wui is 'Временной интервал выхода на связь [10 - 10080]. Время в минутах';
comment on column device_settings.gpst is 'Максимальное время поиска спутников в минутахх';

create index device_settings_device_id_idx on device_settings (device_id);
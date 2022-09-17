create table if not exists user_devices
(
    user_id        integer not null
        constraint fk_users
            references users,
    device_info_id integer not null
        constraint fk_device_info
            references device_info,
    constraint user_device_info_pkey
        primary key (user_id, device_info_id)
);
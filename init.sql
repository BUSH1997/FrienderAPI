create table categories(
                           id serial primary key,
                           name varchar(256)
);

create table users(
                      id serial primary key,
                      uid int
);

create table events(
                       id serial primary key,
                       uid varchar(256) UNIQUE,
                       title varchar(256),
                       description text,
                       images text,
                       avatar_url text,
                       avatar_vk_id text,
                       starts_at bigserial,
                       time_created bigserial,
                       time_updated bigserial,
                       geo varchar(256),
                       category_id int references categories(id),
                       count_members int,
                       is_public bool,
                       is_private bool,
                       owner_id int references users(id),
                       is_deleted bool,
                       photos text,
                       members_limit int,
                       source varchar(256)
);



create table event_sharings(
                       id serial primary key,
                       event_id int references events(id),
                       user_id int references  users(id)
);


create table syncer(
    id serial primary key,
    updated_at timestamptz
);

create table groups(
                       id serial primary key,
                       user_id int references users(uid),
                       group_id int,
                       UNIQUE (user_id, group_id)
);

create table groups_events_sharing (
                                       id serial primary key,
                                       group_id int references groups(id),
                                       event_id int references events(id),
                                       is_admin bool
);

insert into syncer(updated_at) values (now());


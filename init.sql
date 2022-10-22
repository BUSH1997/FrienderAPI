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
    uid varchar(256),
    title varchar(256),
    description text,
    images text,
    starts_at bigserial,
    time_created timestamptz,
    time_updated timestamptz,
    geo varchar(256),
    category_id int references categories(id),
    is_group bool,
    is_public bool,
    owner_id int references users(id)
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


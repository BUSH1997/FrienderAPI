drop table subscribe_profile_sharing cascade;
drop table revindex_words cascade;
drop table revindex_events cascade;
drop table unlocked_awards cascade;
drop table messages cascade;
drop table awards cascade;
drop table unlocked_statuses cascade;
drop table subscribe_sharings cascade;
drop table groups_events_sharing cascade;
drop table groups cascade;
drop table event_sharings cascade;
drop table events cascade;
drop table users cascade;
drop table statuses cascade;
drop table conditions cascade;
drop table categories cascade;
drop table syncer cascade;


create table categories(
                           id serial primary key,
                           name varchar(256)
);

create table conditions(
                           id serial primary key,
                           created_events int,
                           visited_events int
);

create table statuses(
                         id serial primary key,
                         uid int UNIQUE,
                         title varchar(256),
                         condition_id int references  conditions(id)
);

create table users(
                      id serial primary key,
                      uid int UNIQUE,
                      current_status int references statuses(id),
                      created_events int,
                      visited_events int,
                      is_group bool
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
                       geo text,
                       category_id int references categories(id),
                       count_members int,
                       is_public bool,
                       is_private bool,
                       owner_id int references users(id),
                       is_deleted bool,
                       photos text,
                       members_limit int,
                       source varchar(256),
                       ticket varchar(256)
);

create table event_sharings(
                               id serial primary key,
                               event_id int references events(id),
                               user_id int references  users(id),
                               priority int,
                               is_deleted bool
);

create table groups(
                       id serial primary key,
                       user_id int references users(uid),
                       group_id int,
                       allow_user_events bool,
                       UNIQUE (user_id, group_id)
);

create table groups_events_sharing (
                                       id serial primary key,
                                       group_id int references groups(id),
                                       event_id int references events(id),
                                       is_admin bool,
                                       is_need_approve bool
);

create table subscribe_sharings(
                                   id serial primary key,
                                   user_id int references  users(id),
                                   subscriber_id int references  users(id)
);

create table unlocked_statuses(
                                  id serial primary key,
                                  user_id int references users(id),
                                  status_id int references statuses(id)
);

create table awards(
                       id serial primary key,
                       image varchar(256),
                       name varchar(256),
                       description varchar(256),
                       condition_id int references  conditions(id)
);

create table unlocked_awards(
                                id serial primary key,
                                user_id int references users(id),
                                award_id int references awards(id)
);

create table syncer(
                       id serial primary key,
                       updated_at timestamptz
);

create table messages(
                         id serial primary key,
                         user_id int references users(id),
                         user_uid int,
                         event_id int references events(id),
                         event_uid varchar(256),
                         text text,
                         time_created bigint
);

create table revindex_words(
                               id serial primary key,
                               word varchar(256) UNIQUE,
                               events int[]
);

create table revindex_events(
                                id serial primary key,
                                uid varchar(256) UNIQUE
);

create table subscribe_profile_sharing (
                                           id serial primary key,
                                           profile_id int references users(id),
                                           user_id int references users(id)

);

insert into categories(name) values ('Концерт'), ('Выставка'), ('Кино'), ('Экскурсия'), ('Спорт'), ('Театр'), ('Шоу'),
                                    ('Мастер-класс'), ('Бизнес'), ('It'), ('Воркшоп'), ('Флешмоб'), ('Другое');
insert into conditions(created_events, visited_events) values (0, 0);
insert into statuses(uid, title, condition_id) values (1, 'DEFAULT STATUS', 1);
insert into syncer(updated_at) values (now());

insert into users(uid, current_status) values (1,1);
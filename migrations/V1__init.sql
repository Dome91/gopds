create table users
(
    username   varchar primary key,
    password   varchar not null,
    role       varchar not null,
    created_at datetime default current_timestamp
);

create table sources
(
    id         varchar primary key,
    name       varchar not null unique,
    path       varchar not null unique,
    created_at datetime default current_timestamp
);

create table catalog_entries
(
    id                     varchar primary key,
    name                   varchar not null,
    path                   varchar not null,
    is_directory           integer not null,
    found_during_last_sync integer not null default true,
    type                   varchar not null,
    created_at             datetime         default current_timestamp,
    parent_catalog_entry   varchar,
    source                 varchar not null,

    foreign key (parent_catalog_entry) references catalog_entries (id) on delete cascade,
    foreign key (source) references sources (id) on delete cascade
);


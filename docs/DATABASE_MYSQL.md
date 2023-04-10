
```mysql
create table clickhouse_servers
(
    id         varchar(36)  not null primary key,
    org_id     varchar(36)  not null,
    host       varchar(256) not null,
    port       varchar(16),
    username   varchar(256) not null,
    shards     json,
    created_at bigint,
    updated_at bigint
);

CREATE TABLE clickhouse_users
(
    username        varchar(256) not null primary key,
    password        text,
    profile         varchar(64),
    quota           varchar(64),
    any_networks    varchar(64),
    allow_databases text,
    xml_path        text,
    status          varchar(64),
    created_at      bigint,
    updated_at      bigint
);
```
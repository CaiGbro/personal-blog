create table comments
(
    id             bigserial
        primary key,
    nickname       varchar(100) not null,
    content        text         not null,
    created_at     timestamp with time zone default now(),
    parent_id      bigint
        references comments,
    attachment_url text,
    reactions      jsonb
);

alter table comments
    owner to postgres;

create table video_comments
(
    id             bigserial
        primary key,
    video_filename varchar(255)                                       not null,
    nickname       varchar(255)                                       not null,
    content        text                                               not null,
    created_at     timestamp with time zone default CURRENT_TIMESTAMP not null,
    parent_id      bigint
        constraint fk_parent_video_comment
            references video_comments
            on delete cascade,
    attachment_url varchar(2048),
    reactions      jsonb                    default '{}'::jsonb
);

alter table video_comments
    owner to postgres;

create index idx_video_comments_filename_parent_id
    on video_comments (video_filename, parent_id);

create index idx_video_comments_created_at
    on video_comments (created_at desc);

create table video_reactions
(
    id             serial
        primary key,
    video_filename varchar(255) not null
        unique,
    reactions      jsonb default '{}'::jsonb
);

alter table video_reactions
    owner to postgres;

create index idx_video_filename_reactions
    on video_reactions (video_filename);

create table writing_comments
(
    id               serial
        primary key,
    writing_filename varchar(255) not null,
    nickname         varchar(100) not null,
    content          text         not null,
    created_at       timestamp with time zone default now(),
    parent_id        integer
        references writing_comments
            on delete cascade,
    attachment_url   varchar(255),
    reactions        jsonb                    default '{}'::jsonb
);

alter table writing_comments
    owner to postgres;

create table writing_reactions
(
    writing_filename text not null
        primary key,
    reactions        jsonb
);

alter table writing_reactions
    owner to postgres;

create table system_settings
(
    key   varchar(255) not null
        primary key,
    value text
);

alter table system_settings
    owner to postgres;

create table verification_codes
(
    email      varchar(255)             not null
        primary key,
    code       varchar(10)              not null,
    expires_at timestamp with time zone not null
);

alter table verification_codes
    owner to postgres;

create table visitor_logs
(
    id         serial
        primary key,
    hashed_ip  varchar(64)               not null,
    visit_date date default CURRENT_DATE not null,
    unique (hashed_ip, visit_date)
);

alter table visitor_logs
    owner to postgres;

create index idx_visitor_logs_visit_date
    on visitor_logs (visit_date);


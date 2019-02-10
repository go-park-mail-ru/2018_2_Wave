DROP TABLE IF EXISTS app CASCADE;

CREATE TABLE app (
    appid           BIGSERIAL       PRIMARY KEY,
    url             TEXT            DEFAULT '',
    link            TEXT            ,
    name            TEXT            ,
    image           TEXT            ,
    about           TEXT            ,
    installs        INT             DEFAULT 0,
    category        TEXT
);


DROP TABLE IF EXISTS app;

CREATE TABLE app (
    appid           BIGSERIAL       PRIMARY KEY,
    link            TEXT            ,
    name            TEXT            ,
    image           TEXT            ,
    about           TEXT            ,
    installs        INT             DEFAULT 0,
    price           INT             DEFAULT 0,
    category        TEXT
);
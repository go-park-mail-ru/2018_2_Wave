DROP TABLE IF EXISTS app;

CREATE TABLE app (
    appid           BIGSERIAL       PRIMARY KEY,
    name            TEXT            ,
    cover           TEXT            ,
    description     TEXT            ,
    installations   INT             ,
    price           INT             DEFAULT 0
);
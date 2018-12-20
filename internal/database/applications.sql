DROP TABLE IF EXISTS app;

CREATE TABLE app (
    appid           BIGSERIAL       PRIMARY KEY,
    name            TEXT            ,
    cover           TEXT            ,
    description     TEXT            ,
    installations   INT             DEFAULT 0,
    price           INT             DEFAULT 0,
    year            TEXT
);
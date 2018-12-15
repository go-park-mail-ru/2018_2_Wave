DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS userinfo;

-- User Tables Block
CREATE TABLE userinfo (
    uid         BIGSERIAL       PRIMARY KEY,
    username    TEXT            NOT NULL,
    password    VARCHAR(60)     NOT NULL,
    score       INT             DEFAULT 0,
    avatar      TEXT            DEFAULT '/img/avatars/default'

    appid       BIGSERIAL,
);

CREATE TABLE session (
    uid         SERIAL,
    cookie      TEXT            DEFAULT '',

    FOREIGN KEY (uid) REFERENCES userinfo(uid)
);

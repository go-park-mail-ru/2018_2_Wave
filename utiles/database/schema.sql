-- User Tables Block
CREATE TABLE userinfo (
    uid         BIGSERIAL       PRIMARY KEY,
    username    TEXT            NOT NULL,
    password    BYTEA           NOT NULL,
    score       INT             DEFAULT 0, 
    avatar      BYTEA           DEFAULT '/some/path'
);

CREATE TABLE session (
    uid         SERIAL,
    cookie      TEXT            DEFAULT '',

    FOREIGN KEY (uid) REFERENCES userinfo(uid)
);
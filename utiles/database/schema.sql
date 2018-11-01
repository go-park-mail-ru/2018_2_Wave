-- User Tables Block
CREATE TABLE userinfo (
    uid         BIGSERIAL       PRIMARY KEY,
    username    TEXT            NOT NULL,
    password    TEXT            NOT NULL,
    score       INT             DEFAULT 0, 
    avatar      TEXT            DEFAULT '/some/paths/lead/to/beautiful/destinations/'
);

CREATE TABLE session (
    uid         SERIAL,
    cookie   TEXT            DEFAULT '',

    FOREIGN KEY (uid) REFERENCES userinfo(uid)
);
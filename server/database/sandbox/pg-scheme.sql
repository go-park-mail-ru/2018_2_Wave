CREATE TABLE UserInfo (
    uid         SERIAL          PRIMARY KEY,
    username    TEXT            NOT NULL,
    password    TEXT            NOT NULL,
    score       INT             DEFAULT 0, 
    avatar      TEXT,           DEFAULT '/some/paths/lead/to/beautiful/destinations/',
    cookie      TEXT            DEFAULT ''
);
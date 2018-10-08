-- Cleanup
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS cookies;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS score;

-- Main tables
CREATE TABLE users (
    uid         serial  PRIMARY KEY,
    username    text    NOT NULL,
    watchword   text    NOT NULL
);

CREATE TABLE profiles (
    uid         int     PRIMARY KEY,
    avatar      bytea

    FOREIGN KEY(uid) REFERENCES users(uid)
)

CREATE TABLE cookies (
    uid     int     PRIMARY KEY
    cookie  text

    FOREIGN KEY(uid) REFERENCES users(uid)
);

-- Gameplay tables
CREATE TABLE score (
    uid     int     PRIMARY KEY
    value   int     DEFAULT 0

    FOREIGN KEY(uid) REFERENCES users(uid)
)
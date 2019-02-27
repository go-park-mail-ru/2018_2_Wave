DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS userinfo CASCADE;
DROP TABLE IF EXISTS userapp;
DROP TABLE IF EXISTS app CASCADE;

-- User Tables Block
CREATE TABLE userinfo (
    uid         BIGSERIAL       PRIMARY KEY,
    username    TEXT            NOT NULL,
    password    VARCHAR(60)     NOT NULL,
    score       INT             DEFAULT 0,
    avatar      TEXT            DEFAULT 'img/avatars/default'
);

CREATE TABLE session (
    uid         SERIAL,
    cookie      TEXT            DEFAULT '',

    FOREIGN KEY (uid) REFERENCES userinfo(uid)
);

CREATE TABLE app (
    appid           BIGSERIAL       PRIMARY KEY,
    url             TEXT            DEFAULT '',
    link            TEXT            UNIQUE NOT NULL,
    name            TEXT            NOT NULL,
    name_de         TEXT            DEFAULT '',
    name_ru         TEXT            DEFAULT '',
    image           TEXT            DEFAULT '',
    about           TEXT            DEFAULT '',
    about_de        TEXT            DEFAULT '',
    about_ru        TEXT            DEFAULT '',
    installs        INT             DEFAULT 0,
    category        TEXT            NOT NULL
);


CREATE TABLE userapp (
    uid             BIGSERIAL           ,
    appid           BIGSERIAL           ,

    FOREIGN KEY (uid) REFERENCES userinfo(uid),
    FOREIGN KEY (appid) REFERENCES app(appid),
    PRIMARY KEY (uid, appid)
);

CREATE OR REPLACE FUNCTION update_installs()
RETURNS TRIGGER AS '
    BEGIN
        UPDATE app
            SET installs = installs + 1
            WHERE appid = NEW.appid;
            RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER update_installs
AFTER INSERT ON userapp
FOR EACH ROW EXECUTE PROCEDURE update_installs();

INSERT INTO app(link, name, image, about, category) VALUES ('/terminal', 'Terminal', 'img/app_covers/terminal.jpg', 'App by Wave', '2018_2');
INSERT INTO app(link, name, image, about, category) VALUES ('/snake', 'Snake', 'img/app_covers/snake.jpg', 'Game by Wave', '2018_2');

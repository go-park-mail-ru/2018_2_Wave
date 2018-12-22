DROP TABLE IF EXISTS userapp;

CREATE TABLE userapp (
    uid             BIGSERIAL           ,
    appid           BIGSERIAL           ,
    time_total      REAL       DEFAULT 0,
    time_start       REAL       DEFAULT 0,
    time_ping       REAL       DEFAULT 0,б

    FOREIGN KEY (uid) REFERENCES userinfo(uid),
    FOREIGN KEY (appid) REFERENCES app(appid)
);
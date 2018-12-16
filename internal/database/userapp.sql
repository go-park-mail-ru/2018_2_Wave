DROP TABLE IF EXISTS userapp;

CREATE TABLE userapp (
    uid             BIGSERIAL,
    appid           BIGSERIAL,
    time_total      REAL,

    FOREIGN KEY (uid) REFERENCES userinfo(uid),
    FOREIGN KEY (appid) REFERENCES app(appid)
);
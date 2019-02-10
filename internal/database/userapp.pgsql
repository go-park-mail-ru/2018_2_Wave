DROP TABLE IF EXISTS userapp;

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

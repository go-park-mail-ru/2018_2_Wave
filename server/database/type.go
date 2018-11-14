package database

import (
	"Wave/utiles/config"
	lg "Wave/utiles/logger"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //psql driver
)

// DatabaseModel - database facade
type DatabaseModel struct {
	DBconf   config.DatabaseConfiguration
	Database *sqlx.DB
	LG       *lg.Logger
}

func New(dbconf_ config.DatabaseConfiguration, lg_ *lg.Logger) *DatabaseModel {
	postgr := &DatabaseModel{
		DBconf: dbconf_,
		LG:     lg_,
	}

	var err error
	postgr.Database, err = sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname='%s' sslmode=disable", postgr.DBconf.User, os.Getenv("WAVE_DB_PASSWORD"), postgr.DBconf.DBName))
	if err != nil {
		//log.Fatalln(err)
		panic(err)
	}
	postgr.LG.Sugar.Infow("happened form new")
	log.Println("postgres connection established")

	return postgr
}

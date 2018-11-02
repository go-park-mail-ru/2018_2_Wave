package app

import (
	"Wave/app/middleware"
	"Wave/app/middlewareGlob"

	"Wave/utiles/configs"
	"Wave/utiles/logger"
	"Wave/utiles/walhalla"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

//go:generate go run ../utiles/walhalla/main .
//go:generate go run ../utiles/configs/main .

// walhalla:app {
// 	globalMiddlewares  : [ cors, log ],
// 	operationMiddlewars: [ recovery  ]
// }

var MiddlewareGenerators = walhalla.MiddlewareGenerationFunctionMap{
	"auth":     middleware.AuthTrue,
	"recovery": middleware.Recovery,
}

var MiddlewareGeneratorsGlobal = walhalla.GlobalMiddlewareGenerationFunctionMap{
	"log":  middlewareGlob.Logger,
	"cors": middlewareGlob.Cors,
}

func SetupContext(ctx *walhalla.Context) {
	var err error
	{ // read a configuratoin file
		Conf := new(configs.MainConfig)
		Conf.ReadFromFile("config.json")
		ctx.Config = Conf
	}
	{ // setup the logger
		ctx.Log, err = logger.New(logger.Config{
			File:    "log.log",
			BStdOut: true,
			BAsync:  true,
		})
		if err != nil {
			panic(err)
		}
	}
	{ // setup database
		//connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", ctx.Config.database.user, ctx.Config.database.dbname)
		conStr := "user=waveapp password=surf dbname=wave sslmode=disable"
		ctx.DB, err = sqlx.Connect("postgres", conStr)

		if err != nil {
			ctx.Log.Error(err)
		}
		ctx.Log.Info("connection to postgres succesfully established")
	}
}

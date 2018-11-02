package app

import (
	"Wave/app/middleware"
	"Wave/app/middlewareGlob"

	"Wave/utiles/configs"
	"Wave/utiles/logger"
	"Wave/utiles/walhalla"

	_ "github.com/lib/pq" // psql driver
)

//go:generate go run ../utiles/walhalla/main .
//go:generate go run ../utiles/configs/main config.json

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
	var (
		err error
		conf *configs.MainConfig
	)
	{ // read a configuratoin file 
		conf = new(configs.MainConfig)
		conf.ReadFromFile("config.json")
		ctx.Config = conf
	}
	{ // setup the logger
		ctx.Log, err = logger.New(logger.Config{
			File:    conf.Server.Log,
			BStdOut: true,
		})
		if err != nil {
			panic(err)
		}
	}
	{ // setup database
		ctx.InitDatabase(conf.Database)
		ctx.Log.Info("connection to postgres succesfully established")
	}
}

package app

import (
	"Wave/app/middleware"
	"Wave/app/middlewareGlob"
	
	"Wave/utiles/walhalla"
	"Wave/utiles/logger"
	"Wave/utiles/configs"
)

//go:generate go run ../utiles/walhalla/main .
//go:generate go run ../utiles/configs/main .

// walhalla:app {
// 	globalMiddlewares  : [ cors, log ],
// 	operationMiddlewars: [ recovery  ]
// }

var MiddlewareGenerators = walhalla.MiddlewareGenerationFunctionMap{
	"auth": middleware.AuthTrue, 
	"recovery": middleware.Recovery,
}

var MiddlewareGeneratorsGlobal = walhalla.GlobalMiddlewareGenerationFunctionMap{
	"log": middlewareGlob.Logger,
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
			File: "log.log",
			BStdOut: true,
			BAsync:  true,
		})
		if err != nil {
			panic(err)
		}
	}
	{ // setup databse
		
	}
}

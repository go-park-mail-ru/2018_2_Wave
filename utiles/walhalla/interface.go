package walhalla

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"
)

//go:generate go-bindata -pkg assets -o _assets/assets.go templates/...
//go:generate easyjson -output_filename main/types.gen.go main/types.go

// ----------------| global middlewares

// GlobalMiddlewareFunction - pre-router middleware
// Note: the middleware can change an input data
// @param http.ResponseWriter 	- responce object
// @param http.Request 			- request object
// @param Context 				- handler/server conetext
type GlobalMiddlewareFunction func(http.ResponseWriter, *http.Request)

// GlobalMiddlewareGenerationFunction - pre-router middleware
// @param GlobalMiddlewareFunction 	- next function
// @return GlobalMiddlewareFunction	- middleware chain
type GlobalMiddlewareGenerationFunction func(GlobalMiddlewareFunction, *Context) GlobalMiddlewareFunction

// GlobalMiddlewareGenerationFunctionMap - map of global middlewares
// @param string								- middleware name
// @param GlobalMiddlewareGenerationFunctionMap - middleware generator
type GlobalMiddlewareGenerationFunctionMap map[string]GlobalMiddlewareGenerationFunction

// ----------------| operation middlewares

// MiddlewareFunction - middleware cannot modify payload;
// Note: The middleware will be applyed after router
// @param htt.Request 			- request object. Note: All payload objects has been parsed
// @param Conext      			- handler/server conetext
// @return middleware.Responder - simple responce (Code responce)
type MiddlewareFunction func(*http.Request) middleware.Responder

// MiddlewareGenerationFunction appends new middleware to the middleware chain
// @param MiddlewareFunction 	- next function
// @param Context 				- handler/server conetext
// @return MiddlewareFunction	- middleware chain
type MiddlewareGenerationFunction func(MiddlewareFunction, *Context) MiddlewareFunction

// MiddlewareGenerationFunctionMap - map of after router middlewares
// @param string 						- middleware name
// @param MiddlewareGeneratorFunction 	- middleware generator
type MiddlewareGenerationFunctionMap map[string]MiddlewareGenerationFunction

// ----------------| context

// Fields - Log fields
type Fields map[string]interface{}

// ILogger - walhalla logger interface
type ILogger interface {
	WithFields(fields Fields) ILogger

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// Context - operation context
type Context struct {
	OperationID string
	Log         ILogger
	DB          *sqlx.DB
	Config      interface{}
}

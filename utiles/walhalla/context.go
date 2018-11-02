package walhalla

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Context - operation context
type Context struct {
	Log        ILogger
	DB         *sqlx.DB
	Config     interface{}
	request    *http.Request
	outCookies []*http.Cookie
}

// Copy the curent context
func (ctx *Context) Copy() Context {
	var (
		cpy    = *ctx
		length = len(ctx.outCookies)
	)
	cpy.outCookies = make([]*http.Cookie, length)
	copy(cpy.outCookies, ctx.outCookies)
	return cpy
}

// ----------------| cookie

// SetCookie to responce
func (ctx *Context) SetCookie(c *http.Cookie) {
	ctx.outCookies = append(ctx.outCookies, c)
}

// GetCookie from the request
func (ctx *Context) GetCookie(name string) string {
	if cookie, err := ctx.request.Cookie(name); err != nil && cookie != nil {
		return cookie.Value
	}
	return ""
}

// SetRequest to the context
func SetRequest(ctx *Context, r *http.Request) {
	ctx.request = r
}

// ----------------| database

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (dc *DatabaseConfig) Marshal() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dc.Host, dc.Port, dc.User, dc.Password, dc.DBName)
}

func (ctx *Context) InitDatabase(dc DatabaseConfig) {
	ctx.DB = sqlx.MustOpen("postgres", dc.Marshal())
}

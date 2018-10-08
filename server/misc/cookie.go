package misc

import (
	"math/rand"
	"time"

	"github.com/valyala/fasthttp"
)

const CookieStringLenght = 32
const sessionCookieLifeTime = 60 * 24 * 365
const sessionCookieName = "session"

//*****************|

func GenerateCookie() string {
	return RandString(CookieStringLenght)
}

func MakeSessionCookie(value string) *fasthttp.Cookie {
	loginCookie := &fasthttp.Cookie{}
	loginCookie.SetMaxAge(sessionCookieLifeTime)
	loginCookie.SetKey(sessionCookieName)
	loginCookie.SetValue(value)
	loginCookie.SetPath("")
	return loginCookie
}

func GetSessionCookie(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(sessionCookieName))
}

func SetCookie(ctx *fasthttp.RequestCtx, cookie *fasthttp.Cookie) {
	ctx.Response.Header.SetCookie(cookie)
}

//*****************|

func RandString(n int) string {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	src := rand.NewSource(time.Now().UnixNano() + rand.Int63())

	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

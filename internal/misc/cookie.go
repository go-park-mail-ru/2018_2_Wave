package misc

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	cookieStringLenght    = 64
	sessionCookieLifeTime = 60 * 24 * 365
	sessionCookieName     = "session"
)

func GenerateCookie() string {
	return RandString(cookieStringLenght)
}

func MakeSessionCookie(value string) *http.Cookie {
	loginCookie := &http.Cookie{}
	loginCookie.MaxAge = sessionCookieLifeTime
	loginCookie.Name = sessionCookieName
	loginCookie.Value = value
	loginCookie.Path = ""
	loginCookie.Secure = true
	return loginCookie
}

func GetSessionCookie(r *http.Request) string {
	session, err := r.Cookie(sessionCookieName)
	if err != nil || session == nil {
		return ""
	}
	return session.Value
}

func SetCookie(w http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
}

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

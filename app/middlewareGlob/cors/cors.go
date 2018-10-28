package cors

import (
	"strconv"
	"strings"

	"net/http"

	"Wave/utiles/walhalla"
)

// ---------------- |

type Options struct {
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	ExposedHeaders     []string
	OptionsPassthrough bool
	AllowCredentials   bool
	MaxAge             int
}

type CorsHandler struct {
	// origins
	allowedOriginsAll bool
	allowedOrigins    []string
	// headers
	allowedHeadersAll    bool
	allowedHeaders       []string
	exposedHeaders       []string
	joinedExposedHeaders string
	// methods
	allowedMethods []string
	// misc
	allowCredentials  bool
	optionPassthrough bool
	maxAge            int
}

// ---------------- |

type params struct {
	origin        string
	method        string
	headers       []string
	joinedHeaders string
}

// ---------------- |

func New(options Options) *CorsHandler {
	cors := &CorsHandler{
		exposedHeaders:    options.ExposedHeaders,
		allowCredentials:  options.AllowCredentials,
		maxAge:            options.MaxAge,
		optionPassthrough: options.OptionsPassthrough,
	}
	{ // Allowed Origins
		if len(options.AllowedOrigins) == 0 {
			cors.allowedOriginsAll = true
		} else {
			cors.allowedOrigins = []string{}
			for _, origin := range options.AllowedOrigins {
				origin = strings.ToLower(origin)
				if origin == "*" {
					cors.allowedOriginsAll = true
					cors.allowedOrigins = nil
					break
				} else {
					cors.allowedOrigins = append(cors.allowedOrigins, origin)
				}
			}
		}
	}
	{ // Allowed Headers
		if len(options.AllowedHeaders) == 0 {
			cors.allowedHeaders = []string{"Origin", "Accept", "Content-Type", "X-Requested-With"}
		} else {
			cors.allowedHeaders = append(options.AllowedHeaders, "Origin")
			for _, h := range options.AllowedHeaders {
				if h == "*" {
					cors.allowedHeadersAll = true
					cors.allowedHeaders = nil
					break
				}
			}
		}
	}
	{ // Allowed Methods
		if len(options.AllowedMethods) == 0 {
			cors.allowedMethods = []string{"GET", "POST", "HEAD"}
		} else {
			cors.allowedMethods = options.AllowedMethods
			for _, method := range cors.allowedMethods {
				method = strings.ToUpper(method)
			}
		}
	}
	{ // exposed headers
		if len(cors.exposedHeaders) > 0 {
			cors.joinedExposedHeaders = strings.Join(cors.exposedHeaders, ", ")
		}
	}

	return cors
}

func (c *CorsHandler) CorsMiddleware(next walhalla.GlobalMiddlewareFunction) walhalla.GlobalMiddlewareFunction {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			if c.handlePreflight(rw, r) {
				// correct headers
				if c.optionPassthrough {
					next(rw, r)
				} else {
					rw.WriteHeader(http.StatusOK)
				}
			} else {
				// incorrect headers
				rw.WriteHeader(http.StatusForbidden)
			}
		} else {
			if c.handleActual(rw, r) {
				// correct headers
				next(rw, r)
			} else {
				// incorrect headers
				rw.WriteHeader(http.StatusForbidden)
			}
		}
	}
}

// ---------------- |

func extractParams(r *http.Request) *params {
	origin := string(r.Header.Get("Origin"))
	method := string(r.Header.Get("Access-Control-Request-Method"))
	headers := string(r.Header.Get("Access-Control-Request-Headers"))
	return &params{
		origin:        origin,
		method:        method,
		headers:       strings.Split(headers, ","),
		joinedHeaders: headers,
	}
}

func (p *params) updateParams(c *CorsHandler) {
	if c.allowedOriginsAll {
		p.origin = "*"
	}
}

// ---------------- |

func setOrigin(rw http.ResponseWriter, val string) {
	rw.Header().Set("Access-Control-Allow-Origin", val)
}

func setMethods(rw http.ResponseWriter, val string) {
	rw.Header().Set("Access-Control-Allow-Methods", val)
}

func setHeaders(rw http.ResponseWriter, val string) {
	if val != "" {
		rw.Header().Set("Access-Control-Allow-Headers", val)
	}
}

func setExposedHeaders(rw http.ResponseWriter, val string) {
	if val != "" {
		rw.Header().Set("Access-Control-Expose-Headers", val)
	}
}

func setCredentials(rw http.ResponseWriter, val bool) {
	if val {
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func setMaxAge(rw http.ResponseWriter, val int) {
	if val > 0 {
		rw.Header().Set("Access-Control-Max-Age", strconv.Itoa(val))
	}
}

// ---------------- |

func (c *CorsHandler) handlePreflight(rw http.ResponseWriter, r *http.Request) bool {
	params := extractParams(r)
	{ // validate request
		if !c.isAllowedOrigin(params.origin) {
			return false
		}
		if !c.isAllowedMethod(params.method) {
			return false
		}
		if !c.areHeadersAllowed(params.headers) {
			return false
		}
	}
	params.updateParams(c)
	{ // set headers
		setMaxAge(rw, c.maxAge)
		setOrigin(rw, params.origin)
		setMethods(rw, params.method)
		setHeaders(rw, params.joinedHeaders)
		setCredentials(rw, c.allowCredentials)
	}
	return true
}

func (c *CorsHandler) handleActual(rw http.ResponseWriter, r *http.Request) bool {
	params := extractParams(r)
	{ // validate request
		if !c.isAllowedOrigin(params.origin) {
			return false
		}
		if !c.isAllowedMethod(r.Method) {
			return false
		}
	}
	params.updateParams(c)
	{ // set headers
		setOrigin(rw, params.origin)
		setHeaders(rw, params.joinedHeaders)
		setExposedHeaders(rw, c.joinedExposedHeaders)
		setCredentials(rw, c.allowCredentials)
	}
	return true
}

func (c *CorsHandler) isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	if c.allowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	for _, val := range c.allowedOrigins {
		if val == origin {
			return true
		}
	}
	return false
}

func (c *CorsHandler) isAllowedMethod(method string) bool {
	if len(c.allowedMethods) == 0 {
		return false
	}
	method = strings.ToUpper(method)
	if method == "OPTIONS" {
		return true
	}
	for _, m := range c.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

func (c *CorsHandler) areHeadersAllowed(headers []string) bool {
	if c.allowedHeadersAll || len(headers) == 0 {
		return true
	}
	for _, header := range headers {
		found := false
		for _, h := range c.allowedHeaders {
			if h == header {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

package handles

import (
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"

	"github.com/valyala/fasthttp"
)

const sessionCookieLifeTime = 60 * 24 * 365
const sessionCookieName = "session"

func makeSessionCookie(value string) *fasthttp.Cookie {
	loginCookie := &fasthttp.Cookie{}
	loginCookie.SetMaxAge(sessionCookieLifeTime)
	loginCookie.SetKey(sessionCookieName)
	loginCookie.SetSecure(false)
	loginCookie.SetValue(value)
	return loginCookie
}

func getSessionCookie(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(sessionCookieName))
}

func setCookie(ctx *fasthttp.RequestCtx, cookie *fasthttp.Cookie) {
	ctx.Response.Header.SetCookie(cookie)
}

// Pars the form into a flat structure pointer
// @param form 		- target form
// @param target 	- pointer on a flat structure
func parsForm(form *multipart.Form, target interface{}) error {
	if target == nil {
		return errors.New("nilled target is unallowed")
	}
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return errors.New("target must be pointer")
	}

	var (
		object     = reflect.ValueOf(target).Elem()
		objectType = object.Type()
		objectKind = object.Kind()
	)
	if objectKind != reflect.Struct {
		return errors.New("target must be struct pointer")
	}

	for i := 0; i < object.NumField(); i++ {
		var (
			field     = object.Field(i)
			fieldType = objectType.Field(i)
			tag       = fieldType.Tag.Get("json")
			fieldKind = fieldType.Type.Kind()
			fieldName = fieldType.Name
		)
		if tag == "-" {
			continue
		} else if tag != "" {
			fieldName = tag
		}

		switch fieldKind {
		case reflect.Slice:
			var (
				elementType = field.Type().Elem()
				elementKind = elementType.Kind()
			)
			if elementKind != reflect.Uint8 {
				return errors.New("Only []bytes are allowed")
			}

			files, ok := form.File[fieldName]
			if len(files) == 0 || !ok {
				continue
			}
			if len(files) > 1 {
				return fmt.Errorf("Form files %s cannot be multipl", fieldName)
			}

			file := files[0]
			data := make([]byte, file.Size)
			entery, err := file.Open()
			if err != nil {
				return err
			}

			entery.Read(data)
			field.SetBytes(data)

		case reflect.Array:
			return fmt.Errorf("Unsupported field type: %s", fieldKind)

		default:
			values, ok := form.Value[fieldName]
			if len(values) == 0 || !ok {
				continue
			}
			if len(values) > 1 {
				return fmt.Errorf("Form value %s cannot be multipl", fieldName)
			}

			switch fieldKind {
			case reflect.String:
				field.SetString(values[0])

			case reflect.Int:
				value, err := strconv.ParseInt(values[0], 10, 64)
				if err != nil {
					return fmt.Errorf("Inable to parse field %s", fieldName)
				}
				field.SetInt(value)

			default:
				return fmt.Errorf("Unsupported field type: %s", fieldKind)
			}
		}
	}

	return nil
}

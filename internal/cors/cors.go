package cors

func SetOrigin(origin string, allowedOrigins []string) string {
	for _, s := range allowedOrigins {
		if origin == s {
			return s
		}
	}

	return ""
}

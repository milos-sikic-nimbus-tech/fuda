package httpapi

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

type Credentials struct {
	User     string
	Password string
}

func (c Credentials) enabled() bool {
	return c.Password != ""
}

func basicAuth(creds Credentials) middleware {
	return func(next http.Handler) http.Handler {
		if !creds.enabled() {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" || strings.HasPrefix(r.URL.Path, "/api/webhooks/") {
				next.ServeHTTP(w, r)
				return
			}
			user, password, ok := r.BasicAuth()
			if !ok || !equal(user, creds.User) || !equal(password, creds.Password) {
				w.Header().Set("WWW-Authenticate", `Basic realm="fuda", charset="UTF-8"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func equal(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

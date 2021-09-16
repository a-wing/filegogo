package httpd

import (
	"net/http"
	"net/url"
)

// Fork From the: https://pkg.go.dev/net/http#StripPrefix
func NoPrefix(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			h.ServeHTTP(w, r)
		} else {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = ""
			r2.URL.RawPath = ""
			h.ServeHTTP(w, r2)
		}
	})
}

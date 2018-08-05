package server

import (
	"net/http"
)

// Method is a middleware and checks whether the request is sent with the proper HTTP method
func Method(method string, f http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "This endpoint must be called with a POST", http.StatusMethodNotAllowed)
			return
		}
		f.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

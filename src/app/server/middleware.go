package server

import (
	"fmt"
	"net/http"
)

// Method is a middleware and checks whether the request is sent with the proper HTTP method
func Method(method string, f http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("INSIDE METHOD: %+v -- %+v\n\n", method, r.Method)
		if r.Method != method {
			http.Error(w, "This endpoint must be called with a POST", http.StatusMethodNotAllowed)
			return
		}
		f.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

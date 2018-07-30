package server

import "net/http"

// Method is a middleware and checks whether the request is sent with the proper HTTP method
func Method(method string, f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "This endpoint must be called with a POST", http.StatusMethodNotAllowed)
			return
		}
		f(w, r)
	}
}

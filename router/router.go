package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ProvideRouter provides a gorilla mux router
func ProvideRouter() *mux.Router {
	var router = mux.NewRouter()
	router.Use(jsonMiddleware)
	return router
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

var Options = ProvideRouter

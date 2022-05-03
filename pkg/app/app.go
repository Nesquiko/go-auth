package app

import (
	"fmt"
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type GoAuthServer struct{}

func (s GoAuthServer) Signup(w http.ResponseWriter, r *http.Request) {

}

func (s GoAuthServer) Login(w http.ResponseWriter, r *http.Request) {
}

func StartServer() {
	fmt.Println("Starting server...")
	port := "8080"

	r := chi.NewRouter()
	middlewares := []api.MiddlewareFunc{
		func(hf http.HandlerFunc) http.HandlerFunc {
			return http.HandlerFunc(middleware.Logger(hf).ServeHTTP)
		},
	}

	var server GoAuthServer
	servOpts := api.ChiServerOptions{
		BaseRouter:  r,
		Middlewares: middlewares,
	}

	h := api.HandlerWithOptions(server, servOpts)

	fmt.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, h)
}

// Package app is provides functions for starting the Go-Auth application.
package app

import (
	"fmt"
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/db"
	"github.com/Nesquiko/go-auth/pkg/middleware"
	"github.com/Nesquiko/go-auth/pkg/server"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// StartServer starts the whole Go-Auth application. Firstly if tries to connect
// to a MySQL database, if it fails, the app won't start. Then creates new
// router and configures it with middleware and handler. The application listens
// on port 8080.
func StartServer() {
	fmt.Print("Connecting to Database...")
	err := db.ConnectDB(
		"mysql",
		db.MySQLDSNConfig("root", "goAuthDB", "127.0.0.1:3306", "users").FormatDSN(),
	)
	if err != nil {
		fmt.Print(" - \x1b[31;1mFAILED\x1b[0m\n")
		panic(err)
	}
	fmt.Print(" - \x1b[32;1mSUCCESS\x1b[0m\n")

	fmt.Println("Starting server...")
	port := "8080"

	r := chi.NewRouter()
	middlewares := []api.MiddlewareFunc{
		chiMiddleware.Logger,
		middleware.ContentTypeFilter,
	}

	var server server.GoAuthServer
	servOpts := api.ChiServerOptions{
		BaseRouter:  r,
		Middlewares: middlewares,
	}

	h := api.HandlerWithOptions(server, servOpts)

	fmt.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, h)
}

package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartServer() {
	fmt.Println("Starting server...")
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello there!"))
	})

	fmt.Println("Listening on port 3000...")
	http.ListenAndServe(":8080", r)
}

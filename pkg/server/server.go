package server

import (
	"fmt"
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/db"
)

type GoAuthServer struct{}

func (s GoAuthServer) Signup(w http.ResponseWriter, r *http.Request) {
	db, err := db.Connect("root", "goAuthDB", "users", "127.0.0.1:3306")
	if err != nil {
		panic(err)
	}

	users, err := db.FetchUsers()
	if err != nil {
		panic(err)
	}

	for _, u := range users {
		fmt.Println(u)
	}

}

func (s GoAuthServer) Login(w http.ResponseWriter, r *http.Request) {

}

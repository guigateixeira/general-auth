package handler

import (
	"fmt"
	"net/http"
)

type User struct{}

func (o *User) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create user")
}

func (o *User) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get user by ID")
}

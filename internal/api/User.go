package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mandre1899/GO_Webserver/internal/auth"
	"github.com/mandre1899/GO_Webserver/internal/database"
)

type User struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func LoginUser(db *database.Queries) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		decoder := json.NewDecoder(r.Body)
		var user User
		err := decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		if user.Email == "" || user.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		var resUser database.User
		resUser, err = db.GetUserPasswordHashByMail(context.Background(), user.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		valid, err := auth.CheckPasswordHash(user.Password, resUser.HashedPassword)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		if valid == false {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		var response database.GetUserByIdRow
		response, err = db.GetUserById(context.Background(), resUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		res, err := json.Marshal(&response)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(res))
	}
}

func CreateUser(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var user User
		err := decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		if user.Email == "" || user.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		var createParams database.CreateUserParams
		createParams.Email = user.Email
		createParams.HashedPassword, err = auth.HashPassword(user.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Password hashing failed"))
			return

		}
		fmt.Println("Hash:", createParams.HashedPassword)
		var dbUser database.CreateUserRow
		dbUser, err = db.CreateUser(context.Background(), createParams)
		w.Header().Add("Content-Type", "application-json")
		res, err := json.Marshal(&dbUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(res))
	}
}


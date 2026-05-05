package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mandre1899/GO_Webserver/internal/auth"
	"github.com/mandre1899/GO_Webserver/internal/database"
	"github.com/mandre1899/GO_Webserver/internal/middleware"
)

type User struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	User  database.GetUserByIdRow `json:"user"`
	Token string                  `json:"token"`
}

func LoginUser(db *database.Queries, apiConf *middleware.ApiConfig) http.HandlerFunc {
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
		var response database.GetUserByIdRow
		response, err = db.GetUserById(context.Background(), resUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}

		token, err := auth.MakeJWT(resUser.ID, apiConf.JWTSecret, time.Hour)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Token generation failed"))
			return
		}

		loginResp := LoginResponse{
			User:  response,
			Token: token,
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		res, err := json.Marshal(&loginResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Marshal failed"))
			return
		}
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
		jwtStr, err := auth.MakeJWT(dbUser.ID, "1234", time.Hour * 1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Password hashing failed"))
			return
		}
		w.Header().Add("Authorization", "Bearer " + jwtStr)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(res))
	}
}


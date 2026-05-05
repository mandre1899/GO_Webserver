package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mandre1899/GO_Webserver/internal/database"
)

type User struct {
	Email string `json:"email"`
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
		if user.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Client: bad request"))
			return
		}
		var dbUser database.User
		dbUser, err = db.CreateUser(context.Background(), user.Email)
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

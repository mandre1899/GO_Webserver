package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mandre1899/GO_Webserver/internal/database"
)

type req struct {
	Body    string `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

func GetChripById(db *database.Queries) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		val := r.PathValue("id")
		uuidVal, err := uuid.Parse(val)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Couldn't convert to UUID"))
			return
		}
		var chrip database.Chirp
		chrip, err = db.GetChripByID(context.Background(), uuidVal)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		res, err := json.Marshal(&chrip)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
}

func GetChrips(db *database.Queries) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		var chirps []database.Chirp
		chirps, err := db.GetAllChirps(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("DB query failed"))
			return
		}
		res, err := json.Marshal(&chirps)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("DB query failed"))
			return
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
}

func ValidateChirpHandler(db *database.Queries) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {

		}
		replacer := strings.NewReplacer(
			"Kerfuffle", "****",
			"Sharbert", "****",
			"Fornax", "****",
			"kerfuffle", "****",
			"sharbert", "****",
			"fornax", "****",
		)

		reqBody = []byte(replacer.Replace(string(reqBody)))
		var request req
		err = json.Unmarshal(reqBody, &request)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if 140 < len(reqBody) {
			respondWithJSON(w, http.StatusBadRequest, "{\"error\":\"Something went wrong\"}")
			return
		}
		var chripCreate database.CreateChirpParams
		chripCreate.Body = request.Body
		chripCreate.UserID = request.User_id
		w.Header().Add("content-type", "application/json")
		var dbRes database.Chirp
		dbRes, err = db.CreateChirp(context.Background(), chripCreate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("UUID not found"))
		}
		
		w.WriteHeader(http.StatusCreated)
		resp, err := json.Marshal(&dbRes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Marshal failed"))
		}
		w.Write([]byte(resp))
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	resBody := `{"error":\"`+ msg + `\"}`
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(resBody))
}

func respondWithJSON(w http.ResponseWriter, code int, payload string) {
	w.WriteHeader(code)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(payload))
}


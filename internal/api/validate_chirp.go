package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type req struct {
	Body string `json:"body"`
}

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
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

	respondWithJSON(w, http.StatusOK, "{\"cleaned_body\":\"" + request.Body + "\"}")
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


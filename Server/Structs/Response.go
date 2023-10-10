package Structs

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Response struct {
	Message string `json:"outputString"`
	Code    int    `json:"-"`
}

func (response *Response) Json(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

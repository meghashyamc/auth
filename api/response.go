package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type response struct {
	Success bool        `json:"success"`
	Errors  []string    `json:"errors"`
	Data    interface{} `json:"data"`
}

func writeResponse(w http.ResponseWriter, statusCode int, success bool, errors []string, data interface{}) {

	jsonBytes, err := json.Marshal(response{Success: success, Errors: errors, Data: data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not marshal response")
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonBytes)
	return

}

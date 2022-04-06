package api

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func getTokenValidity() (time.Time, error) {

	tokenValidityDuration, err := time.ParseDuration(os.Getenv("CONFIRMATION_TOKEN_VALIDITY_DURATION"))
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not parse confirmation token validity duration")
		return time.Time{}, err
	}

	return time.Now().UTC().Add(tokenValidityDuration), nil
}

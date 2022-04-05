package api

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var tokenValidityDuration = 10 * time.Minute

// Hash password using the bcrypt hashing algorithm
func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not create password digest")
		return "", err
	}
	return string(hashedPasswordBytes), nil
}

func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func generateAuthToken(id string) (string, error) {
	timeNow := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"issue_time":  timeNow.Unix(),
		"expiry_time": timeNow.Add(tokenValidityDuration).Unix(),
	})

	tokenString, err := token.SignedString(os.Getenv("AUTH_SECRET"))
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not sign token")
		return "", err
	}

	return tokenString, nil

}

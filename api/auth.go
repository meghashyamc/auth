package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"html/template"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type emailConfirmation struct {
	UID   string
	Token string
}

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

func generateConfirmationLink(id uuid.UUID) (string, string, error) {

	userIDDigest := getInsecureStringDigest(id.String())
	confirmationToken, err := generateConfirmationToken()
	if err != nil {
		return "", "", err
	}
	confirmationLink, err := formConfirmationLink(userIDDigest, confirmationToken)
	return confirmationToken, confirmationLink, err

}

func formConfirmationLink(userIDDigest, confirmationToken string) (string, error) {
	t := template.Must(template.New("confirmationEmailLinkTemplate").Parse(os.Getenv("CONFIRMATION_EMAIL_LINK_TEMPLATE")))
	buf := &bytes.Buffer{}

	if err := t.Execute(buf, &emailConfirmation{UID: userIDDigest, Token: confirmationToken}); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not form email confirmation link")
		return "", err
	}
	return buf.String(), nil
}

func getInsecureStringDigest(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))

}

func generateConfirmationToken() (string, error) {

	newUUID, err := uuid.NewV4()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not generate email confirmation token")
		return "", err
	}

	return getInsecureStringDigest(newUUID.String()), nil

}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meghashyamc/auth/models"
	log "github.com/sirupsen/logrus"
)

func (l *HTTPListener) homeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hello world!"))
	return
}

func (l *HTTPListener) registerHandler(w http.ResponseWriter, r *http.Request) {

	registerReq := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(registerReq); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not unmarshal request body")
		writeResponse(w, http.StatusBadRequest, false, []string{"could not unmarshal request body"}, nil)
		return
	}
	defer r.Body.Close()

	if err := l.validate.Struct(registerReq); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("request validation failed")
		writeResponse(w, http.StatusBadRequest, false, getValidationErrors(err), nil)
		return
	}

	var err error
	registerReq.PasswordDigest, err = hashPassword(*registerReq.Password)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, false, []string{"error creating new user"}, nil)
		return
	}
	registerReq.Password = nil

	if err := l.dbClient.CreateUser(registerReq); err != nil {
		writeResponse(w, http.StatusInternalServerError, false, []string{"error creating new user"}, nil)
		return
	}

	writeResponse(w, http.StatusOK, true, []string{}, map[string]string{"message": "user created successfully"})
	return
}
func (l *HTTPListener) loginHandler(w http.ResponseWriter, r *http.Request) {

	loginReq := models.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not unmarshal request body")
		writeResponse(w, http.StatusBadRequest, false, []string{"could not unmarshal request body"}, nil)
		return
	}
	defer r.Body.Close()

	if err := l.validate.Struct(loginReq); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("request validation failed")
		writeResponse(w, http.StatusBadRequest, false, getValidationErrors(err), nil)
		return
	}

	userFound, user, err := l.dbClient.GetUserByEmail(loginReq.Email)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, false, []string{"error authenticating user"}, nil)
		return
	}
	if !userFound {
		writeResponse(w, http.StatusNotFound, false, []string{fmt.Sprintf("no user exists corresponding to the email %s", loginReq.Email)}, nil)
		return
	}

	if !doPasswordsMatch(user.PasswordDigest, loginReq.Password) {
		writeResponse(w, http.StatusUnauthorized, false, []string{"the credentials sent are invalid"}, nil)
		return
	}

	token, err := l.loginUser(user.ID.String())
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, false, []string{"error logging in user"}, nil)
		return
	}
	writeResponse(w, http.StatusOK, true, []string{}, map[string]string{"message": "user logged in successfully", "token": token})
	return

}

func (l *HTTPListener) logoutHandler(w http.ResponseWriter, r *http.Request)         {}
func (l *HTTPListener) renewTokenHandler(w http.ResponseWriter, r *http.Request)     {}
func (l *HTTPListener) listUsersHandler(w http.ResponseWriter, r *http.Request)      {}
func (l *HTTPListener) getUserDetailsHandler(w http.ResponseWriter, r *http.Request) {}

func (l *HTTPListener) loginUser(id string) (string, error) {

	token, err := generateAuthToken(id)
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(models.Session{LoggedIn: true})
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not marshal session data")
		return "", err
	}
	if err := l.cacheClient.Write(id, string(jsonBytes)); err != nil {
		return "", err
	}

	return token, nil

}

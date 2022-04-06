package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (l *HTTPListener) newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", l.homeHandler)
	r.HandleFunc("/users/register", l.registerHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/confirm/sendmail", l.sendConfirmationMailHandler).Methods(http.MethodPut)
	r.HandleFunc("/users/confirm", l.confirmUserHandler).Methods(http.MethodPut)
	r.HandleFunc("/users/login", l.loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/logout", l.logoutHandler).Methods(http.MethodDelete)
	r.HandleFunc("/users/renewtoken", l.renewTokenHandler).Methods(http.MethodPut)
	r.HandleFunc("/users", l.listUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/users/:id", l.getUserDetailsHandler).Methods(http.MethodGet)

	return r
}

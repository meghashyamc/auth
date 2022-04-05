package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/meghashyamc/auth/services/cache"
	"github.com/meghashyamc/auth/services/db"
	log "github.com/sirupsen/logrus"
)

var (
	shutdownTime       = 5 * time.Second
	serverWriteTimeout = 60 * time.Second
	serverReadTimeout  = 60 * time.Second
)

type HTTPListener struct {
	dbClient    *db.DBClient
	cacheClient *cache.CacheClient
	server      *http.Server
	validate    *validator.Validate
}

func NewHTTPListener() (*HTTPListener, error) {

	dbClient, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	cacheClient, err := cache.NewClient()
	if err != nil {
		return nil, err
	}

	listener := &HTTPListener{dbClient: dbClient, cacheClient: cacheClient, validate: newValidator()}
	server := &http.Server{
		Handler:      listener.newRouter(),
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	listener.server = server
	return listener, nil

}

func (l *HTTPListener) Listen() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := l.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{"err": err.Error}).Error("HTTP listener exited with an error")
		}
	}()

	log.WithFields(log.Fields{"address": l.server.Addr}).Info("server started, listening successfully")
	signalReceived := <-done
	log.WithFields(log.Fields{"signal": signalReceived.String()}).Info("server stopped because of signal")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer func() {
		cancel()
	}()

	if err := l.server.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("server shutdown failed")
		return
	}
	log.Info("server exited gracefully")

}

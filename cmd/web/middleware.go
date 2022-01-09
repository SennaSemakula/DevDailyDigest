package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var warnLogger = log.New(os.Stdout, "WARN:\t", log.Ldate|log.Ltime)

func (app *App) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		// Code executed here will be down the chain e.g. middleware -> servemux -> handler

		// Handle api authentication
		if !authenticate(*app.bearer, w, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
		// Code executed after will be back up the chain
	})
}

func (app *App) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log every request
		app.infoLog.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// On panics during each request it will send a 500 to a user
func (app *App) handlePanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", rec))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func authenticate(token string, w http.ResponseWriter, r *http.Request) bool {
	currentToken := r.Header.Get("Authorization")
	validToken := fmt.Sprintf("Bearer %s", token)

	if currentToken != validToken {
		warnLogger.Println("Not authorised for endpoint", r.URL)
		return false
	}

	return true

}

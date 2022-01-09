package main

import (
	"net/http"
)

func (app *App) routes() http.Handler {
	// Create a local servemux to act as a router for our URL patterns
	mux := http.NewServeMux()

	// mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/api/v1/articles", app.articleHandler)
	mux.HandleFunc("/api/v1/hashnode_articles", app.hashnodeArticlesHandler)
	mux.HandleFunc("/api/v1/hackernoon_articles", app.hackernoonArticleHandler)
	mux.HandleFunc("/api/v1/users", app.usersHandler)
	mux.HandleFunc("/api/v1/users/create", app.usersCreate)
	mux.HandleFunc("/api/v1/users/unsubscribe", app.usersDelete)
	mux.HandleFunc("/api/v1/users/validate", app.validateUserHandler)
	mux.HandleFunc("/api/v1/jobs", app.jobHandler)
	mux.HandleFunc("/api/v1/events", app.eventHandler)
	mux.HandleFunc("/bundle", app.homeHandler)

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){w.Write([]byte("404: Could not find page"))})

	// buildHandler := http.FileServer(http.Dir("./ui/build/"))          // load the main index
	// staticHandler := http.FileServer(http.Dir("./ui/build/static/")) //serve the statics

	// mux.Handle("/", buildHandler)
	// mux.Handle("/static/", http.StripPrefix("/static", staticHandler))

	return app.logRequest(app.auth(app.handlePanics(mux)))
}

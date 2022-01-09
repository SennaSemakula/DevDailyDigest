package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	nl "github.com/Pioneersltd/DevDailyDigest/v1/pkg/newsletter"
)

type internalError interface {
	Error() string
}

// Helper function to write 500 error messages with stack trace
func (a *App) serverError(w http.ResponseWriter, err internalError) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errLog.Println(err.Error())
	a.errLog.Output(2, trace) // go one step back in the stack trace

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Helper function that sends an error status code to the user with message
func (a *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *App) notFoundError(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a *App) filterJobs(jobs []nl.Job, level string) []nl.Job {
	// Jobs in the hood
	var filteredJobs = make([]nl.Job, 0)
	for _, item := range jobs {
		if item.Level == level {
			filteredJobs = append(filteredJobs, item)
		}
	}

	return filteredJobs
}

package main

import (
	"github.com/Pioneersltd/DevDailyDigest/v1/pkg/client"
	"github.com/Pioneersltd/DevDailyDigest/v1/pkg/models/users"
	nl "github.com/Pioneersltd/DevDailyDigest/v1/pkg/newsletter"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	userList         = make(map[string]users.User)
	articles         []nl.Article
	jobs             []nl.Job
	jobsMap          = make(map[string]map[string][]nl.Job)
	countryMap       = make(map[string]map[string]map[string][]nl.Job)
	events           []nl.Event
	eventsMap        = make(map[string]map[string][]nl.Event)
	hashnodeArticles = make(map[string][]nl.Article)
	hackernoonArticles = make(map[string][]nl.Article)
)

func (a *App) usersCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if method := r.Method; method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // set header to show client what method is allowed
		a.clientError(w, 405)
		return
	}

	switch method := r.Method; method {
	case "GET":
		resp, err := json.Marshal(userList)
		if err != nil {
			a.serverError(w, err)
			return
		}
		w.Write(resp)
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		newUser := users.User{}
		newUser.Lock()

		defer newUser.Unlock()

		if err != nil {
			a.serverError(w, err)
			return
		}

		if err := json.Unmarshal([]byte(body), &newUser); err != nil {
			a.clientError(w, 400)
			return
		} else {
			// TODO: check for optional data e.g. resume
			lang := strings.ToLower(newUser.Language)

			entry, err := a.model.Insert(newUser.Email, newUser.Country, newUser.City,
				lang, newUser.Level, newUser.Frameworks, newUser.Resume,
				newUser.Git, newUser.Date, newUser.JobStatus, newUser.Referral,
				newUser.Recruit)

			if err != nil {
				a.serverError(w, err)
				return
			}

			// User exists; updating
			if entry == 2 {
				a.infoLog.Printf("Updating user %v in Database", newUser.Email)
				w.WriteHeader(200)
			} else {
				a.infoLog.Printf("Saving user %v to Database\n", newUser.Email)
				w.WriteHeader(201)
				// Send welcome email
				err = client.SendWelcome(newUser.Email, lang, newUser.Referral)
			}

			userList[newUser.Email] = newUser

			if err != nil {
				a.errLog.Println("unable to send welcome email to %v", newUser.Email)
			}

		}
	}
}

func (a *App) usersDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if method := r.Method; method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // set header to show client what method is allowed
		a.clientError(w, 405)
		return
	}

	switch method := r.Method; method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		newUser := users.User{}
		newUser.Lock()

		defer newUser.Unlock()

		if err != nil {
			a.serverError(w, err)
			return
		}

		if err := json.Unmarshal([]byte(body), &newUser); err != nil {
			a.clientError(w, 400)
			return
		} else {
			a.infoLog.Printf("Unsubscribing user %v\n", newUser.Email)
			if err := a.model.Delete(newUser.Email); err != nil {
				a.serverError(w, err)
			}

			w.WriteHeader(204)
		}
	}
}

func (a *App) usersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := a.model.GetAll()

	if err != nil {
		a.serverError(w, err)
		return
	}
	resp, err := json.Marshal(users)

	if err != nil {
		a.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(resp)
}

func (a *App) validateUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp string
	email := r.URL.Query().Get("email")
	user, err := a.model.Get(email)

	w.Header().Set("Content-Type", "application/json")

	if user != nil {
		a.warnLog.Printf("user %v already exists in database", email)
		resp = fmt.Sprintf(`{"status": true}`)
		w.Write([]byte(resp))
		return
	}

	// Tidy and and use errors.Is instead
	if err.Error() == "no matching record found" {
		resp = fmt.Sprintf(`{"status": false}`)
		w.Write([]byte(resp))
		return
	}

	if err != nil {
		a.serverError(w, err)
		return
	}
}

// Serve the articles
func (a *App) articleHandler(w http.ResponseWriter, r *http.Request) {
	if method := r.Method; method != "GET" {
		w.Header().Set("Allow", "GET")
		a.clientError(w, 405)
		return
	}
	// TODO: santise input by checking if string is returned
	if len(r.URL.Query().Get("language")) > 0 {
		// Using RawQuery to fetch escaped characters e.g. c++
		language := r.URL.RawQuery[9:]
		if language == "c++" {
			language = "cpp"
		}

		articles := client.Fetch("https://dev.to/api/articles", language)

		w.Header().Set("Content-Type", "application/json") // must set content type because http.DetectContentType does
		//not understand JSON so gives plain text
		resp, err := json.Marshal(articles)

		if err != nil {
			a.serverError(w, err)
			return

		}
		w.Write(resp)
	} else {
		w.Write([]byte("No articles"))
	}
}

// Serve the jobs
func (a *App) jobHandler(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
	case "GET":
		country := r.URL.Query().Get("country")
		city := r.URL.Query().Get("city")
		language := r.URL.Query().Get("language")
		level := r.URL.Query().Get("level")

		if strings.Contains(language, "c") && language != "javascript" {
			language = "c++"
		}

		w.Header().Set("Content-Type", "application/json")

		// TODO: refactor
		if len(r.URL.Query()) > 0 {
			if _, ok := countryMap[country][city]; ok {
				// TODO: CLEAN THIS UP
				if len(level) > 0 {
					filteredJobs := a.filterJobs(countryMap[country][city][language], level)
					resp, err := json.Marshal(filteredJobs)

					if err != nil {
						a.errLog.Println(err.Error())
						a.serverError(w, err)
						return
					}

					if len(filteredJobs) > 3 {
						w.Write(resp)
						return
					}
				}
				resp, err := json.Marshal(countryMap[country][city][language])
				if err != nil {
					a.errLog.Println(err.Error())
					a.serverError(w, err)
					return
				}
				w.Write(resp)
				return
			} else {
				if len(level) > 0 {
					filteredJobs := a.filterJobs(countryMap["remote"]["remote"][language], level)
					resp, err := json.Marshal(filteredJobs)

					if err != nil {
						a.errLog.Println(err.Error())
						a.serverError(w, err)
						return
					}

					if len(filteredJobs) > 3 {
						w.Write(resp)
						return
					}
				}
				resp, err := json.Marshal(countryMap["remote"]["remote"][language])
				if err != nil {
					a.errLog.Println(err.Error())
					a.serverError(w, err)
					return
				}
				w.Write(resp)
				return
			}

		}
		resp, err := json.Marshal(countryMap)

		if err != nil {
			a.errLog.Println(err.Error())
			a.serverError(w, err)
			return
		}

		w.Write(resp)
	case "POST":
		var (
			country  = r.Header.Get("country")
			city     = r.Header.Get("city")
			language = r.Header.Get("category")
		)
		var newJobs []nl.Job
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			a.serverError(w, err)
			return
		}

		w.Header().Set("allow", http.MethodPost)

		// Handle ALL JOBS POST
		if len(city) == 0 {
			var newAllJobs = make(map[string]map[string]map[string][]nl.Job)
			if err := json.Unmarshal([]byte(body), &newAllJobs); err != nil {
				log.Println(err)
				a.clientError(w, 400)
				return
			} else {
				// Concatenate by unpacking a slice
				countryMap = newAllJobs

				log.Printf("Added ALL jobs to API\n")
	
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		if err := json.Unmarshal([]byte(body), &newJobs); err != nil {
			a.clientError(w, 405)
			return
		} else {
			log.Printf("Adding (%v, %v, %v) jobs to API\n", language, country, city)
			// Concatenate by unpacking a slice
			jobs = append(jobs, newJobs...)

			cityMap := make(map[string]map[string][]nl.Job)
			languageMap := make(map[string][]nl.Job)

			if len(countryMap[country][city]) > 0 {
				countryMap[country][city][language] = newJobs
			} else if len(countryMap[country]) > 0 {
				countryMap[country][city] = languageMap
				languageMap[language] = newJobs
			} else {
				countryMap[country] = cityMap
				cityMap[city] = languageMap
				languageMap[language] = newJobs
			}

			w.WriteHeader(http.StatusOK)
		}
	}
}

func (a *App) hashnodeArticlesHandler(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
	case "GET":
		language := r.URL.Query().Get("language")

		if strings.Contains(language, "c") && language != "javascript" {
			language = "c++"
		}

		w.Header().Set("Content-Type", "application/json")

		// TODO: refactor
		resp, err := json.Marshal(hashnodeArticles[language])
		if err != nil {
			a.errLog.Println(err.Error())
			a.serverError(w, err)
			return
		}
		w.Write(resp)

	case "POST":
		var (
			language = r.Header.Get("category")
		)
		var newArticles []nl.Article
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			a.serverError(w, err)
			return
		}

		w.Header().Set("allow", http.MethodPost)
		if err := json.Unmarshal([]byte(body), &newArticles); err != nil {
			a.clientError(w, 405)
			return
		} else {
			log.Printf("Adding (%v) hashnode articles to API\n", language)

			hashnodeArticles[language] = newArticles

			w.WriteHeader(http.StatusOK)
		}
	}
}

func (a *App) hackernoonArticleHandler(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
	case "GET":
		language := r.URL.Query().Get("language")

		if strings.Contains(language, "c") && language != "javascript" {
			language = "c++"
		}

		w.Header().Set("Content-Type", "application/json")

		// TODO: refactor
		resp, err := json.Marshal(hackernoonArticles[language])
		if err != nil {
			a.errLog.Println(err.Error())
			a.serverError(w, err)
			return
		}
		w.Write(resp)

	case "POST":
		var newArticles = make(map[string][]nl.Article)
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			a.serverError(w, err)
			return
		}

		w.Header().Set("allow", http.MethodPost)
		if err := json.Unmarshal([]byte(body), &newArticles); err != nil {
			a.clientError(w, 405)
			return
		} else {
			log.Printf("Adding hackernoon articles to API\n")

			hackernoonArticles = newArticles

			w.WriteHeader(http.StatusOK)
		}
	}
}

// Serve the events
func (a *App) eventHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
		city := r.URL.Query().Get("city")
		category := r.URL.Query().Get("category")

		if strings.Contains(category, "c") && category != "javascript" {
			category = "c++"
		}

		w.Header().Set("Content-Type", "application/json")

		if len(r.URL.Query()) > 0 {
			resp, err := json.Marshal(eventsMap[city][category])
			if err != nil {
				a.errLog.Println(err.Error())
				a.serverError(w, err)
				return
			}
			w.Write(resp)
		} else {
			resp, err := json.Marshal(eventsMap)
			if err != nil {
				a.errLog.Println(err.Error())
				a.serverError(w, err)
				return
			}
			w.Write(resp)
		}

	case "POST":
		var newEvents []nl.Event
		body, err := ioutil.ReadAll(r.Body)

		var (
			city     = r.Header.Get("city")
			language = r.Header.Get("category")
		)

		if err != nil {
			a.serverError(w, err)
			return
		}

		w.Header().Set("allow", http.MethodPost)

		// POST ALL EVENTS
		if len(city) == 0 {
			var newAllEvents = make(map[string]map[string][]nl.Event)
			if err := json.Unmarshal([]byte(body), &newAllEvents); err != nil {
				log.Println(err)
				a.clientError(w, 400)
				return
			} else {
				// Concatenate by unpacking a slice
				eventsMap = newAllEvents

				log.Printf("Added all events to API\n")
	
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		if err := json.Unmarshal([]byte(body), &newEvents); err != nil {
			log.Println(err)
			a.clientError(w, 400)
			return
		} else {
			// Concatenate by unpacking a slice
			events = append(events, newEvents...)

			languageMap := make(map[string][]nl.Event)

			if len(eventsMap[city]) > 0 {
				eventsMap[city][language] = newEvents
			} else {
				eventsMap[city] = languageMap
				languageMap[language] = newEvents
			}

			log.Printf("Added (%v, %v) events to API\n", language, city)

			w.WriteHeader(http.StatusOK)
		}
	}
}

func (a *App) bundleHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/build/index.html")
}

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	a.notFoundError(w)
	// 	return
	// }

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.page.tmpl",
		"./ui/html/footer.page.tmpl",
	}

	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		a.serverError(w, err)
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		a.serverError(w, err)
	}
}

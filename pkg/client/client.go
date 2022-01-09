package client

import (
	"bytes"
	nl "github.com/Pioneersltd/DevDailyDigest/v1/pkg/newsletter"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	start      time.Time
	url        = "https://dev.to/api/articles"
	emailApi   = "https://api.mailjet.com/v3.1/send"
	api        = "https://devdaily.io"
	authToken  = os.Getenv("DEV_AUTH_TOKEN")
	bearer     = os.Getenv("BEARER_TOKEN")
	client     *http.Client
	infoLogger = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)
	errLogger  = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	start = time.Now()
	client = &http.Client{}
}

type ArticlesMap map[string]interface{}
type usersMap map[string]interface{}

type Articles struct {
	articles       []ArticlesMap
	storedArticles []ArticlesMap
}

type Users struct {
	records []usersMap
}

func FetchUsers() ([]usersMap, error) {
	url := fmt.Sprintf("%v/api/v1/users", api)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bearer))
	resp, err := client.Do(req)

	var users Users

	if err != nil {
		errLogger.Println("unable to fetch users")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &users.records); err != nil {
		return nil, err
	}

	return users.records, nil
}

func FetchArticles(path, tag string) ([]ArticlesMap, error) {
	url := fmt.Sprintf("%v/%v?language=%v", api, path, tag)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bearer))
	resp, err := client.Do(req)

	var articles Articles

	if err != nil {
		errLogger.Printf("unable to fetch %v articles", tag)
	}

	if err := checkHTTPStatus(resp.StatusCode); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &articles.storedArticles); err != nil {
		return nil, err
	}

	return articles.storedArticles, nil
}

func FetchJobs(country, city, language string) ([]nl.Job, error) {
	country = strings.ReplaceAll(country, " ", "+")
	city = strings.ReplaceAll(city, " ", "+")
	url := fmt.Sprintf("%v/api/v1/jobs?country=%v&city=%v&language=%v", api, country, city, language)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bearer))
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if err := checkHTTPStatus(resp.StatusCode); err != nil {
		return nil, err
	}

	var jobs []nl.Job

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}


func FetchEvents(city, language string) ([]nl.Event, error) {
	city = strings.ReplaceAll(city, " ", "+")
	url := fmt.Sprintf("%v/api/v1/events?category=%v&city=%v", api, language, city)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bearer))
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var events []nl.Event

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := checkHTTPStatus(resp.StatusCode); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &events); err != nil {
		return nil, err
	}

	return events, nil
}
func Fetch(url string, tag string) []ArticlesMap {
	// TODO: get the right query parameters to always get fresh articles
	api := url + "?tag=" + tag + "&per_page=200&state=fresh&top=20"

	req, err := http.NewRequest("GET", api+"", nil)
	req.Header.Add("api-key", authToken)
	resp, err := client.Do(req)

	if err != nil {
		resp.Body.Close()
		errLogger.Printf("unable to fetch %v articles", tag)
	}

	if err := checkHTTPStatus(resp.StatusCode); err != nil {
		errLogger.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		errLogger.Println("unable to fetch users")
	}

	var articles Articles
	if err := json.Unmarshal(body, &articles.articles); err != nil {
		errLogger.Println(err)
	}
	return articles.articles
}

func SendWelcome(email, language, referral string) error {
	now := time.Now()
	currentDate := fmt.Sprintf("%v, %v %v %v", now.Weekday(), now.Day(), now.Month(), now.Year())

	language = strings.Split(language, ",")[0]

	body := []byte(fmt.Sprintf(`
	{
		"Messages":[
		  {
			"From": {
			  "Email": "no-reply@devdaily.io",
			  "Name": "DevDaily"
			},
			"To": [
			  {
				"Email": %q,
				"Name": ""
			  }
			],
			"TemplateID": 2051628,
			"TemplateLanguage": true,
			"Subject": "Welcome to DevDaily",
			"Variables": {
		"date": %q,
		"language": %q,
		"name": "DevDaily",
		"referral": %q,
		"sender_name": "DevDaily",
		"sender_address": "12th Brixton Avenue",
		"sender_city": "London",
		"sender_state": "United Kingdom",
		"sender_postcode": "SW2 32H",
		"unsubscribe": "https://devdaily.io/unsubscribe",
		"unsubscribe_preferences": "https://devdaily.io",
		"contact_us": "https://devdaily.io"
		}
		  }
		]
	  }
	`, email, currentDate, language, referral))
	req, err := http.NewRequest("POST", emailApi, bytes.NewBuffer(body))
	req.SetBasicAuth("e00beea1bb1882667308052232bc6a37", "3d9905d71a213c8649fce864e3518973")
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		errLogger.Printf("unable to send %v welcome email to %v", language, email)
		return err
	}

	infoLogger.Printf("sent welcome email to %v", email)

	return nil
}

func checkHTTPStatus(status int) error{
	if status < 200 || status > 299 {
		return fmt.Errorf("unable to fetch articles. http status code: %v", status)
	}
	
	return nil
}


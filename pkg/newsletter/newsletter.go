package newsletter

import (
	"html/template"
	"log"
)

type Template template.Template

type Location struct {
	Lat, Long float32
}

type Job struct {
	Role, Company string
	Level         string
	Salary        string
	Location      string
	Url           string
	Date          string
}

type Article struct {
	Title, Desc, Date, Image, Url, Category string
}

type Question struct {
	Question string
}

type Event struct {
	Title     string
	Group     string
	Attendees string
	Date      string // Change to more appropiate data type
	Time      string
	Location  string
	Url       string
	Image     string
}

type Letter struct {
	Questions []Question
	Jobs      []Job
	Articles  []Article
	Events    []Event
}

type File struct {
	name     string
	contents []byte
}

func LoadTemplate(tmpl string) template.Template {
	t, err := template.ParseFiles(tmpl)

	if err != nil {
		log.Fatal(err)
	}

	return *t
}

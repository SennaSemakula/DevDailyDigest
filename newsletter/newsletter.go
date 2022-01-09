package newsletter

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

type location struct {
	lat, long float32
}

type Job struct {
	Role, Company string
	Salary        string
	Url           string
}

type Article struct {
	Title, Desc, Date, Image, Url string
}

type Question struct {
	Question string
}

type Event struct {
	Name string
	Loc  location
	Date string
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

func loadTemplate(tmpl string) *template.Template {
	t, err := template.ParseFiles(tmpl)

	if err != nil {
		log.Fatal(err)
	}

	return t
}

func main() {
	tmpl := loadTemplate("templates/index.html")

	questions := []Question{
		Question{Question: "What is a channel in Go?"},
	}

	articles := []Article{
		Article{
			Title: "Test article",
			Desc:  "My Description",
			Date:  "27/09/20",
			Image: "Test Image",
			Url:   "https://www.google.co.uk",
		},
	}

	jobs := []Job{
		Job{Role: "Golang developer", Company: "CodeBundle", Salary: "1241232.7"},
	}

	events := []Event{
		Event{Name: "Python event", Loc: location{48.4284, 123.3656}, Date: "27/09/20"},
	}

	data := Letter{Questions: questions, Articles: articles, Jobs: jobs, Events: events}

	letter := tmpl.Execute(os.Stdout, data)
	fmt.Println(letter)
}

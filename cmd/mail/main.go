package main

import (
	"github.com/Pioneersltd/DevDailyDigest/v1/cmd/mail/digest"
	"github.com/Pioneersltd/DevDailyDigest/v1/pkg/client"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	emailUser   = os.Getenv("EMAIL_USER")
	emailToken  = os.Getenv("EMAIL_TOKEN")
	bearerToken = os.Getenv("BEARER_TOKEN")
	errLogger   = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
)

const (
	templateId = "2052202"
)

type config struct {
	mode     string
	email    string
	level    string
	city     string
	country  string
	language string
}

type emailStatus struct {
	email  string
	status int
}

func init() {
	if emailUser == "" || emailToken == "" {
		log.Fatal("unable to find MailJet environment credentials. Are EMAIL_USER and EMAIL_TOKEN variables exported?")
	}

	if bearerToken == "" {
		log.Fatal("unable to find bearer token. Is BEARER_TOKEN variables exported?")
	}
}

func sendEmails(c *config) {

	if c.mode == "dev" {
		emailBody := &digest.Email{c.email, c.country, c.city, c.language, c.level, ""}
		// log.Fatal("please use a mode from the following options: 'dev' | 'production'")
		status, err := digest.SendMail(emailBody, templateId)

		if err != nil {
			errLogger.Println(err.Error())
		}

		if status != 200 {
			fmt.Printf("need to retry sending email to %v\n", c.email)
		}
		return
	}

	var input string

	if c.mode == "integration" {
		users, _ := client.FetchUsers()
		fmt.Printf("Distributing %v DevDaily emails...\n", len(users))

		for i, user := range users {
			emailBody := &digest.Email{user["Email"].(string), user["Country"].(string), user["City"].(string), user["Language"].(string), user["Level"].(string), user["Frameworks"].(string)}
			if emailBody.Email == "senna.semakula@gmail.com" {
				status, err := digest.SendMail(emailBody, templateId)

				if err != nil {
					errLogger.Println(err.Error())
				}

				if status != 200 {
					fmt.Printf("need to retry sending email to %v", emailBody.Email)
				} else {
					fmt.Println("number ", i, " just got sent")
				}
			}
		}
	}

	if c.mode == "production" {
		fmt.Println("Are you sure you want to distribute emails to all DevDaily customers? (NO/YES)")
		fmt.Scanln(&input)
		if input != "YES" {
			return
		}
		users, _ := client.FetchUsers()
		fmt.Printf("Distributing %v DevDaily emails...\n", len(users))

		for i, user := range users {
			emailBody := &digest.Email{user["Email"].(string), user["Country"].(string), user["City"].(string), user["Language"].(string), user["Level"].(string), user["Frameworks"].(string)}
			status, err := digest.SendMail(emailBody, templateId)

			if err != nil {
				errLogger.Println(err.Error())
			}

			if status != 200 {
				fmt.Printf("need to retry sending email to %v", emailBody.Email)
			} else {
				fmt.Println("number ", i, " just got sent")
			}

		}
	}

}

func main() {

	mode := flag.String("mode", "", "use 'dev' or 'production'. Only use production to send out emails to customers")
	email := flag.String("email", "", "email address to send digest to e.g. 'apple@gmail.com'")
	level := flag.String("level", "", "level e.g. beginner, mid, senior")
	city := flag.String("city", "", "e.g. London")
	country := flag.String("country", "", "e.g. United Kingdom")
	language := flag.String("language", "", "programming lanuage. Choices: javascript, python, golang, ruby, php, c++")
	flag.Parse()

	setup := config{*mode, *email, *level, *city, *country, *language}

	sendEmails(&setup)

	// TODO: implement concurrency

}

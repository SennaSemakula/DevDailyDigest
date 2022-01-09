/*
user package that represents Users and Recruiters

User:

`
{
	"Email": "test@gmail.com",
	"Location": "london",
	"Git": "https://github.com",
	"Interests": {
		"Coding": {
			"Go": true,
			"JavaScript": true,
			"Python": true
		},
		"Bundle": {
			"Job": true,
			"Article": true,
			"Event": false
		}
	},
	"Resume": "senna_resume"
}
`

*/
package users

import "sync"

type User struct {
	Email      string
	Country    string
	City       string
	Git        string
	Interests  Bundle
	Language   string
	Level      string
	Frameworks string
	Resume     string
	Recruit    bool
	Date       string
	JobStatus  string
	Referral   string
	Mu         sync.RWMutex
}

func (u *User) Lock() {
	u.Mu.Lock()
}

func (u *User) Unlock() {
	u.Mu.Unlock()
}

type Bundle struct {
	Coding  Code
	Content Content
}

type Code struct {
	Language   string
	Frameworks string
	Level      string
}

type Content struct {
	Job     bool
	Article bool
	Event   bool
}

type Recruiter struct {
	Name  string
	Email string
}

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

type User struct {
	Email     string
	Location  string
	Git       string
	Interests Bundle
	Resume    string
}

type Bundle struct {
	Coding  Code
	Content Content
}

type Code struct {
	Go         bool
	Python     bool
	JavaScript bool
}

type Content struct {
	Job bool
	Article bool
	Event bool
}

type Recruiter struct {
	Name  string
	Email string
}

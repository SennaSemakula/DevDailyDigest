package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

// TODO: refactor language/frameworks into different struct
type User struct {
	Email      string
	Country    string
	City       string
	Git        string
	Resume     string
	Recruit    bool
	Language   string
	Level      string
	Frameworks string
	Date       string
	JobStatus  string
	Referral string
}

type Language struct {
	JavaScript bool
	Python     bool
	Go         bool
}

type Interest struct {
	Event   bool
	Job     bool
	Article bool
}

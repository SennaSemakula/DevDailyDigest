package email

import (
	"log"
	"net/smtp"
	"os"
	"strings"
)

type Client struct {
	Username string
	Token    string
	Server   string
	auth     smtp.Auth
}

type Mail struct {
	From    string
	To      []string
	Subject string
	Body    []byte
}

var creds = make(map[string]string)

// Initialise a client with def
func NewClient() *Client {
	newClient := &Client{creds["EMAIL_USER"], creds["EMAIL_TOKEN"], creds["EMAIL_SERVER"], nil}
	return newClient
}

// Receiver for Client struct to authenticate smtp Client
func (c *Client) Authenticate() {
	c.auth = smtp.PlainAuth("", c.Username, c.Token, c.Server)
}

// Create a reciever to get mail decode the body
func (m *Mail) GetBody(body []byte) string {
	msg := string(body)

	return msg
}

// Receiver for Client func to send emails
func (c *Client) Send(m *Mail) error {
	server := c.Server + ":" + creds["EMAIL_PORT"]
	err := smtp.SendMail(server, c.auth, c.Username, m.To, m.Body)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sending email to", m.To)

	return nil
}

func init() {
	// Intialise email Client with credentials
	log.Println("Initialising email credentials")
	for _, e := range os.Environ() {
		if e := e; e[:5] == "EMAIL" {
			k := strings.Split(e, "=")

			creds[k[0]] = os.Getenv(k[0])
		}
	}

	// Terminate if sufficient credentials are not provided
	if len(creds) != 4 {
		log.Fatal("Do not have sufficient email environment variables set.\r\nRequires EMAIL_USER, EMAIL_TOKEN, EMAIL_SERVER, EMAIL_PORT")
	}
}

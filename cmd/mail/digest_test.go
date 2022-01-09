package digest

import (
	"github.com/Pioneersltd/DevDailyDigest/v1/cmd/mail/digest"
	"github.com/Pioneersltd/DevDailyDigest/v1/pkg/client"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"os"
	"testing"
	"strings"
)

// TODO: change mailer.go package to be mailer
// TODO: Move testing of articles into a different package
// Create a main.go and call mailer.SendEmail()

func init() {
	if os.Getenv("BEARER_TOKEN") == "" {
		panic("environment variable not exported. Please export environment variable 'BEARER_TOKEN'")
	}
}

func TestFetchDevToArticles(t *testing.T) {
	expected := 3

	languages := []string{"javascript", "python", "golang", "ruby", "php", "java", "c++"}

	for _, lang := range languages {
		t.Run(lang+" Articles", func(t *testing.T) {
			articles, err := digest.GetArticles(lang, "", "")

			if err != nil {
				t.Fatalf(`FetchArticles("api/v1/articles", %q) got err: %v`, lang, err)
			}

			actual := len(articles)

			if actual < expected {
				t.Fatalf(`FetchArticles("api/v1/articles", %q) got length %v but want greater than %v`, lang, actual, expected)
			}
		})
	}

	frameworks := []string{"react", "node", "angular", "vue", "django", "Tornado", "flask"}
	for _, val := range frameworks {
		t.Run(val+" Articles", func(t *testing.T) {
			lang := val
			articles, err := digest.GetArticles("javascript", val, "")

			if err != nil {
				t.Fatalf(`FetchArticles("api/v1/articles", %q) got err: %v`, lang, err)
			}

			actual := len(articles)

			if actual < expected {
				t.Fatalf(`FetchArticles("api/v1/articles", %q) got length %v but want greater than %v`, lang, actual, expected)
			}
		})
	}
}

func TestBeginnerArticles(t *testing.T) {
	expected := 3

	languages := []string{"javascript", "python", "golang"}

	for _, lang := range languages {
		t.Run(lang+" Beginner Articles", func(t *testing.T) {
			articles, err := digest.GetBeginnerArticles(lang)

			if err != nil {
				t.Fatalf(`GetBeginnerArticles(%q) got err: %v`, lang, err)
			}

			actual := len(articles)

			if actual < expected {
				t.Fatalf(`GetBeginnerArticles(%q) got length %v but want greater than %v`, lang, actual, expected)
			}
		})
	}

	invalidLanguages := []string{"java", "ruby", "c++", "ruby", "php"}
	expected = 0

	for _, lang := range invalidLanguages {
		t.Run(lang+" Beginner Articles", func(t *testing.T) {
			articles, err := digest.GetBeginnerArticles(lang)

			if err != nil {
				t.Fatalf(`GetBeginnerArticles(%q) got err: %v`, lang, err)
			}

			actual := len(articles)

			if actual < expected {
				t.Fatalf(`GetBeginnerArticles(%q) got length %v but want  %v`, lang, actual, expected)
			}
		})
	}
}

func TestFetchEvents(t *testing.T) {
	expected := 2

	cases := map[string]string{
		"London":       "javascript",
		"Kyoto":        "python",
		"Birmingham":   "golang",
		"Paris":        "java",
		"Lagos":        "ruby",
		"Glasgow":      "php",
		"Johannesburg": "c++",
		"remote":       "javascript",
	}

	for key, val := range cases {
		t.Run(fmt.Sprintf("Fetch %v Events in London", val), func(t *testing.T) {

			events, err := digest.FetchEvents(key, val)
			actual := len(events)

			if err != nil {
				t.Fatalf(`FetchEvents(%q, %q) returned %v`, key, val, err)
			}

			if actual < expected {
				t.Fatalf(`FetchEvents(%q, %q) got length %v but want greater than %v`, key, val, actual, expected)
			}
		})
	}
}

func TestFetchHackernoonArticles(t *testing.T) {
	expected := 3

	languages := []string{"javascript", "python", "golang", "ruby", "php", "java", "c++"}

	for _, lang := range languages {
		t.Run(lang+" Articles", func(t *testing.T) {
			articles, err := client.FetchArticles("api/v1/hackernoon_articles", lang)

			if err != nil {
				t.Fatalf(`FetchArticles("api/v1/hackernoon_articles", %q) got err: %v`, lang, err)
			}

			actual := len(articles)

			if actual < expected {
				t.Fatalf(`FetchArticles("api/v1/hackernoon_articles", %q) got length %v but want greater than %v`, lang, actual, expected)
			}
		})
	}
}

// // func TestFetchMiscToArticles(t *testing.T) {
// // 	expected := 3

// // }

// TestGetJobs tests FetchJobs
func TestFetchJobs(t *testing.T) {
	expected := 3

	t.Run("JavaScript London Jobs", func(t *testing.T) {
		jobs, err := client.FetchJobs("United Kingdom", "London", "javascript")

		if err != nil {
			t.Fatalf(`FetchJobs("United Kingdom", "London", "javascript") got %q`, err)
		}

		actual := len(jobs)

		if len(jobs) < 3 {
			t.Fatalf(`FetchJobs("United Kingdom", "London", "javascript") got length %v but want greater than %v`, actual, expected)
		}

	})

	t.Run("Python Manchester Jobs", func(t *testing.T) {
		jobs, err := client.FetchJobs("United Kingdom", "Manchester", "python")

		if err != nil {
			t.Fatalf(`FetchJobs("United Kingdom", "Manchester", "python") got %q`, err)
		}

		actual := len(jobs)

		if len(jobs) < 3 {
			t.Fatalf(`FetchJobs("United Kingdom", "Manchester", "python") got length %v but want greater than %v`, actual, expected)
		}

	})

	t.Run("Golang Liverpool Jobs", func(t *testing.T) {
		jobs, err := client.FetchJobs("United Kingdom", "Liverpool", "golang")

		if err != nil {
			t.Fatalf(`FetchJobs("United Kingdom", "Liverpool", "golang") got %q`, err)
		}

		actual := len(jobs)

		if len(jobs) < 3 {
			t.Fatalf(`FetchJobs("United Kingdom", "Liverpool", "golang") got length %v but want greater than %v`, actual, expected)
		}

	})

	t.Run("Java Reading Jobs", func(t *testing.T) {
		lang := "java"
		jobs, err := client.FetchJobs("United Kingdom", "Hounslow", lang)

		if err != nil {
			t.Fatalf(`FetchJobs("United Kingdom", "Hounslow", %q) got %q`, lang, err)
		}

		actual := len(jobs)

		if len(jobs) < 3 {
			t.Fatalf(`FetchJobs("United Kingdom", "Hounslow", %q) got length %v but want greater than %v`, lang, actual, expected)
		}

	})

	t.Run("Ruby Dallas Jobs", func(t *testing.T) {
		lang := "ruby"
		jobs, err := client.FetchJobs("United States", "Dallas", lang)

		if err != nil {
			t.Fatalf(`FetchJobs("United States", "Dallas", %q) got %q`, lang, err)
		}

		actual := len(jobs)

		if len(jobs) < 3 {
			t.Fatalf(`FetchJobs("United States", "Dallas", %q) got length %v but want greater than %v`, lang, actual, expected)
		}

	})

	t.Run("PHP Toronto Jobs", func(t *testing.T) {
		lang := "php"
		jobs, err := client.FetchJobs("Canada", "Toronto", lang)

		if err != nil {
			t.Fatalf(`FetchJobs("Canada", "Toronto", %q) got %q`, lang, err)
		}

		actual := len(jobs)

		if len(jobs) < 3 {
			t.Fatalf(`FetchJobs("Canada", "Toronto", %q) got length %v but want greater than %v`, lang, actual, expected)
		}

	})

	t.Run("C++ Sydney Jobs", func(t *testing.T) {
		lang := "c++"
		jobs, err := client.FetchJobs("Australia", "Sydney", lang)

		if err != nil {
			t.Fatalf(`FetchJobs("Australia", "Sydney", %q) got %q`, lang, err)
		}

		actual := len(jobs)

		if len(jobs) < 3 {
			t.Fatalf(`FetchJobs("Australia", "Sydney", %q) got length %v but want greater than %v`, lang, actual, expected)
		}

	})

	t.Run("JavaScript Remote Jobs", func(t *testing.T) {
		lang := "javascript"
		jobs, err := client.FetchJobs("remote", "remote", lang)

		if err != nil {
			t.Fatalf(`FetchJobs("remote", "remote", %q) got %q`, lang, err)
		}

		actual := len(jobs)

		if len(jobs) < 3 {
			t.Fatalf(`FetchJobs("remote", "remote", %q) got length %v but want greater than %v`, lang, actual, expected)
		}

	})

}

// TODO: seperate these unit tests out into a different package

func TestBuildBodySingle(t *testing.T) {
	const (
		templateId = "fakeId"
		to         = "fakeemail@com"
		country    = "United Kingdom"
		city       = "London"
		language   = "javascript"
		level      = "mid"
		frameworks = ""
	)
	emailBody, err := digest.BuildBody(to, templateId, country, city, language, level, frameworks)

	if err != nil {
		t.Fatalf(`buildBody(%q, %q, %q, %q, %q, %q, %q) got err %q`, to, templateId, country, city, language, level, frameworks, err)
	}

	t.Run("Test Email Subject", func(t *testing.T) {
		expected := fmt.Sprintf("☕ Top %v Resources You Have To Read", strings.Title(language))
		actual, err := jsonparser.GetString(emailBody, "Messages", "[0]", "Subject")

		if err != nil {
			t.Fatalf("unable to get json string. Got %q", err)
		}

		if actual != expected {
			t.Fatalf("got %q but expected %q for email subject", actual, expected)
		}
	})

	t.Run("Test Email Articles length", func(t *testing.T) {
		articles := make([]map[string]string, 0)
		expected := 5

		emailArticles, _, _, err := jsonparser.Get(emailBody, "Messages", "[0]", "Variables", "articles")

		if err != nil {
			t.Fatalf(`Unable to parse email articles json. got %q`, err)
		}

		if err := json.Unmarshal(emailArticles, &articles); err != nil {
			t.Fatalf(`Unable to unmarshal email articles. got %q`, err)
		}

		actual := len(articles)

		if actual < expected {
			t.Fatalf(`email articles length got %q but expected %q`, actual, expected)
		}

	})

	t.Run("Test Email Jobs length", func(t *testing.T) {
		jobs := make([]map[string]string, 0)
		expected := 3

		emailJobs, _, _, err := jsonparser.Get(emailBody, "Messages", "[0]", "Variables", "jobs")

		if err != nil {
			t.Fatalf(`Unable to parse email jobs json. got %q`, err)
		}

		if err := json.Unmarshal(emailJobs, &jobs); err != nil {
			t.Fatalf(`Unable to unmarshal email jobs. got %q`, err)
		}

		actual := len(jobs)

		if actual < expected {
			t.Fatalf(`email jobs length got %q but expected %q`, actual, expected)
		}

	})

	t.Run("Test Email Events length", func(t *testing.T) {
		events := make([]map[string]string, 0)
		expected := 2

		emailEvents, _, _, err := jsonparser.Get(emailBody, "Messages", "[0]", "Variables", "events")

		if err != nil {
			t.Fatalf(`Unable to parse email events json. got %q`, err)
		}

		if err := json.Unmarshal(emailEvents, &events); err != nil {
			t.Fatalf(`Unable to unmarshal email events. got %q`, err)
		}

		actual := len(events)

		if actual < expected {
			t.Fatalf(`email events length got %q but expected %q`, actual, expected)
		}

	})
}

func TestBuildBodyPolygot(t *testing.T) {
	const (
		templateId = "fakeId"
		to         = "fakeemail@com"
		country    = "United Kingdom"
		city       = "London"
		language   = "javascript,golang"
		level      = "beginner"
		frameworks = ""
	)
	emailBody, err := digest.BuildBody(to, templateId, country, city, language, level, frameworks)

	if err != nil {
		t.Fatalf(`buildBody(%q, %q, %q, %q, %q, %q, %q) got err %q`, to, templateId, country, city, language, level, frameworks, err)
	}

	t.Run("Test Email Subject", func(t *testing.T) {
		expected := "☕ Top Tech Stack Resources You Have To Read"
		actual, err := jsonparser.GetString(emailBody, "Messages", "[0]", "Subject")

		if err != nil {
			t.Fatalf("unable to get json string. Got %q", err)
		}

		if actual != expected {
			t.Fatalf("got %q but expected %q for email subject", actual, expected)
		}
	})
}

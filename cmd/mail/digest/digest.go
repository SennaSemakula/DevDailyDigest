package digest

import (
	"bytes"
	"github.com/Pioneersltd/DevDailyDigest/v1/pkg/client"
	nl "github.com/Pioneersltd/DevDailyDigest/v1/pkg/newsletter"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	emailUser  = os.Getenv("EMAIL_USER")
	emailToken = os.Getenv("EMAIL_TOKEN")
	tmpl       template.Template
	letter     bytes.Buffer
	articles   []nl.Article
	errLogger  = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	statsMap = make(map[string]int)
)

// Move this
type Podcast struct {
	Title string
	Url string
	Desc string
	Image string

}

var newPodcasts map[string][]Podcast
var newBeginnerArticles map[string][]nl.Article

var beginnerArticles = `
{
	"javascript": [
	  {
		"title": "7 Repos I Didn’t Know I Needed For Front-End",
		"url": "https://javascript.plainenglish.io/7-repos-i-didnt-know-i-needed-for-front-end-389bf498afaa",
		"desc": "",
		"image": "",
		"category": "Javascript 🔥"
	  },
	  {
		"title": "The 50 Best Websites to Learn JavaScript",
		"url": "https://www.codeconquest.com/blog/top-50-websites-to-learn-javascript/",
		"desc": "Did you know that JavaScript is considered the language of the web? And that it is used for a wide variety of online and mobile applications?",
		"image": ""
	  },
	  {
		"title": "How to use the JavaScript console: going beyond console.log()",
		"url": "https://www.freecodecamp.org/news/how-to-use-the-javascript-console-going-beyond-console-log-5128af9d573b/?utm_source=mybridge&utm_medium=blog&utm_campaign=read_more",
		"desc": "If you’re reading this, chances are that you’ve quite recently started learning Web Development...",
		"image": ""
	  },
	  {
		"title": "JavaScript basics",
		"url": "https://developer.mozilla.org/en-US/docs/Learn/Getting_started_with_the_web/JavaScript_basics",
		"desc": "JavaScript is a programming language that adds interactivity to your website. This happens in games, in the behavior of responses when buttons are pressed or with data entry on forms; with dynamic styling; with animation, etc."
	  }
  
	],
	"python": [
	  {
		"title": "4 Ways To Level Up Your Python Code",
		"url": "https://betterprogramming.pub/4-ways-to-level-up-your-python-code-f148a50efeea",
		"desc": "Avoid unnecessary for loops, access dictionary items more effectively, and more",
		"image": "",
		"category": "Python 🔥"
	  },
	  {
		"title": "What is Python? Powerful, intuitive programming",
		"url": "https://www.infoworld.com/article/3204016/what-is-python-powerful-intuitive-programming.html",
		"desc": "Why the Python programming language shines for data science, machine learning, systems automation, web and API development, and more",
		"image": ""
	  },
	  {
		"title": "Learn Python: Tutorials for Beginners, Intermediate, and Advanced Programmers",
		"url": "https://stackify.com/learn-python-tutorials/",
		"desc": "As you know, computers totally depend on program code to function properly. There are so many programming languages available that helps developers create applications.",
		"image": ""
	  }
	],
	"golang": [
	  {
		"title": "Getting started with Golang: a tutorial for beginners",
		"url": "https://www.educative.io/blog/golang-tutorial",
		"desc": "Golang, also known as Go, is an open-source programming language created by Google developers Robert Griesemer, Ken Thompson, and Rob Pike in 2007. It was created for ease, and many developers praise it for building simple, reliable programs.",
		"image": "",
		"category": "Golang 🔥"
	  },
	  {
		"title": "WRITING A SIMPLE PERSISTENT KEY-VALUE STORE IN GO",
		"url": "https://www.opsdash.com/blog/persistent-key-value-store-golang.html",
		"desc": "Want a simple, persistent, key-value store in Go? Something handy to have in your toolbox that’s easy to use and performant? Let’s write one!",
		"image": ""
	  },
	  {
		"title": "Understand Go pointers in less than 800 words or your money back",
		"url": "https://dave.cheney.net/2017/04/26/understand-go-pointers-in-less-than-800-words-or-your-money-back",
		"desc": "This post is for programmers coming to Go who are unfamiliar with the idea of pointers or a pointer type in Go.",
		"image": ""
	  }
	]
}
`

var podcasts = `
{
	"javascript": [
	  {
		"title": "JavaScript Jabber",
		"episode": "JSJ 476: Understanding Search Engines and SEO (for devs) – Part 1",
		"url": "JSJ 476: Understanding Search Engines and SEO (for devs) – Part 1",
		"desc": "If you're building a website or web-app, there's a good chance that you want people to find it so that they will access it. ",
		"image": "https://devchat.tv/wp-content/uploads/2020/06/javascript_jabber_thumb2800-scaled-1500x250.jpg"
	  }
	],
	"python": [
	  {
		"title": "The Real Python Podcast",
		"episode": "Episode 53: Improving the Learning Experience on Real Python",
		"url": "https://realpython.com/podcasts/rpp/53/",
		"desc": "If you haven’t visited the website lately, then you’re missing out on the updates to realpython.com!",
		"image": "https://files.realpython.com/media/real-python-logo-square.28474fda9228.png"
	  }
	],
	"golang": [
	  {
		"title": "Go Time",
		"episode": "Design philosophy",
		"url": "https://changelog.com/gotime/172",
		"desc": "n this insight-filled episode, Bill Kennedy joins Johnny and Kris to discuss best practices around the design of software in Go",
		"image": "https://cdn.changelog.com/uploads/covers/go-time-original.png?v=63725770357"
	  }
	],
	"ruby": [
	  {
		"title": "Ruby Rogues",
		"episode": "BONUS: Continuing Your Learning Journey by Finding Mentors as an Influencer",
		"url": "https://devchat.tv/ruby-rogues/bonus-continuing-your-learning-journey-by-finding-mentors-as-an-influencer-10/",
		"desc": "Chuck outlines how he's used his podcasts to find mentors to continue his learning journey over 12 years of podcasting.",
		"image": "https://devchat.tv/wp-content/uploads/2020/06/ruby-rogues-thumb2800-scaled-768x768.jpg"
	  }
	],
	"php": [
	  {
		"title": "‎php[podcast] episodes from php[architect]",
		"episode": "Editor Bytes – Lambda PHP",
		"url": "https://www.phparch.com/podcast/editor-bytes-lambda-php/",
		"desc": "Oscar Merida, the Editor-in-Chief, looks at the March 2021 issue, Lambda PHP.",
		"image": "https://cdn-images-1.listennotes.com/podcasts/phppodcast-episodes-from-phparchitect-LdHU52lKvV8.1400x1400.jpg"
	  }
	],
	"java": [
	  {
		"title": "Inside Java",
		"episode": "Episode 15 “Java 16 is Here!” with Mikael Vidstedt",
		"url": "https://inside.java/2021/03/16/podcast-015/",
		"desc": "The release of Java 16 was a good reason to invite Mikael Vidstedt (Director of JVM Engineering) again to the show. In this episode, Chad and Mikael discuss some of the new JDK 16 features, the 6 months release cadence but also how some Valhalla initial bits are starting to gradually appear into the platform, and more!",
		"image": "https://inside.java/images/InsideJavaPodcast1000.png"
	  }
	],
	"c++": [
	  {
		"title": "CPPCast",
		"episode": "Freestanding Update",
		"url": "https://cppcast.com/freestanding-update/",
		"desc": "Rob and Jason are joined by Ben Craig. They first discuss blog posts detailing how function call resolution works in C++ and algorithm selection",
		"image": "https://cppcast.com/img/logo-square.png"
	  }
	]
}  
`

type Email struct {
	Email    string
	Country  string
	City     string
	Language string
	Level    string
	Frameworks string
}

type emailStatus struct {
	email  string
	status int
}

func GetBeginnerArticles(language string) ([]nl.Article, error) {

	// Add beginner friendly articles
	if err := json.Unmarshal([]byte(beginnerArticles), &newBeginnerArticles); err != nil {
		return nil, err
	}

	if val, ok := newBeginnerArticles[language]; ok {
		return val, nil
	}

	// TODO: change this to return an error for languages that don't have beginner articles
	return []nl.Article{}, nil
}


func GetArticles(language, frameworks, level string) ([]nl.Article, error) {
	// TODO: this is a monolithic function that needs to get cleared up
	emailArticles := make([]nl.Article, 0)
	tool := ""

	if level == "beginner" {
		beginnerBlogs, err := GetBeginnerArticles(language)

		if err != nil {
			return nil, err
		}
		emailArticles = append(emailArticles, beginnerBlogs...)
	}

	// Only fetch these articles for mid-senior
	if level != "beginner" {
		hackernoonArticles, err := client.FetchArticles("api/v1/hackernoon_articles", language)

		if len(hackernoonArticles) < 3 {
			return nil, fmt.Errorf("Unable to find enough %v hackernoon articles", language)
		}
	
		if err != nil {
			return nil, err
		}

		header := ""

		for i, article := range hackernoonArticles {
			if i == 0 {
				header = strings.Title(language)
			} else {
				header = ""
			}
			newArticle := nl.Article{
				Title: article["Title"].(string),
				Desc:  article["Desc"].(string),
				Image: "",
				Url:   article["Url"].(string),
				Category: header,
			}
	
			emailArticles = append(emailArticles, newArticle)
		}
	}

	if language == "golang" {
		language = language[:2]
	}

	//articles, _ := client.FetchArticles("api/v1/articles", language)
	articles := client.Fetch("https://dev.to/api/articles", language)

	var frameworkArticles []client.ArticlesMap

	if len(frameworks) > 0 {
		// tool := strings.Split(frameworks, ",")[0]
		frameworksSlice := strings.Split(frameworks, ",")

		if tool == "ReactJS" {
			tool = "React"
		}

		// Loop through
		for _, tool := range frameworksSlice {
			if tool == "Tornado" {
				// No dev.to tornado articles
				continue
			}
			if tool == "ReactJS" {
				tool = "react"
			}
			articles, err := client.FetchArticles("api/v1/articles", tool)
			if tool == "Tornado" || tool == "Hibernate" {
				continue
			}

			if err != nil {
				return nil, fmt.Errorf("Unable to fetch %v dev.to articles", tool)
			}

			frameworkArticles = append(frameworkArticles, articles[:3]...)

		}

		// Populate the emailArticles with the framework content 
		for _, article := range frameworkArticles {
			// if article["comments_count"].(float64) < 30 {
			// 	continue
			// }
	
			if article["cover_image"] == nil {
				continue
			}

			newArticle := nl.Article{
				Title: article["title"].(string),
				Desc:  article["description"].(string),
				Date:  article["created_at"].(string),
				Image: "",
				Url:   article["url"].(string),
			}
	
			emailArticles = append(emailArticles, newArticle)
		}

		// articles, _ = client.FetchArticles("api/v1/articles", tool)
	}

	if len(articles) < 3 {
		articles, _ = client.FetchArticles("api/v1/articles", language)
	}

	// For now only show 3 articles; main and two sub articles
	if len(articles) < 3 {
		return nil, fmt.Errorf("Unable to find enough %v dev.to articles", language)
	}

	count := 0

	if len(frameworkArticles) == 0 {

		for _, article := range articles {

			if article["cover_image"] == nil {
				continue
			}
	
			if count >= 4 {
				break
			}
	
			newArticle := nl.Article{
				Title: article["title"].(string),
				Desc:  article["description"].(string),
				Date:  article["created_at"].(string),
				Image: "",
				Url:   article["url"].(string),
			}
	
			emailArticles = append(emailArticles, newArticle)
			count += 1
		}
	}

	// Quick hack
	if language == "go" {
		language = "golang"
	}

	// hashnodeArticles, err := client.FetchArticles("api/v1/hashnode_articles", language)

	// if err != nil {
	// 	errLogger.Printf("Unable to fetch %v hashnode articles")
	// }

	// count = 0
	// for _, article := range hashnodeArticles[2:] {
	// 	if count >= 2 {
	// 		break
	// 	}
	// 	newArticle := nl.Article{
	// 		Title: article["Title"].(string),
	// 		Desc:  article["Desc"].(string),
	// 		Date:  "",
	// 		Image: "",
	// 		Url:  article["Url"].(string),
	// 	}

	// 	emailArticles = append(emailArticles, newArticle)
	// 	count += 1
	// }

	return emailArticles, nil
}

func getJobs(country, city, language, level string) ([]nl.Job, error) {
	var jobs []nl.Job

	jobs, err := client.FetchJobs(country, city, language)

	if err != nil {
		return nil, err
	}

	count := 0
	for _, job := range jobs {
		if count >= 3 {
			break
		}

		if job.Level == level {
			jobs = append(jobs, job)
			count += 1
		}

	}

	return jobs, nil
}

func buildTemplate(tpl template.Template) {
	questions := []nl.Question{
		nl.Question{Question: "What is a channel in Go?"},
	}

	articles, _ = GetArticles("javascript", "gdfgfd", "")

	jobs := []nl.Job{
		nl.Job{Role: "Frontend Engineer", Company: "Monzo", Salary: "£75,000 a year", Location: "London"},
		nl.Job{Role: "Python Developer", Company: "JP Morgan Chase", Salary: "£80,000 a year", Location: "Manchester"},
		nl.Job{Role: "Python Developer", Company: "JP Morgan Chase", Salary: "£80,000 a year", Location: "Manchester"},
	}

	events := []nl.Event{
		nl.Event{Title: "Javascript Event", Location: "Shoreditch", Date: "03/12/20", Time: "17:00 UTC"},
	}

	bundle := nl.Letter{Questions: questions, Articles: articles, Jobs: jobs, Events: events}

	tmpl.Execute(&letter, bundle)
}

func fetchStats(level, language string) map[string]int {
	jobsStat := 0
	eventsStat := 0
	repoStat := 0

	switch level {
	case "beginner":
		jobsStat += 25
		eventsStat += 14
		repoStat += 24
	case "mid":
		jobsStat += 34
		eventsStat += 12
		repoStat += 57
	case "senior":
		jobsStat += 17
		eventsStat += 8
		repoStat += 118
	default:
		jobsStat += 15
		eventsStat += 12
		repoStat += 25
	}

	switch language {
	case "javascript":
		jobsStat += 6
		eventsStat += 6
		repoStat += 6
	case "python":
		jobsStat += 4
		eventsStat += 4
		repoStat += 4
	case "ruby":
		jobsStat += 2
		eventsStat += 2
		repoStat += 2
	case "c++":
		jobsStat += 2
		eventsStat += 2
		repoStat += 2
	case "java":
		jobsStat += 3
		eventsStat += 3
		repoStat += 3
	case "php":
		jobsStat += 1
		eventsStat += 1
		repoStat += 1
	default:
		jobsStat += 4
		eventsStat += 5
		repoStat += 8
	}

	statsMap["jobs"] = jobsStat
	statsMap["events"] = eventsStat
	statsMap["repos"] = repoStat

	return statsMap
}

func FetchEvents(city, language string) ([]nl.Event, error) {
	events, err := client.FetchEvents(city, language)

	if err != nil {
		return nil, err
	}

	if len(events) < 3 {
		events, err = client.FetchEvents("London", language)
	}

	return events, nil

}

func BuildBody(to, templateId, country, city, language, level, frameworks string) ([]byte, error) {
	testArticles := make([]nl.Article, 0)
	articles := make([]nl.Article, 0)
	stats := fetchStats(level, language)

	statsJson, err := json.Marshal(stats)

	if err != nil {
		return nil, err
	}

	// If multi articles
	formattedLang := strings.Split(language, ",")

	subject := ""
	if len(formattedLang) > 1 {
		subject = fmt.Sprintf("☕ Top Tech Stack Resources You Have To Read")
		for _, choice := range formattedLang {
			articles, err := GetArticles(choice, frameworks, level)
			if err != nil {
				return nil, err
			}
			testArticles = append(testArticles, articles...)
		}
		articles = testArticles
		// set language to primary
		language = formattedLang[0]

	} else {
		subject = fmt.Sprintf("☕ Top %v Resources You Have To Read", strings.Title(language))
		holderArticles, err := GetArticles(language, frameworks, level)
		if err != nil {
			return nil, err
		}

		articles = holderArticles
	}

	articlesJson, err := json.Marshal(articles)

	if err != nil {
		return nil, err
	}

	jobs, err := getJobs(country, city, language, level)

	if err != nil {
		return nil, err
	}

	if len(jobs) < 3 {
		jobs, err = getJobs("remote", "remote", language, level)
	}

	if err != nil {
		return nil, err
	} else if len(jobs) < 3 {
		return nil, fmt.Errorf("unable to find %q enough jobs", language)
	}

	jobsJson, err := json.Marshal(jobs[:3])

	if err != nil {
		return nil, err
	}
	// TODO HERE:
	events, err := FetchEvents(city, language)

	if err != nil {
		return nil, err
	} else if len(events) < 3 {
		return nil, fmt.Errorf("unable to find %q enough events", language)
	}

	eventsJson, err := json.Marshal(events[:2])

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(podcasts), &newPodcasts); err != nil {
		fmt.Println(err)
		return nil, err
	}

	podcastsJson, err := json.Marshal(newPodcasts[language]) 

	if err != nil {
		return nil, err
	}


	body := fmt.Sprintf(`
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
				"Name": "DevDaily"
			  }
			],
			"TemplateID": %v,
			"TemplateLanguage": true,
			"Subject": %q,
			"TemplateErrorReporting": {
				"Email": "senna.semakula@gmail.com",
				"Name": "Senna"
			},
			"Variables": {
				"unsubscribe": "https://devdaily.io/unsubscribe",
				"unsubscribe_preferences": "https://devdaily.io/unsubscribe",
				"articles": %v,
				"jobs": %v,
				"events": %v,
				"level": %q,
				"stats": %v,
				"podcasts": %v
				}
		  }
		]
	  }
	`, to, templateId, subject, string(articlesJson), string(jobsJson), string(eventsJson), level, string(statsJson), string(podcastsJson))

	return []byte(body), nil

}

func SendMail(e *Email, templateId string) (int, error) {
	// req := sendgrid.GetRequest(os.Getenv(apiKey), "/v3/mail/send", api)
	// req.Method = "POST"
	url := "https://api.mailjet.com/v3.1/send"

	body, err := BuildBody(e.Email, templateId, e.Country, e.City, e.Language, e.Level, e.Frameworks)
	if err != nil {
		return 1, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	// Auth
	req.SetBasicAuth(emailUser, emailToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		return 1, err
	}

	if resp.StatusCode != 200 {
		errLogger.Printf("unable to send email to %v. status code %v", e.Email, resp.StatusCode)
	} else {
		log.Println("Sent email to ", e.Email, "status:", resp.StatusCode)
	}

	return resp.StatusCode, nil
}

package client

import (
	"time"
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	start time.Time
	url = "https://dev.to/api/articles"
	authToken = os.Getenv("DEV_AUTH_TOKEN")
	client *http.Client
)

func init() {
	start = time.Now()
	client = &http.Client{}
}


func Fetch(url string, tag string) *http.Response{
	api := url + "?=" + tag + "&per_page=100&state=fresh&top=7"
	req, err := http.NewRequest("GET", api + "", nil)
	req.Header.Add("api-key", authToken)
	resp, err := client.Do(req)

	if err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	return resp
}

// Takes in a send only channel
// func AlotOfRequests(url string, soc chan<-*http.Response) {
// 	for i := 0; i <= 100; i++ {
// 		res := Fetch(url)
// 		soc <- res
// 		fmt.Printf("idx %v Wrote to channel\n", i)
// 	}

// 	close(soc)
// }


// func main() {
// 	fmt.Println("Main Goroutine started")
// 	// c := make(chan http.Response)
// 	fmt.Println(authToken)
// 	Fetch(url, "golang")
// 	// c := make(chan *http.Response, 11)

// 	// for i:= 0; i < 10; i++ {
// 	// 	go func(url string){
// 	// 		res := Fetch(url)
// 	// 		c <- res
// 	// 	}(url)
// 	// }

// 	// // Closes channel by default
// 	// for resp := range c {
// 	// 	defer resp.Body.Close()
// 	// 	_, err := ioutil.ReadAll(resp.Body)

// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}
// 	// 	fmt.Printf("Popping %v from the channel\r\n", resp.Status)
// 	// }

// 	end := time.Now()

// 	fmt.Printf("Operation took %v\n", end.Sub(start))
// 	fmt.Println("Main Goroutine ended")
// 	// log.Println(res)
// }
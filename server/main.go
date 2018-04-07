package main

import (
	"encoding/json"
  "os"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Config struct {
  ConsumerKey string
  ConsumerSecret string
  Token string
  TokenSecret string
}

func increment(alphanum string) string {
	chars := strings.Split(alphanum, "")
	result := make([]byte, len(chars))
	lastchar := len(chars) - 1
	for index, element := range chars {
		result[index] = element[0]
	}
	for lastchar >= 0 {
		if result[lastchar] != 'z' && result[lastchar] != '9' {
			result[lastchar] = result[lastchar] + 1
			break
		} else if result[lastchar] == 'z' {
			result[lastchar] = '0'
		} else {
			result[lastchar] = 'a'
		}
		lastchar--
	}
	return string(result)
}

func decrement(alphanum string) string {
	chars := strings.Split(alphanum, "")
	result := make([]byte, len(chars))
	lastchar := len(chars) - 1
	for index, element := range chars {
		result[index] = element[0]
	}
	for lastchar >= 0 {
		if result[lastchar] != 'a' && result[lastchar] != '0' {
			result[lastchar] = result[lastchar] - 1
			break
		} else if result[lastchar] == 'a' {
			result[lastchar] = '9'
		} else {
			result[lastchar] = 'z'
		}
		lastchar--
	}
	return string(result)
}

func image(res http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/image/")
	fmt.Println(id)
	if len(id) != 6 {
		io.WriteString(res, "<!DOCTYPE html><html><body><p>Invalid Request</p></body></html>")
		return
	}
	client := http.Client{}
	myreq, _ := http.NewRequest("GET", "https://prnt.sc/"+id, nil)
	myreq.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0")
	resp, _ := client.Do(myreq)
	bytes, _ := ioutil.ReadAll(resp.Body)
	code := string(bytes)
	r1, _ := regexp.Compile("og:image\" content=\".*png\"/>")
	linkr, _ := regexp.Compile("https:.*png")
	meta := r1.FindString(code)
	image := linkr.FindString(meta)
	resp.Body.Close()
	data := make(map[string]string)
	data["image"] = image
	data["next"] = increment(id)
	data["prev"] = decrement(id)
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(data)
}

func getUrls(tweet twitter.Tweet) []string {
	result := make([]string, 0, 10)
	for _, url := range tweet.Entities.Urls {
		if strings.HasPrefix(url.ExpandedURL, "http://prntscr.com") ||
		   strings.HasPrefix(url.ExpandedURL, "https://prntscr.com") ||
		   strings.HasPrefix(url.ExpandedURL, "http://prnt.sc") ||
		   strings.HasPrefix(url.ExpandedURL, "https://prnt.sc") {
			result = append(result, url.ExpandedURL)
		}
	}
	return result
}

func getRecent() string {
  configfile, _ := os.Open("config.json")
  defer configfile.Close()
  decoder := json.NewDecoder(configfile)
  creds := Config{}
  decoder.Decode(&creds)
  fmt.Println(creds)
  config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
  token := oauth1.NewToken(creds.Token, creds.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	search, _, _ := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: "prnt.sc",
		Count: 100,
	})

  // Try and find a recent link
  link := ""
  for _, tweet := range search.Statuses {
    urls := getUrls(tweet)
    if len(urls) > 0 {
      link = urls[0]
      break
    }
  }

  return link
}

func recent(res http.ResponseWriter, req *http.Request) {
  data := make(map[string]string)
  data["recent"] = strings.Split(getRecent(), "/")[3]
  res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(data)
}

func main() {
	http.HandleFunc("/image/", image)
  http.HandleFunc("/recent/", recent)
	http.ListenAndServe(":12345", nil)
}

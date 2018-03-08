package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

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

func main() {
	http.HandleFunc("/image/", image)
	http.ListenAndServe(":12345", nil)
}

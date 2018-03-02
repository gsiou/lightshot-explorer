package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "regexp"
  "strings"
)

func image (res http.ResponseWriter, req *http.Request) {
  id := strings.TrimPrefix(req.URL.Path, "/image/")
  fmt.Println(id)
  if len(id) != 6 {
    io.WriteString(res, "<!DOCTYPE html><html><body><p>Invalid Request</p></body></html>")
    return 
  }
  client := http.Client{}
  myreq, _ := http.NewRequest("GET", "https://prnt.sc/" + id, nil)
  myreq.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0")
  resp, _ := client.Do(myreq)
  bytes, _ := ioutil.ReadAll(resp.Body)
  code := string(bytes)
  r1, _ := regexp.Compile("og:image\" content=\".*png\"/>")
  linkr, _ := regexp.Compile("https:.*png")
  meta := r1.FindString(code)
  image := linkr.FindString(meta)
  fmt.Println(image)
  resp.Body.Close()
  io.WriteString(res, "<!DOCTYPE html><html><body><img src='" + image + "'/></body></html>")  
}

func main () {
     http.HandleFunc("/image/", image)
     http.ListenAndServe(":12345", nil)
     fmt.Println("Listening ...")
}

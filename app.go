package main

import (
    "io"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", "http://173.192.61.226" + r.URL.Path, nil)
    req.Host = "disqus.com"
    resp, err := client.Do(req)
    if err != nil {
        //
    }
    defer resp.Body.Close()
    io.Copy(w, resp.Body)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":6081", nil)
}

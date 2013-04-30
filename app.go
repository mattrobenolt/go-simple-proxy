package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
    "encoding/json"
)

type Stats struct {
    NumFetches int
    NumRequests int
}

var (
    // Map to hold all responses
    // Queue is mapped based on request path naiively
    responseQueue = map[string][]chan []byte{}

    // Total number of origin fetches
    numFetches = 0

    // Total number of requests
    numRequests = 0

    stats = Stats{}
)

func fetchURL(path string) {
    stats.NumFetches++

    client := &http.Client{}
    req, err := http.NewRequest("GET", "http://75.101.142.194" + path, nil)
    req.Host = "httpbin.org"
    resp, err := client.Do(req)
    if err != nil {
        //
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    // We're manipuating the global map, so we want to lock while mutating
    mutex := new(sync.Mutex)
    mutex.Lock()
    // Copy the existing queue
    queue := make([]chan []byte, len(responseQueue[path]))
    copy(queue, responseQueue[path])
    // Now delete to free it up for new requests
    delete(responseQueue, path)
    mutex.Unlock()

    // Respond to all waiting channels
    for _, ch := range queue {
        ch <- body
    }
}

func handleMain(w http.ResponseWriter, r *http.Request) {
    stats.NumRequests++

    // Start with a new response channel
    ch := make(chan []byte)
    path := r.URL.Path

    // Lock down the map while we manipulate it
    mutex := new(sync.Mutex)
    mutex.Lock()
    if queue, ok := responseQueue[path]; ok {
        // A queue already exists, just append this channel to it
        queue = append(queue, ch)
    } else {
        // No queue, create one with your channel already in it
        queue := []chan []byte{ch}
        responseQueue[path] = queue
        // Start the origin fetch
        go fetchURL(path)
    }
    mutex.Unlock()
    // Write out the response
    fmt.Fprintf(w, "%s", <- ch)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
    b, _ := json.Marshal(stats)
    fmt.Fprintf(w, "%s", b)
}

func main() {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/stats", handleStats)
    http.ListenAndServe(":6081", nil)
}

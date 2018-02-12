package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", handle)
    appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, world!")
}

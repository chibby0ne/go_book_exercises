package main

import (
    "github.com/chibby0ne/go_book_exercises/chapter12/params"
    "net/http"
    "fmt"
)

func search(resp http.ResponseWriter, req *http.Request) {
    var data struct {
        Labels []string `http:"l"`
        MaxResults  int `http:"max"`
        Exact bool `http:"x"`
    }
    data.MaxResults = 10 // set default
    if err := params.Unpack(req, &data); err != nil {
        http.Error(resp, err.Error(), http.StatusBadRequest) // 400
        return
    }
    // ... rest of handler
    fmt.Fprintf(resp, "Search: %+v\n", data)
}


func main() {
    http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        search(w, r)
    })
    http.ListenAndServe("localhost:12345", nil)
}

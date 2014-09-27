package main

import (
  // "html/template"
  "net/http"
  "fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "hi")
}

func main() {
  // http.HandleFunc("/", indexHandler)
  // http.ListenAndServe(":8080", nil)
  readSections()
}

package main

import (
  "net/http"
  "fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  sections := readSections()
  indexTempl := pageCache.get("main", "_base", "main")
  indexTempl.ExecuteTemplate(w, "page", sections)
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
  avatarTempl := pageCache.get("avatar", "_base", "avatar")
  avatarTempl.ExecuteTemplate(w, "page", nil)
}

func loggingHandler(handler func(http.ResponseWriter, *http.Request) ) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s  %s", r.RemoteAddr, r.URL)
    handler(w, r)
    fmt.Print("\n")
  }
}

func main() {
  http.HandleFunc("/", loggingHandler(indexHandler))
  http.HandleFunc("/avatar", loggingHandler(avatarHandler))
  http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("public"))))
  http.ListenAndServe(":8080", nil)
}

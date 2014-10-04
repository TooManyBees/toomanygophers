package main

import (
  "net/http"
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

func main() {
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/avatar", avatarHandler)
  http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("public"))))
  http.ListenAndServe(":8080", nil)
}

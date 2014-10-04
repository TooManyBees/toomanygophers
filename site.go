package main

import (
  "html/template"
  "net/http"
)

const templateDir = "templates"

func indexHandler(w http.ResponseWriter, r *http.Request) {
  sections := readSections()
  indexTempl := template.Must(template.ParseFiles("templates/_base.html", "templates/main.html"))
  indexTempl.ExecuteTemplate(w, "page", sections)
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
  avatarTempl := template.Must(template.ParseFiles("templates/_base.html", "templates/avatar.html"))
  avatarTempl.ExecuteTemplate(w, "page", nil)
}

func main() {
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/avatar", avatarHandler)
  http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("public"))))
  http.ListenAndServe(":8080", nil)
}

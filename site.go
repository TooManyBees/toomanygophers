package main

import (
  "html/template"
  "net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  sections := readSections()
  indexTempl := template.Must(template.ParseFiles("templates/main/index.html"))
  indexTempl.Execute(w, sections)
}

func main() {
  http.HandleFunc("/", indexHandler)
  http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("public"))))
  http.ListenAndServe(":8080", nil)
}

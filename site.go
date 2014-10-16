package main

import (
  "net/http"
  "fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  sections := readSections()
  indexTempl := pageCache.get("main")
  indexTempl.ExecuteTemplate(w, "page", sections)
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
  avatarTempl := pageCache.get("avatar")
  avatarTempl.ExecuteTemplate(w, "page", nil)
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    comics := comicStore.random(20)
    // for _, c := range comics {
    //   fmt.Fprintln(w, c.Title)
    // }
    quizTempl := pageCache.get("quiz")
    quizTempl.ExecuteTemplate(w, "page", comics)
  } else if r.Method == "POST" {
    //
  }
}

// Wraps an endpoint handler in logging statements, as well as
// enforcing the allowed HTTP methods for it
func loggingHandler(handler func(http.ResponseWriter, *http.Request), methods ...string) http.HandlerFunc {
  includes := func(needle string, haystack ...string) bool {
    for _, str := range haystack {
      if str == needle {
        return true
      }
    }
    return false
  }

  return func(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s  %s", r.RemoteAddr, r.URL)
    if len(methods) == 0 || includes(r.Method, methods...) {
      handler(w, r)
    } else {
      fmt.Print("  unsupported method")
      w.WriteHeader(http.StatusMethodNotAllowed)
    }
    fmt.Print("\n")
  }
}

var pageCache = PageCache{}
var comicStore = ComicStore{}
func main() {
  comicStore.init()
  http.HandleFunc("/", loggingHandler(indexHandler, "GET"))
  http.HandleFunc("/avatar", loggingHandler(avatarHandler, "GET"))
  http.HandleFunc("/quiz", loggingHandler(quizHandler, "GET", "POST"))
  http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("public"))))
  http.ListenAndServe(":8080", nil)
}

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
func main() {
  http.HandleFunc("/", loggingHandler(indexHandler, "GET"))
  http.HandleFunc("/avatar", loggingHandler(avatarHandler, "GET"))
  http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("public"))))
  http.ListenAndServe(":8080", nil)
}

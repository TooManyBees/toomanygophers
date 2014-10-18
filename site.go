package main

import (
  "net/http"
  "fmt"
  "os"
  "path/filepath"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path == "/" {
    sections := readSections()
    indexTempl := pageCache.get("main")
    indexTempl.ExecuteTemplate(w, "page", sections)
  } else {
    fmt.Print("  not found")
    w.WriteHeader(http.StatusNotFound)
  }
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
    fmt.Printf("%s  %s", r.RemoteAddr, r.URL.Path)
    if len(methods) == 0 || includes(r.Method, methods...) {
      handler(w, r)
    } else {
      fmt.Print("  unsupported method")
      w.WriteHeader(http.StatusMethodNotAllowed)
    }
    fmt.Print("\n")
  }
}

func deployAssets(env string) {
  // TODO: make this depend on environment (uglify instead of link in prod)
  scripts, _ := filepath.Glob("assets/javascripts/*.js")
  cwd, _ := os.Getwd()
  for _, script := range scripts {
    filename := filepath.Base(script)
    publicScript := filepath.Join("public", filename)
    _, err := os.Lstat(publicScript)
    if os.IsNotExist(err) {
      linkErr := os.Symlink(filepath.Join(cwd, script), publicScript)
      if linkErr != nil {
        fmt.Println(linkErr)
      }
    }
  }
}

func parseOptions() (string, string) {
  env := "development"
  port := "8080"
  if len(os.Args) > 1 {
    args := os.Args[1:]
    lastIndex := len(args) - 1
    for i, arg := range args {
      if arg == "-e" && i < lastIndex {
        env = args[i+1]
      } else if arg == "-p" && i < lastIndex {
        port = args[i+1]
      }
    }
  }
  return env, port
}

var pageCache = PageCache{}
var comicStore = ComicStore{}
func main() {
  env, port := parseOptions()
  deployAssets(env)
  comicStore.init()
  http.HandleFunc("/", loggingHandler(indexHandler, "GET"))
  http.HandleFunc("/avatar", loggingHandler(avatarHandler, "GET"))
  http.HandleFunc("/quiz", loggingHandler(quizHandler, "GET", "POST"))
  http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("public"))))
  fmt.Printf("Running in %s on port %s\n", env, port)
  http.ListenAndServe(":"+port, nil)
}

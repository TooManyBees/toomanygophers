package main

import (
  _ "github.com/lib/pq"
  "database/sql"
  "fmt"
  "math/rand"
  "time"
  "strings"
)

type Quiz struct {
  Id string
  Comics string
}

type Comic struct {
  Id int
  Title string
  Image string
}

type ComicStore struct {
  db *sql.DB
  count int
  dirty bool
}

func (cs *ComicStore) init() {
  db, err := sql.Open("postgres", "dbname=spam_quiz sslmode=disable")
  if err != nil {
    fmt.Printf("Couldn't connect to database: %s\n", err)
    return
  }
  cs.dirty = true
  cs.db = db
}

func (cs *ComicStore) debug() {
  rows, _ := cs.db.Query("SELECT * FROM comics ORDER BY id")
  defer rows.Close()
  for rows.Next() {
    var id int
    var title string
    var image string
    var altImage string
    rows.Scan(&id, &title, &image, &altImage)
    fmt.Printf("%d: %s => %s (%s)\n", id, title, image, altImage)
  }
}

func (cs *ComicStore) Count() int {
  if cs.dirty {
    row := cs.db.QueryRow("SELECT count(id) FROM comics")
    row.Scan(&cs.count)
    cs.dirty = false
  }
  return cs.count
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
func sqlIn(arr []int) string {
  raw := fmt.Sprint(arr)
  replaced := strings.Replace(raw, " ", ",", -1)
  return replaced[1:len(replaced)-1]
}

func (cs *ComicStore) random(n int) []Comic {

  if n > cs.Count() {
    return nil
  }

  p := rng.Perm(int(cs.Count()))
  q := fmt.Sprintf("SELECT * FROM comics WHERE id IN (%s)", sqlIn(p[:n]))
  rows, err := cs.db.Query(q)
  if err != nil {
    fmt.Println(err)
  }

  var comics = make([]Comic, 0, n)

  for rows.Next() {
    var id int
    var title string
    var image string
    var alt_image string
    rows.Scan(&id, &title, &image, &alt_image)
    comics = append(comics, Comic{id, title, image})
  }
  return comics
}

func main() {
  cs := ComicStore{}
  cs.init()
  fmt.Println(cs.Count())
  fmt.Println(cs.random(5))
}
